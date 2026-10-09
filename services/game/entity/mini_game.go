package entity

import (
	"fmt"
	"math/rand"

	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
)

const (
	MiniGameOmokSets     = 12
	MiniGameMatchCardSet = 4080100
	miniGameOmokSetBase  = 4080000
	miniGameOmokSize     = 15
	miniGameOmokLine     = 5
)

var miniGameCardPairs = []int{6, 10, 15}

type MiniGame struct {
	ObjectCore
	Kind      uint8
	Title     string
	Password  string
	Piece     uint8
	Turn      uint8
	Cards     []uint32
	owner     *Character
	visitor   *Character
	ready     bool
	playing   bool
	exitAfter [2]bool
	tieAsked  bool
	tieSlot   uint8
	board     [miniGameOmokSize][miniGameOmokSize]uint8
	stones    [2]uint8
	placed    int
	matched   []bool
	picked    bool
	firstCard uint8
	points    [2]int
}

func (g *MiniGame) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeMiniGame
}

func (g *MiniGame) Is(typ constant.ObjectType) bool {
	return g.GetObjectType().Has(typ)
}

func (g *MiniGame) SlotOf(ch *Character) (uint8, bool) {
	switch {
	case ch == g.owner:
		return 0, true
	case ch == g.visitor && ch != nil:
		return 1, true
	}
	return 0, false
}

func (g *MiniGame) player(slot uint8) *Character {
	if slot == 0 {
		return g.owner
	}
	return g.visitor
}

func (g *MiniGame) Members() []*Character {
	if g.visitor == nil {
		return []*Character{g.owner}
	}
	return []*Character{g.owner, g.visitor}
}

func (g *MiniGame) recordKey(stat string) string {
	if g.Kind == pconst.MiniRoomTypeOmok {
		return "omok." + stat
	}
	return "match_card." + stat
}

func (g *MiniGame) Record(ch *Character) response.MiniGameRecord {
	wins := int32(ch.Records.Get(g.recordKey("wins")))
	ties := int32(ch.Records.Get(g.recordKey("ties")))
	losses := int32(ch.Records.Get(g.recordKey("losses")))
	return response.MiniGameRecord{
		Type:   uint32(g.Kind),
		Wins:   wins,
		Ties:   ties,
		Losses: losses,
		Score:  2000 + wins*2 + ties - losses*2,
	}
}

func (g *MiniGame) balloon() response.MiniRoomBalloon {
	users := uint8(1)
	if g.visitor != nil {
		users++
	}
	return response.MiniRoomBalloon{
		Type:     g.Kind,
		SN:       g.OID,
		Title:    g.Title,
		Private:  g.Password != "",
		Spec:     g.Piece,
		Users:    users,
		MaxUsers: pconst.MiniRoomGameUsers,
		Playing:  g.playing,
	}
}

func (g *MiniGame) updateBalloon() {
	balloon := g.balloon()
	g.owner.Listener.OnMiniRoomBalloon(g.owner, &balloon)
}

func (ch *Character) CreateMiniGame(kind uint8, title string, password string, piece uint8) error {
	m := ch.GetMap()
	if m == nil || ch.MiniRoom != nil || ch.miniRoomPending {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterCannotOpen}
	}
	if m.Wz.PersonalShop || m.Wz.Limits(constant.FieldLimitMinigame) {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterCannotOpen}
	}
	if title == "" {
		return ErrMiniRoomInvalid
	}

	var itemID uint32
	switch kind {
	case pconst.MiniRoomTypeOmok:
		if piece >= MiniGameOmokSets {
			return ErrMiniRoomInvalid
		}
		itemID = miniGameOmokSetBase + uint32(piece)
	case pconst.MiniRoomTypeMatchCard:
		if int(piece) >= len(miniGameCardPairs) {
			return ErrMiniRoomInvalid
		}
		itemID = MiniGameMatchCardSet
	default:
		return ErrMiniRoomInvalid
	}
	if !ch.Inventory.HasItem(itemID) {
		return ErrMiniRoomInvalid
	}
	if err := m.checkShopSpot(ch.Position); err != nil {
		return err
	}

	g := &MiniGame{
		Kind:     kind,
		Title:    title,
		Password: password,
		Piece:    piece,
		owner:    ch,
	}
	g.ObjectCore.self = g
	g.Position = ch.Position
	m.AddMiniGame(g)
	ch.MiniRoom = g
	ch.Listener.OnMiniGameEntered(ch, g)
	g.updateBalloon()
	return nil
}

func (g *MiniGame) Visit(ch *Character, password string) error {
	if ch.MiniRoom != nil || ch.miniRoomPending {
		return ErrMiniRoomInvalid
	}
	if g.visitor != nil {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterFull}
	}
	if g.Password != "" && g.Password != password {
		ch.Listener.OnMessage(ch, constant.MsgPopup, "입력한 비밀번호가 올바르지 않습니다.")
		return ErrMiniRoomInvalid
	}

	g.owner.Listener.OnMiniGameVisited(g.owner, g, ch)
	g.visitor = ch
	ch.MiniRoom = g
	ch.Listener.OnMiniGameEntered(ch, g)
	g.updateBalloon()
	return nil
}

func (g *MiniGame) Chat(ch *Character, message string) error {
	slot, ok := g.SlotOf(ch)
	if !ok {
		return ErrMiniRoomInvalid
	}

	text := fmt.Sprintf("%s : %s", ch.GetName(), message)
	for _, member := range g.Members() {
		member.Listener.OnMiniRoomChat(member, slot, text)
	}
	return nil
}

func (g *MiniGame) Leave(ch *Character) {
	slot, ok := g.SlotOf(ch)
	if !ok {
		return
	}
	if g.playing {
		g.finish(pconst.MiniGameGiveUp, 1-slot)
	}
	if ch.MiniRoom == nil {
		return
	}

	if slot == 0 {
		g.close()
		return
	}
	g.removeVisitor()
}

func (g *MiniGame) removeVisitor() {
	g.visitor.MiniRoom = nil
	g.visitor = nil
	g.ready = false
	g.exitAfter[1] = false
	g.owner.Listener.OnMiniRoomLeft(g.owner, 1, pconst.MiniRoomLeaveExit)
	g.updateBalloon()
}

func (g *MiniGame) close() {
	if g.visitor != nil {
		visitor := g.visitor
		g.visitor = nil
		visitor.MiniRoom = nil
		visitor.Listener.OnMiniRoomLeft(visitor, 1, pconst.MiniRoomLeaveShopClosed)
	}
	owner := g.owner
	owner.MiniRoom = nil
	owner.Listener.OnMiniRoomBalloon(owner, nil)
	if g.Map != nil {
		g.Map.RemoveMiniGame(g)
	}
}

func (g *MiniGame) Ready(ch *Character, ready bool) error {
	if ch != g.visitor || g.playing {
		return ErrMiniRoomInvalid
	}

	g.ready = ready
	for _, member := range g.Members() {
		member.Listener.OnMiniGameReady(member, ready)
	}
	return nil
}

func (g *MiniGame) Expel(ch *Character) error {
	if ch != g.owner || g.visitor == nil || g.playing {
		return ErrMiniRoomInvalid
	}

	visitor := g.visitor
	g.removeVisitor()
	visitor.Listener.OnMiniRoomLeft(visitor, 1, pconst.MiniRoomLeaveKicked)
	return nil
}

func (g *MiniGame) Start(ch *Character) error {
	if ch != g.owner || g.visitor == nil || !g.ready || g.playing {
		return ErrMiniRoomInvalid
	}

	g.playing = true
	g.ready = false
	g.tieAsked = false
	g.board = [miniGameOmokSize][miniGameOmokSize]uint8{}
	g.stones = [2]uint8{}
	g.placed = 0
	g.points = [2]int{}
	g.picked = false
	g.Cards = nil

	if g.Kind == pconst.MiniRoomTypeMatchCard {
		pairs := miniGameCardPairs[g.Piece]
		g.Cards = make([]uint32, 0, pairs*2)
		for i := 0; i < pairs; i++ {
			g.Cards = append(g.Cards, uint32(i), uint32(i))
		}
		rand.Shuffle(len(g.Cards), func(i, j int) {
			g.Cards[i], g.Cards[j] = g.Cards[j], g.Cards[i]
		})
		g.matched = make([]bool, len(g.Cards))
	}

	for _, member := range g.Members() {
		member.Listener.OnMiniGameStarted(member, g)
	}
	g.updateBalloon()
	return nil
}

func (g *MiniGame) turnOf(ch *Character) (uint8, error) {
	slot, ok := g.SlotOf(ch)
	if !ok || !g.playing || slot != g.Turn {
		return 0, ErrMiniRoomInvalid
	}
	return slot, nil
}

func (g *MiniGame) MoveOmok(ch *Character, x int32, y int32, stone uint8) error {
	slot, err := g.turnOf(ch)
	if err != nil {
		return err
	}
	if g.Kind != pconst.MiniRoomTypeOmok || x < 0 || y < 0 || x >= miniGameOmokSize || y >= miniGameOmokSize {
		return ErrMiniRoomInvalid
	}
	if stone != 1 && stone != 2 {
		return ErrMiniRoomInvalid
	}
	if g.stones[slot] != 0 && g.stones[slot] != stone {
		return ErrMiniRoomInvalid
	}
	if g.stones[1-slot] == stone || g.board[x][y] != 0 {
		return ErrMiniRoomInvalid
	}

	g.stones[slot] = stone
	g.board[x][y] = stone
	g.placed++
	for _, member := range g.Members() {
		member.Listener.OnMiniGameOmokMoved(member, x, y, stone)
	}

	switch {
	case g.fiveInRow(int(x), int(y), stone):
		g.finish(pconst.MiniGameWin, slot)
	case g.placed == miniGameOmokSize*miniGameOmokSize:
		g.finish(pconst.MiniGameTie, 0)
	default:
		g.Turn = 1 - slot
	}
	return nil
}

func (g *MiniGame) fiveInRow(x int, y int, stone uint8) bool {
	for _, dir := range [][2]int{{1, 0}, {0, 1}, {1, 1}, {1, -1}} {
		count := 1
		for _, sign := range []int{1, -1} {
			for step := 1; ; step++ {
				nx, ny := x+dir[0]*step*sign, y+dir[1]*step*sign
				if nx < 0 || ny < 0 || nx >= miniGameOmokSize || ny >= miniGameOmokSize || g.board[nx][ny] != stone {
					break
				}
				count++
			}
		}
		if count >= miniGameOmokLine {
			return true
		}
	}
	return false
}

func (g *MiniGame) SelectCard(ch *Character, firstPick bool, card uint8) error {
	slot, err := g.turnOf(ch)
	if err != nil {
		return err
	}
	if g.Kind != pconst.MiniRoomTypeMatchCard || int(card) >= len(g.Cards) || g.matched[card] {
		return ErrMiniRoomInvalid
	}
	if firstPick == g.picked || (g.picked && card == g.firstCard) {
		return ErrMiniRoomInvalid
	}

	if firstPick {
		g.picked = true
		g.firstCard = card
		opponent := g.player(1 - slot)
		opponent.Listener.OnMiniGameCardSelected(opponent, true, card, 0, 0)
		return nil
	}

	first := g.firstCard
	g.picked = false
	result := slot
	matched := g.Cards[first] == g.Cards[card]
	if matched {
		g.matched[first] = true
		g.matched[card] = true
		g.points[slot]++
		result += 2
	} else {
		g.Turn = 1 - slot
	}
	for _, member := range g.Members() {
		member.Listener.OnMiniGameCardSelected(member, false, card, first, result)
	}
	if !matched || g.points[0]+g.points[1] < len(g.Cards)/2 {
		return nil
	}

	switch {
	case g.points[0] == g.points[1]:
		g.finish(pconst.MiniGameTie, 0)
	case g.points[0] > g.points[1]:
		g.finish(pconst.MiniGameWin, 0)
	default:
		g.finish(pconst.MiniGameWin, 1)
	}
	return nil
}

func (g *MiniGame) Skip(ch *Character) error {
	slot, err := g.turnOf(ch)
	if err != nil {
		return err
	}

	g.picked = false
	g.Turn = 1 - slot
	for _, member := range g.Members() {
		member.Listener.OnMiniGameSkipped(member, g.Turn)
	}
	return nil
}

func (g *MiniGame) GiveUp(ch *Character) error {
	slot, ok := g.SlotOf(ch)
	if !ok || !g.playing {
		return ErrMiniRoomInvalid
	}

	g.finish(pconst.MiniGameGiveUp, 1-slot)
	return nil
}

func (g *MiniGame) RequestTie(ch *Character) error {
	slot, ok := g.SlotOf(ch)
	if !ok || !g.playing || g.tieAsked {
		return ErrMiniRoomInvalid
	}

	g.tieAsked = true
	g.tieSlot = slot
	opponent := g.player(1 - slot)
	opponent.Listener.OnMiniGameTieRequested(opponent)
	return nil
}

func (g *MiniGame) AnswerTie(ch *Character, accept bool) error {
	slot, ok := g.SlotOf(ch)
	if !ok || !g.playing || !g.tieAsked || g.tieSlot == slot {
		return ErrMiniRoomInvalid
	}

	g.tieAsked = false
	if accept {
		g.finish(pconst.MiniGameTie, 0)
		return nil
	}
	requester := g.player(g.tieSlot)
	requester.Listener.OnMiniGameTieDenied(requester)
	return nil
}

func (g *MiniGame) ExitAfterGame(ch *Character, reserve bool) error {
	slot, ok := g.SlotOf(ch)
	if !ok || !g.playing {
		return ErrMiniRoomInvalid
	}

	g.exitAfter[slot] = reserve
	ch.Listener.OnMiniGameExitAfter(ch, reserve)
	return nil
}

func (g *MiniGame) finish(outcome pconst.MiniGameOutcome, winner uint8) {
	g.playing = false
	g.tieAsked = false
	g.picked = false

	reset := RecordReset{}
	switch outcome {
	case pconst.MiniGameTie:
		g.owner.Records.Add(g.recordKey("ties"), 1, reset)
		g.visitor.Records.Add(g.recordKey("ties"), 1, reset)
		g.Turn = 1 - g.Turn
	default:
		g.player(winner).Records.Add(g.recordKey("wins"), 1, reset)
		g.player(1-winner).Records.Add(g.recordKey("losses"), 1, reset)
		g.Turn = 1 - winner
	}
	for _, member := range g.Members() {
		member.Listener.OnMiniGameOver(member, g, outcome, winner)
	}
	g.updateBalloon()

	switch {
	case g.exitAfter[0]:
		g.exitAfter = [2]bool{}
		g.owner.Listener.OnMiniRoomLeft(g.owner, 0, pconst.MiniRoomLeaveExit)
		g.close()
	case g.exitAfter[1]:
		visitor := g.visitor
		g.removeVisitor()
		visitor.Listener.OnMiniRoomLeft(visitor, 1, pconst.MiniRoomLeaveExit)
	}
}
