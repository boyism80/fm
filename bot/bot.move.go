package bot

import (
	"cmp"
	"container/list"
	"maps"
	"math"
	"slices"
	"sync"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/types"
)

// Map.wz/Physics.img as the client reads it, with its default shoe (mass 100, every other factor 1) at 100% speed
// and jump. The jump, rope, and climb factors are literals of the client's movement code.
const (
	walkForce        = 140000.0
	walkSpeed        = 125.0
	walkDrag         = 80000.0
	floatDrag2       = 10000.0
	floatCoefficient = 0.01
	gravityAcc       = 2000.0
	fallSpeed        = 670.0
	jumpSpeed        = 555.0
	shoeMass         = 100.0
	jumpBoost        = 0.8
	jumpDownSpeed    = -0.35355339 * jumpSpeed
	ropeJumpSpeed    = 0.5 * jumpSpeed
	ropeJumpBoost    = 1.3 * walkSpeed
	climbStep        = 3.0
	floatDrag1       = 100000.0
	maxFriction      = 2.0
	minFriction      = 0.05
	conveyorAgainst  = 0.2
	swimForce        = 120000.0
	swimSpeed        = 140.0
	swimSpeedDec     = 0.9
	swimJump         = 0.7
	swimStroke       = 5.0
	swimHover        = 0.3
	swimDive         = 1.5
)

const (
	// The client steps its physics every 30 ms and sends the local user's path every 500 ms.
	moveTick  = 30
	moveFlush = 500
	climbRate = climbStep * 1000 / moveTick

	moveLimit     = 300000
	trialTime     = 15000
	chainSegment  = 120.0
	ropeSegment   = 90.0
	routeMapCount = 8

	attrNormal   = 0
	attrJumpDown = 15
)

type moveKind uint8

const (
	moveWalkOff moveKind = iota
	moveRunJump
	moveJump
	moveJumpDown
	moveGrab
	moveRopeTop
	moveRopeBottom
	moveRopeJump
	movePortal
)

// moveAction is one step of a route. It starts at x on a foothold chain, or at y on rope for the rope actions; a
// run jump runs from x and jumps at takeoff. dir is the side key, or the up or down key of a grab; catch holds up
// in the air to catch a rope. A portal action enters portal, standing under it.
type moveAction struct {
	kind    moveKind
	x, y    float64
	takeoff float64
	dir     int
	rope    wz.LadderRope
	catch   bool
	portal  innerPortal
}

// innerPortal leads from the portal named name at from to the portal of the same map at to.
type innerPortal struct {
	name     string
	from, to types.Point[int16]
}

// moveEdge leaves a route node by act and ends at at in node to (x on a chain, y on a rope) after cost ms.
type moveEdge struct {
	act  moveAction
	to   int
	at   float64
	cost float64
}

// routeNode is a stretch of a foothold chain from lo to hi in x, or of a rope (chain -1) from lo to hi in y.
type routeNode struct {
	chain  int
	rope   wz.LadderRope
	lo, hi float64
}

type footholdChain struct {
	footholds   []wz.Foothold
	left, right float64
}

// routeMaps keeps the route graphs of the maps bots walked last. Suites walk from their own actors, so the cache
// and each graph take a lock.
var routeMaps = &routeCache{order: list.New(), maps: make(map[routeKey]*list.Element)}

type routeCache struct {
	mu    sync.Mutex
	order *list.List
	maps  map[routeKey]*list.Element
}

// routeKey is a map walked at a speed and jump; where a jump lands differs with them, so each has its own graph.
type routeKey struct {
	mapID       uint32
	speed, jump int
}

type routeEntry struct {
	key   routeKey
	graph *routeMap
}

func (c *routeCache) find(key routeKey, m *wz.Map) *routeMap {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e, ok := c.maps[key]; ok {
		c.order.MoveToFront(e)
		return e.Value.(*routeEntry).graph
	}
	g := newRouteMap(key, m)
	c.maps[key] = c.order.PushFront(&routeEntry{key: key, graph: g})
	if c.order.Len() > routeMapCount {
		last := c.order.Back()
		c.order.Remove(last)
		delete(c.maps, last.Value.(*routeEntry).key)
	}
	return g
}

// routeMap is the route graph of one map. An edge is found the first time a search leaves its node, by trying the
// action on a walker, so a search holds mu.
type routeMap struct {
	mu         sync.Mutex
	tree       *types.QuadTreeNode[int16, wz.Foothold]
	footholds  map[int16]wz.Foothold
	ropes      []wz.LadderRope
	portals    []innerPortal
	fs         float64
	speed      float64
	jump       float64
	swim       bool
	waterCuts  map[int]bool
	chains     []footholdChain
	chainOf    map[int16]int
	nodes      []routeNode
	chainNodes [][]int
	ropeNodes  map[int16][]int
	links      map[int][]moveEdge
}

// newRouteMap builds the graph of m for a character of the speed and jump of key, which the client turns into
// factors of its walk, swim, and jump by a hundredth.
func newRouteMap(key routeKey, m *wz.Map) *routeMap {
	g := &routeMap{fs: 1, speed: float64(key.speed) / 100, jump: float64(key.jump) / 100, waterCuts: make(map[int]bool), footholds: make(map[int16]wz.Foothold), chainOf: make(map[int16]int), ropeNodes: make(map[int16][]int), links: make(map[int][]moveEdge)}
	id := key.mapID
	if m == nil || m.Footholds == nil {
		return g
	}
	g.tree, g.ropes, g.swim = m.Footholds, m.LadderRopes, m.Swim
	if m.FS > 0 {
		g.fs = m.FS
	}
	for _, from := range m.Portals {
		if from.TargetMapId != int32(id) || from.ScriptName != "" {
			continue
		}
		for _, to := range m.Portals {
			if to.Name == from.Target && to.ID != from.ID {
				g.portals = append(g.portals, innerPortal{name: from.Name, from: from.Position, to: to.Position})
				break
			}
		}
	}
	var all []wz.Foothold
	m.Footholds.QueryRange(m.Footholds.Bounds, &all)
	for _, fh := range all {
		g.footholds[fh.ID] = fh
	}
	for _, id := range slices.Sorted(maps.Keys(g.footholds)) {
		if fh := g.footholds[id]; fh.Wall || fh.X1 >= fh.X2 {
			continue
		}
		if _, ok := g.chainOf[id]; ok {
			continue
		}
		c := footholdChain{left: math.Inf(1), right: math.Inf(-1)}
		for _, at := range slices.Sorted(maps.Keys(g.chain(id))) {
			g.chainOf[at] = len(g.chains)
			c.footholds = append(c.footholds, g.footholds[at])
			c.left, c.right = min(c.left, float64(g.footholds[at].X1)), max(c.right, float64(g.footholds[at].X2))
		}
		g.chains = append(g.chains, c)
	}

	g.chainNodes = make([][]int, len(g.chains))
	for c, span := range g.chains {
		n := max(1, int(math.Ceil((span.right-span.left)/chainSegment)))
		width := (span.right - span.left) / float64(n)
		for i := range n {
			g.chainNodes[c] = append(g.chainNodes[c], len(g.nodes))
			g.nodes = append(g.nodes, routeNode{chain: c, lo: span.left + float64(i)*width, hi: span.left + float64(i+1)*width})
		}
	}
	for _, rope := range g.ropes {
		n := max(1, int(math.Ceil(float64(rope.Y2-rope.Y1)/ropeSegment)))
		height := float64(rope.Y2-rope.Y1) / float64(n)
		for i := range n {
			g.ropeNodes[rope.ID] = append(g.ropeNodes[rope.ID], len(g.nodes))
			g.nodes = append(g.nodes, routeNode{chain: -1, rope: rope, lo: float64(rope.Y1) + float64(i)*height, hi: float64(rope.Y1) + float64(i+1)*height})
		}
	}
	return g
}

// edges are the actions that leave node n for another chain or rope, each kept only when a trial shows where it
// ends.
func (g *routeMap) edges(n int) []moveEdge {
	if links, ok := g.links[n]; ok {
		return links
	}
	var links []moveEdge
	for _, act := range g.actions(n) {
		if e, ok := g.try(n, act); ok {
			links = append(links, e)
		}
	}
	g.links[n] = links
	return links
}

// actions are the moves worth trying from node n: on a rope, stepping off its top or bottom and jumping off to
// either side; on a chain, jumps in place, running jumps, a jump down, walking off its ends, and grabbing or
// catching nearby ropes from several distances.
func (g *routeMap) actions(n int) []moveAction {
	node := g.nodes[n]
	var acts []moveAction
	if node.chain < 0 {
		rope, y := node.rope, (node.lo+node.hi)/2
		if node.lo <= float64(rope.Y1) && rope.UpperFloor {
			acts = append(acts, moveAction{kind: moveRopeTop, y: y, rope: rope})
		}
		if node.hi >= float64(rope.Y2) {
			acts = append(acts, moveAction{kind: moveRopeBottom, y: y, rope: rope})
		}
		for _, dir := range []int{-1, 1} {
			for _, catch := range []bool{false, true} {
				acts = append(acts, moveAction{kind: moveRopeJump, y: y, dir: dir, rope: rope, catch: catch})
			}
		}
		return acts
	}

	span := g.chains[node.chain]
	x := (node.lo + node.hi) / 2
	s, _ := g.surface(node.chain, x)
	reach := g.reach()
	var near []wz.LadderRope
	for _, rope := range g.ropes {
		if math.Abs(float64(rope.X)-x) <= chainSegment+90 && float64(rope.Y2) >= s-reach-60 && float64(rope.Y1) <= s+50 {
			near = append(near, rope)
		}
	}
	// A walk stops a few pixels off its mark, so no action starts at the very end of a chain.
	inside := func(x float64) float64 {
		if span.right-span.left < 12 {
			return (span.left + span.right) / 2
		}
		return min(max(x, span.left+6), span.right-6)
	}
	jump := func(x float64, dir int, catch bool) {
		x = inside(x)
		acts = append(acts, moveAction{kind: moveJump, x: x, dir: dir, catch: catch})
		if run := x - float64(dir)*40; dir != 0 && run >= span.left+6 && run <= span.right-6 {
			acts = append(acts, moveAction{kind: moveRunJump, x: run, takeoff: x, dir: dir, catch: catch})
		}
	}
	for _, dir := range []int{-1, 0, 1} {
		jump(x, dir, false)
		if len(near) > 0 {
			jump(x, dir, true)
		}
	}
	jumpDown := func(x float64) {
		x = inside(x)
		fh, ok := g.footholdAt(node.chain, x)
		if !ok || fh.ForbidFallDown {
			return
		}
		if _, _, ok := g.below(x, float64(fh.SurfaceY(int16(x)))+1, fh.ID); ok {
			acts = append(acts, moveAction{kind: moveJumpDown, x: x})
		}
	}
	jumpDown(x)
	for d, other := range g.chains {
		if d == node.chain {
			continue
		}
		if lo, hi := max(node.lo, other.left), min(node.hi, other.right); lo <= hi {
			mid := (lo + hi) / 2
			here, _ := g.surface(node.chain, mid)
			there, _ := g.surface(d, mid)
			switch {
			case there < here-8 && there >= here-reach:
				jump(mid, 0, false)
			case there > here+8 && math.Abs(mid-x) > 10:
				jumpDown(mid)
			}
			continue
		}
		from, dir, to := node.hi, 1, other.left
		if other.right < node.lo {
			from, dir, to = node.lo, -1, other.right
		}
		here, _ := g.surface(node.chain, from)
		there, _ := g.surface(d, to)
		if math.Abs(to-from) <= 80*g.speed && there >= here-reach && there <= here+200 {
			jump(from, dir, false)
		}
	}
	for _, dir := range []int{-1, 1} {
		end := span.left
		if dir > 0 {
			end = span.right
		}
		if end < node.lo || end > node.hi {
			continue
		}
		takeoff := end - float64(dir)*2
		acts = append(acts, moveAction{kind: moveWalkOff, x: inside(end - float64(dir)*40), dir: dir})
		for _, catch := range []bool{false, true} {
			jump(takeoff, dir, catch)
		}
	}
	for _, rope := range near {
		rx := float64(rope.X)
		if rx >= node.lo && rx <= node.hi {
			if rs, ok := g.surface(node.chain, rx); ok {
				if float64(rope.Y1) < rs-8 && float64(rope.Y2) >= rs-reach {
					acts = append(acts, moveAction{kind: moveGrab, x: rx, dir: -1, rope: rope})
				}
				if float64(rope.Y1) >= rs && float64(rope.Y1) <= rs+10 {
					acts = append(acts, moveAction{kind: moveGrab, x: rx, dir: 1, rope: rope})
				}
			}
		}
		for _, off := range []float64{30, 50, 70, 90} {
			for _, side := range []int{-1, 1} {
				if from := rx + float64(side)*off; from >= node.lo && from <= node.hi {
					jump(from, -side, true)
				}
			}
		}
	}
	// The server takes a portal within 150 px of it; the bot uses one from the floor right under it.
	for _, portal := range g.portals {
		px := float64(portal.from.X)
		if px < node.lo || px > node.hi {
			continue
		}
		if floor, ok := g.surface(node.chain, inside(px)); ok && floor-float64(portal.from.Y) >= -20 && floor-float64(portal.from.Y) <= 100 {
			acts = append(acts, moveAction{kind: movePortal, x: inside(px), portal: portal})
		}
	}
	return acts
}

// try carries out act on walkers standing at its start and a few pixels to either side, since the bot stops a
// little off where it walks to; the edge holds only when all of them end on the same chain or rope.
func (g *routeMap) try(n int, act moveAction) (moveEdge, bool) {
	node := g.nodes[n]
	var edge moveEdge
	var place routeNode
	for i, off := range []float64{0, -3, 3} {
		trial := &walker{g: g, limit: trialTime}
		at := act
		if node.chain < 0 {
			at.y = min(max(at.y+off, float64(node.rope.Y1)), float64(node.rope.Y2))
			trial.y = at.y
			trial.hold(node.rope)
		} else {
			span := g.chains[node.chain]
			at.x = min(max(at.x+off, span.left), span.right)
			at.takeoff = min(max(at.takeoff+off, span.left), span.right)
			fh, ok := g.footholdAt(node.chain, at.x)
			if !ok {
				return moveEdge{}, false
			}
			trial.x = at.x
			trial.stand(fh)
		}
		if !trial.do(at) {
			return moveEdge{}, false
		}
		to := trial.node()
		if to < 0 {
			return moveEdge{}, false
		}
		dest := g.nodes[to]
		if dest.chain == node.chain && (dest.chain >= 0 || dest.rope.ID == node.rope.ID) {
			return moveEdge{}, false
		}
		if i > 0 {
			if dest.chain != place.chain || dest.rope.ID != place.rope.ID {
				return moveEdge{}, false
			}
			continue
		}
		place = dest
		edge = moveEdge{act: act, to: to, at: trial.x, cost: float64(trial.spent)}
		if dest.chain < 0 {
			edge.at = trial.y
		}
	}
	return edge, true
}

func (g *routeMap) siblings(n int) []int {
	node := g.nodes[n]
	if node.chain < 0 {
		return g.ropeNodes[node.rope.ID]
	}
	return g.chainNodes[node.chain]
}

func (g *routeMap) nodeAt(chain int, x float64) int {
	nodes := g.chainNodes[chain]
	for _, n := range nodes {
		if x <= g.nodes[n].hi {
			return n
		}
	}
	return nodes[len(nodes)-1]
}

func (g *routeMap) ropeNodeAt(id int16, y float64) int {
	nodes := g.ropeNodes[id]
	for _, n := range nodes {
		if y <= g.nodes[n].hi {
			return n
		}
	}
	return nodes[len(nodes)-1]
}

// estimate is a lower bound of the ms from at in node n to the target: no move crosses faster than a jump off a
// rope sideways, a fall down, or a jump up.
func (g *routeMap) estimate(n int, at, tx, ty float64) float64 {
	node := g.nodes[n]
	x, y := float64(node.rope.X), at
	if node.chain >= 0 {
		x = at
		y, _ = g.surface(node.chain, at)
	}
	t := math.Abs(tx-x) / (ropeJumpBoost * g.speed)
	if dy := ty - y; dy > 0 {
		t = max(t, dy/fallSpeed)
	} else {
		t = max(t, -dy/(jumpSpeed*g.jump))
	}
	return t * 1000
}

func (g *routeMap) query(left, top, right, bottom float64) []wz.Foothold {
	var found []wz.Foothold
	if g.tree != nil {
		g.tree.QueryRange(types.Rect[int16]{Left: int16(math.Floor(left)), Top: int16(math.Floor(top)), Right: int16(math.Ceil(right)), Bottom: int16(math.Ceil(bottom))}, &found)
	}
	return found
}

// below finds the foothold a fall at x from y lands on, passing through skip as a jump down does.
func (g *routeMap) below(x, y float64, skip int16) (wz.Foothold, float64, bool) {
	var found wz.Foothold
	top := math.Inf(1)
	if g.tree == nil {
		return found, top, false
	}
	for _, fh := range g.query(x, y, x, float64(g.tree.Bounds.Bottom)) {
		if fh.Wall || fh.X1 > fh.X2 || fh.ID == skip || x < float64(fh.X1) || x > float64(fh.X2) {
			continue
		}
		if y2 := float64(fh.SurfaceY(int16(x))); y2 >= y && y2 < top {
			found, top = fh, y2
		}
	}
	return found, top, found.ID != 0
}

// findRope finds the first rope within 10 px of the x range whose span meets the y range, as the client searches.
func (g *routeMap) findRope(x1, y1, x2, y2 float64) (wz.LadderRope, bool) {
	for _, rope := range g.ropes {
		if float64(rope.X) >= min(x1, x2)-10 && float64(rope.X) <= max(x1, x2)+10 && min(y1, y2) <= float64(rope.Y2) && float64(rope.Y1) <= max(y1, y2) {
			return rope, true
		}
	}
	return wz.LadderRope{}, false
}

func (g *routeMap) findRopeByID(id int16) wz.LadderRope {
	for _, rope := range g.ropes {
		if rope.ID == id {
			return rope
		}
	}
	return wz.LadderRope{}
}

func (g *routeMap) footholdAt(c int, x float64) (wz.Foothold, bool) {
	for _, fh := range g.chains[c].footholds {
		if x >= float64(fh.X1) && x <= float64(fh.X2) {
			return fh, true
		}
	}
	return wz.Foothold{}, false
}

func (g *routeMap) surface(c int, x float64) (float64, bool) {
	fh, ok := g.footholdAt(c, x)
	if !ok {
		return math.Inf(1), false
	}
	return float64(fh.SurfaceY(int16(x))), true
}

// chain is the run of footholds a walk follows from id both ways, without a wall or a gap. A link that holds only
// one way, as some maps have, ends the chain: the walk across it goes one way only.
func (g *routeMap) chain(id int16) map[int16]bool {
	chain := map[int16]bool{id: true}
	for _, dir := range []int{-1, 1} {
		for at := g.footholds[id]; ; {
			next, ok := g.joined(at, dir)
			if !ok || next.Wall || chain[next.ID] {
				break
			}
			if back, ok := g.joined(next, -dir); !ok || back.ID != at.ID {
				break
			}
			chain[next.ID] = true
			at = next
		}
	}
	return chain
}

// crosses reports a move from (x0, y0) to (x1, y1) passing through fh from its blocking side, as the client tests a
// fall: from above for ground, from below for a ceiling, and from one side for a wall, by the way it runs.
func (g *routeMap) crosses(fh wz.Foothold, x0, y0, x1, y1 float64) bool {
	ax, ay, bx, by := float64(fh.X1), float64(fh.Y1), float64(fh.X2), float64(fh.Y2)
	from := (bx-ax)*(y0-ay) - (by-ay)*(x0-ax)
	to := (bx-ax)*(y1-ay) - (by-ay)*(x1-ax)
	if from > 0 || to < 0 || from == 0 && to == 0 {
		return false
	}
	return ((x1-x0)*(ay-y0)-(y1-y0)*(ax-x0))*((x1-x0)*(by-y0)-(y1-y0)*(bx-x0)) <= 0
}

// joined is the foothold fh links to on the dir side when it starts where fh ends there; a wall counts.
func (g *routeMap) joined(fh wz.Foothold, dir int) (wz.Foothold, bool) {
	id, end := fh.Next, fh.X2
	if dir < 0 {
		id, end = fh.Prev, fh.X1
	}
	next, ok := g.footholds[id]
	switch {
	case !ok:
		return next, false
	case next.Wall:
		return next, next.X1 == end
	case dir > 0:
		return next, next.X1 < next.X2 && next.X1 == end
	default:
		return next, next.X1 < next.X2 && next.X2 == end
	}
}

type pathFragment struct {
	attr         uint8
	x, y, vx, vy int16
	fh, fallFrom int16
	stance       uint8
	duration     int16
}

// walker steps a character like the client's vector controller and keeps the move packets the client would send.
type walker struct {
	g        *routeMap
	x, y     float64
	vx, vy   float64
	speed    float64
	fh       int16
	fallFrom int16
	zmass    int16
	rope     int16
	left     bool
	spent    int
	limit    int
	start    types.Point[int16]
	path     []pathFragment
	elapsed  int
	keys     []uint8
	packets  []outbound
}

func newWalker(g *routeMap, from types.Point[int16]) *walker {
	w := &walker{g: g, x: float64(from.X), y: float64(from.Y), limit: moveLimit, start: from}
	if fh, top, ok := g.below(w.x, w.y-5, 0); ok && top-w.y <= 5 {
		w.stand(fh)
	} else if rope, ok := g.findRope(w.x, w.y, w.x, w.y); ok && float64(rope.X) == w.x {
		w.hold(rope)
	}
	return w
}

// moveTo takes the bot to target along the cheapest route of walks, jumps, drops, climbs, and portals of the map,
// searching again after each action, then lets the client's last half second of path go out. A target up to 60 px
// over a floor is reached on that floor, like an NPC or a portal drawn above it; one higher counts only on a rope,
// where the bot stops holding it.
func (w *walker) moveTo(target types.Point[int16]) bool {
	tx, ty := float64(target.X), float64(target.Y)
	var goal int
	if fh, top, ok := w.g.below(tx, ty-20, 0); ok && top-ty <= 60 {
		goal = w.g.nodeAt(w.g.chainOf[fh.ID], tx)
	} else if rope, ok := w.g.findRope(tx, ty, tx, ty); ok {
		goal = w.g.ropeNodeAt(rope.ID, ty)
	} else {
		return false
	}
	if w.g.swim {
		return w.swimTo(goal, tx, ty)
	}
	end := w.g.nodes[goal]

	// An action that keeps landing somewhere its trials did not, as on conveyors where the bot never stands still,
	// is picked again from the same node; the bot gives up on it the third time.
	type attempt struct {
		node int
		act  moveAction
	}
	tried := make(map[attempt]int)
	for range 100 {
		if !w.land() {
			return false
		}
		switch here := w.g.nodes[w.node()]; {
		case end.chain >= 0 && here.chain == end.chain:
			return w.walkTo(tx) && w.settle()
		case end.chain < 0 && w.rope == end.rope.ID:
			return w.climbTo(ty) && w.settle()
		}
		act, ok := w.plan(goal, tx, ty)
		if !ok {
			return false
		}
		key := attempt{node: w.node(), act: act}
		if tried[key]++; tried[key] > 2 || !w.do(act) {
			return false
		}
	}
	return false
}

// plan searches the route graph by A* from where the bot is to goal and returns the first action of the route;
// moving along the same chain or rope between nodes needs no action.
func (w *walker) plan(goal int, tx, ty float64) (moveAction, bool) {
	type label struct {
		cost, at float64
		first    moveAction
		acted    bool
	}
	start := w.node()
	labels := map[int]*label{start: {at: w.x}}
	if w.rope != 0 {
		labels[start].at = w.y
	}
	done := make(map[int]bool)
	for {
		cur, best := -1, math.Inf(1)
		for n, l := range labels {
			if f := l.cost + w.g.estimate(n, l.at, tx, ty); !done[n] && f < best {
				cur, best = n, f
			}
		}
		if cur < 0 {
			return moveAction{}, false
		}
		l := labels[cur]
		if cur == goal {
			return l.first, l.acted
		}
		done[cur] = true

		rate := walkSpeed * w.g.speed
		if w.g.nodes[cur].chain < 0 {
			rate = climbRate
		}
		relax := func(to int, at, cost float64, act *moveAction) {
			if old, ok := labels[to]; done[to] || ok && old.cost <= l.cost+cost {
				return
			}
			next := &label{cost: l.cost + cost, at: at, first: l.first, acted: l.acted}
			if act != nil && !l.acted {
				next.first, next.acted = *act, true
			}
			labels[to] = next
		}
		for _, n := range w.g.siblings(cur) {
			at := min(max(l.at, w.g.nodes[n].lo), w.g.nodes[n].hi)
			relax(n, at, math.Abs(at-l.at)/rate*1000, nil)
		}
		for _, e := range w.g.edges(cur) {
			from := e.act.x
			if w.g.nodes[cur].chain < 0 {
				from = e.act.y
			}
			relax(e.to, e.at, math.Abs(from-l.at)/rate*1000+e.cost, &e.act)
		}
	}
}

// do walks or climbs to where act starts and carries it out, returning once the bot stands on a foothold again or
// holds a rope.
func (w *walker) do(act moveAction) bool {
	switch act.kind {
	case moveRopeTop, moveRopeBottom, moveRopeJump:
		if w.rope != act.rope.ID || !w.climbTo(act.y) {
			return false
		}
	default:
		if !w.walkTo(act.x) {
			return false
		}
	}

	air, grab := 0, 0
	if act.catch {
		grab = -1
	}
	switch act.kind {
	case moveWalkOff:
		air = act.dir
		for w.fh != 0 {
			if !w.step(act.dir, 0, false) {
				return false
			}
		}
	case moveRunJump:
		air = act.dir
		for w.fh != 0 && float64(act.dir)*(act.takeoff-w.x) > 0 {
			if !w.step(act.dir, 0, false) {
				return false
			}
		}
		if w.fh != 0 && !w.step(act.dir, grab, true) {
			return false
		}
	case moveJump:
		air = act.dir
		if !w.step(act.dir, grab, true) {
			return false
		}
	case moveJumpDown:
		if !w.step(0, 1, true) {
			return false
		}
	case moveGrab:
		jump := false
		if act.dir < 0 {
			grab, jump = -1, float64(act.rope.Y2) < w.y-20 || float64(act.rope.Y1) > w.y
		}
		if !w.step(0, act.dir, jump) {
			return false
		}
	case moveRopeTop:
		if !w.climbTo(math.Inf(-1)) {
			return false
		}
	case movePortal:
		// The client sends the portal while its path may still be held; the bot sends the path first so the
		// server's 150 px check sees where it stands. The client lands at the target 120 ms later, 10 px above it.
		if len(w.path) > 0 {
			w.flush()
		}
		w.packets = append(w.packets, &request.UseInnerPortal{PortalName: act.portal.name, FromX: int16(w.x), FromY: int16(w.y), ToX: act.portal.to.X, ToY: act.portal.to.Y})
		for range 4 {
			if !w.step(0, 0, false) {
				return false
			}
		}
		w.x, w.y = float64(act.portal.to.X), float64(act.portal.to.Y)-10
		w.fh, w.fallFrom, w.vx, w.vy, w.speed = 0, 0, 0, 0, 0
	case moveRopeBottom:
		if !w.climbTo(math.Inf(1)) {
			return false
		}
	case moveRopeJump:
		air = act.dir
		if !w.step(act.dir, 0, true) {
			return false
		}
	}
	for w.fh == 0 && w.rope == 0 {
		if !w.step(air, grab, false) {
			return false
		}
	}
	return true
}

func (w *walker) node() int {
	switch {
	case w.rope != 0:
		return w.g.ropeNodeAt(w.rope, w.y)
	case w.fh != 0:
		return w.g.nodeAt(w.g.chainOf[w.fh], w.x)
	}
	return -1
}

// walkTo walks along the chain to x and stops there. A conveyor never lets the bot stand still, so on one it only
// needs to be at x.
func (w *walker) walkTo(x float64) bool {
	for math.Abs(x-w.x) > 2 || w.speed != 0 && w.g.footholds[w.fh].Force == 0 {
		if w.fh == 0 {
			return false
		}
		dir := cmp.Compare(x-w.x, 0)
		stop := w.speed * w.speed / (2 * w.g.walkDrag() / shoeMass)
		if math.Abs(x-w.x) <= 2 || cmp.Compare(w.speed, 0) == dir && math.Abs(x-w.x) <= stop {
			dir = 0
		}
		if !w.step(dir, 0, false) {
			return false
		}
	}
	return true
}

// climbTo climbs the held rope toward y and returns once the bot is there or has left the rope at an end.
func (w *walker) climbTo(y float64) bool {
	for w.rope != 0 && math.Abs(w.y-y) > climbStep/2 {
		if !w.step(0, cmp.Compare(y, w.y), false) {
			return false
		}
	}
	return true
}

func (w *walker) land() bool {
	for w.fh == 0 && w.rope == 0 {
		if !w.step(0, 0, false) {
			return false
		}
	}
	return true
}

// settle lets the bot stand until the path it walked has gone out.
func (w *walker) settle() bool {
	for len(w.path) > 0 || len(w.packets) == 0 {
		w.step(0, 0, false)
	}
	return true
}

// step runs one tick with the arrow keys x and y and the jump key held; it reports false once the walk runs out of
// time or falls out of the map.
func (w *walker) step(x, y int, jump bool) bool {
	px, py := w.x, w.y
	switch {
	case w.rope != 0:
		w.climb(x, y, jump)
	case jump && y > 0 && w.fh != 0:
		w.fallFrom, w.fh = w.fh, 0
		w.vx, w.vy = 0, jumpDownSpeed
		w.fall(x, y)
	case jump && w.fh != 0:
		w.jump(x)
		w.fall(x, y)
	case jump && w.g.swim:
		w.vy = -swimSpeed * w.g.speed * swimStroke
		w.fall(x, y)
	case w.fh != 0:
		w.walk(x)
	default:
		w.fall(x, y)
	}
	if x != 0 {
		w.left = x < 0
	}
	w.record(x, y)
	if w.rope == 0 && y != 0 {
		w.grab(y, px, py)
	}

	w.spent += moveTick
	return w.spent <= w.limit && (w.g.tree == nil || w.y <= float64(w.g.tree.Bounds.Bottom)+100)
}

func (w *walker) jump(input int) {
	w.fh = 0
	w.vy = -jumpSpeed * w.g.jump
	if w.g.swim {
		w.vy *= swimJump
	}
	if input == 0 {
		return
	}
	dir, limit := float64(input), walkSpeed*w.g.speed
	if limit*jumpBoost > dir*w.vx {
		w.vx += dir * limit * jumpBoost
	}
	if dir*w.vx > limit {
		w.vx = dir * limit
	}
}

// climb moves along the held rope. The bot steps off the top only onto an upper floor and lets go below the
// bottom; a jump off the rope needs a side key and moves the bot only from the next tick.
func (w *walker) climb(x, y int, jump bool) {
	rope := w.g.findRopeByID(w.rope)
	if jump && x != 0 {
		w.rope = 0
		w.vx, w.vy = float64(x)*ropeJumpBoost*w.g.speed, -ropeJumpSpeed*w.g.jump
		return
	}
	w.y += float64(y) * climbStep
	switch {
	case y < 0 && w.y < float64(rope.Y1) && rope.UpperFloor:
		w.y, w.rope = float64(rope.Y1)-5, 0
	case y < 0 && w.y < float64(rope.Y1):
		w.y = float64(rope.Y1)
	case y > 0 && w.y > float64(rope.Y2):
		w.y, w.rope = float64(rope.Y2)+1, 0
	}
}

// grab catches a rope the way the client does after a tick: up while falling catches one the fall passed, up on a
// foothold one overhead, and down on a foothold one whose top is at the feet.
func (w *walker) grab(y int, px, py float64) {
	switch {
	case w.fh == 0 && y < 0 && w.vy > 0:
		if rope, ok := w.g.findRope(px, py, w.x, w.y); ok {
			w.hold(rope)
		}
	case w.fh != 0 && y < 0:
		if rope, ok := w.g.findRope(w.x, w.y-20, w.x, w.y); ok {
			w.y = min(float64(rope.Y2), w.y)
			w.hold(rope)
		}
	case w.fh != 0 && y > 0:
		if rope, ok := w.g.findRope(w.x, w.y, w.x, w.y+10); ok && w.y <= float64(rope.Y1) && float64(rope.Y1) <= w.y+10 {
			w.y = float64(rope.Y1)
			w.hold(rope)
		}
	}
}

func (w *walker) hold(rope wz.LadderRope) {
	w.rope, w.fh, w.fallFrom = rope.ID, 0, 0
	w.x, w.vx, w.vy, w.speed = float64(rope.X), 0, 0, 0
}

// walk moves along the foothold chain; the bot stops at a wall and leaves the chain where it ends. The map's fs
// scales the push and the friction, water slows the walk, and a conveyor pushes the bot, faster along it and slower
// against it by the whole part of its force.
func (w *walker) walk(input int) {
	const t = moveTick / 1000.0
	fh := w.g.footholds[w.fh]
	force, limit := walkForce*w.g.fs, walkSpeed*w.g.speed
	if w.g.swim {
		force, limit = force*swimSpeedDec, limit*swimSpeedDec
	}
	push := float64(input) * force
	if conveyor := math.Abs(math.Trunc(fh.Force)); conveyor > 0 {
		switch {
		case input == 0:
			push, limit = fh.Force*force, limit*conveyor
		case float64(input)*fh.Force > 0:
			push, limit = push*2*conveyor, limit*2*conveyor
		default:
			push, limit = push*conveyorAgainst/conveyor, limit*conveyorAgainst/conveyor
		}
	}
	drag := w.g.walkDrag() / shoeMass * t
	speed := w.speed
	switch {
	case speed > limit:
		speed = max(speed-drag, limit)
	case speed < -limit:
		speed = min(speed+drag, -limit)
	}
	switch {
	case push > 0 && speed < limit:
		speed = min(speed+push/shoeMass*t, limit)
	case push < 0 && speed > -limit:
		speed = max(speed+push/shoeMass*t, -limit)
	case push == 0 && speed > 0:
		speed = max(speed-drag, 0)
	case push == 0 && speed < 0:
		speed = min(speed+drag, 0)
	}
	cos := float64(fh.X2-fh.X1) / math.Hypot(float64(fh.X2-fh.X1), float64(fh.Y2-fh.Y1))
	w.x += (w.speed + speed) / 2 * t * cos
	w.speed = speed

	for w.x > float64(fh.X2) || w.x < float64(fh.X1) {
		dir := cmp.Compare(w.x, float64(fh.X2))
		end, endY := fh.X2, fh.Y2
		if w.x < float64(fh.X1) {
			end, endY = fh.X1, fh.Y1
		}
		next, ok := w.g.joined(fh, dir)
		switch {
		case ok && next.Wall && min(next.Y1, next.Y2) < endY:
			w.x, w.speed = float64(end), 0
		case ok && !next.Wall:
			fh = next
		default:
			w.stand(fh)
			w.fh, w.y = 0, float64(endY)
			return
		}
	}
	w.stand(fh)
}

// fall moves through the air, or through water on a swim map. The walls and ceilings of the group the bot last
// stood on stop it, and a foothold under the drop catches it, except the one it jumped down through.
func (w *walker) fall(input, vertical int) {
	const t = moveTick / 1000.0
	vx, vy := w.vx, w.vy
	switch {
	case w.g.swim:
		vx, vy = w.swim(input, vertical)
	case input != 0:
		limit := walkSpeed / walkForce * floatDrag2 * w.g.fs
		force := 2 * floatDrag2 * w.g.fs / shoeMass * t
		if input > 0 && vx < limit {
			vx = min(vx+force, limit)
		}
		if input < 0 && vx > -limit {
			vx = max(vx-force, -limit)
		}
	default:
		drag := floatDrag2 * w.g.fs / shoeMass * t
		if vy < fallSpeed {
			drag *= floatCoefficient
		}
		if vx > 0 {
			vx = max(vx-drag, 0)
		} else {
			vx = min(vx+drag, 0)
		}
	}
	if !w.g.swim && vy < fallSpeed {
		vy = min(vy+gravityAcc*t, fallSpeed)
	}
	x := w.x + (w.vx+vx)/2*t
	y := w.y + (w.vy+vy)/2*t
	for _, fh := range w.g.query(min(w.x, x), min(w.y, y), max(w.x, x), max(w.y, y)) {
		if fh.X1 < fh.X2 || fh.ZMass != w.zmass || !w.g.crosses(fh, w.x, w.y, x, y) {
			continue
		}
		if fh.Wall {
			x, vx = float64(fh.X1)+math.Copysign(0.5, w.x-float64(fh.X1)), 0
		} else {
			y, vy = max(y, float64(fh.SurfaceY(int16(x)))+0.5), 0
		}
	}
	w.vx, w.vy = vx, vy

	if vy > 0 {
		if fh, top, ok := w.g.below(x, w.y, w.fallFrom); ok && top <= y {
			w.x, w.speed = x, vx
			w.stand(fh)
			return
		}
	}
	w.x, w.y = x, y
}

// swim is the velocity after a tick in water. Side keys swim up to the swim speed and water drags anything past it;
// the bot always sinks, slower with up held and faster with down, and rises only by the stroke of a jump.
func (w *walker) swim(input, vertical int) (float64, float64) {
	const t = moveTick / 1000.0
	drag := floatDrag1 * w.g.fs / shoeMass * t
	force := swimForce / shoeMass * t
	limit := swimSpeed * w.g.speed
	vx, vy := w.vx, w.vy
	switch {
	case vx > limit:
		vx = max(vx-drag, limit)
	case vx < -limit:
		vx = min(vx+drag, -limit)
	}
	switch {
	case input > 0 && vx < limit:
		vx = min(vx+force, limit)
	case input < 0 && vx > -limit:
		vx = max(vx-force, -limit)
	case input == 0 && vx > 0:
		vx = max(vx-drag, 0)
	case input == 0 && vx < 0:
		vx = min(vx+drag, 0)
	}

	switch {
	case vy > limit:
		vy = max(vy-drag, limit)
	case vy < -limit:
		vy = min(vy+drag, -limit)
	}
	sink, rise := limit, force
	switch {
	case vertical < 0:
		sink, rise = limit*swimHover, force/2
	case vertical > 0:
		sink, rise = limit*swimDive, force/2
	}
	if vy < sink {
		vy = min(vy+rise, sink)
	} else {
		vy = max(vy-force, sink)
	}
	return vx, vy
}

// reach is how high a jump lifts the bot, with room for the trials to decide the last pixels.
func (g *routeMap) reach() float64 {
	rise := jumpSpeed * g.jump
	return rise*rise/(2*gravityAcc) + 13
}

// walkDrag is the friction of walking: the map's fs within the physics limits, halved when below 1.
func (g *routeMap) walkDrag() float64 {
	friction := min(max(g.fs, minFriction), maxFriction)
	if friction < 1 {
		friction /= 2
	}
	return friction * walkDrag
}

// stand puts the bot on fh at its x, turning its speed along the foothold into the velocity the client reports.
func (w *walker) stand(fh wz.Foothold) {
	length := math.Hypot(float64(fh.X2-fh.X1), float64(fh.Y2-fh.Y1))
	w.fh, w.fallFrom, w.zmass = fh.ID, 0, fh.ZMass
	w.y = float64(fh.SurfaceY(int16(w.x)))
	w.vx = w.speed * float64(fh.X2-fh.X1) / length
	w.vy = w.speed * float64(fh.Y2-fh.Y1) / length
}

// record adds the tick to the path as the client does: it stretches the last fragment while the attribute,
// foothold, stance, and velocity signs hold, and sends the path once half a second of it gathers. On a rope the
// fragment's foothold is the rope's id negated.
func (w *walker) record(x, y int) {
	var stance uint8
	switch rope := w.g.findRopeByID(w.rope); {
	case w.rope != 0 && rope.Ladder:
		stance = constant.StanceLadderRight
	case w.rope != 0:
		stance = constant.StanceRopeRight
	case w.fh != 0 && x != 0:
		stance = constant.StanceWalkRight
	case w.fh != 0:
		stance = constant.StanceStandingRight
	default:
		stance = constant.StanceJumpRight
	}
	if w.left {
		stance++
	}
	attr := uint8(attrNormal)
	if w.fallFrom != 0 {
		attr = attrJumpDown
	}
	fh := w.fh
	if w.rope != 0 {
		fh = -w.rope
	}
	f := pathFragment{attr: attr, x: int16(w.x), y: int16(w.y), vx: int16(w.vx), vy: int16(w.vy), fh: fh, fallFrom: w.fallFrom, stance: stance, duration: moveTick}
	if n := len(w.path); n > 0 {
		last := w.path[n-1]
		if last.attr == f.attr && last.fh == f.fh && last.fallFrom == f.fallFrom && last.stance == f.stance &&
			cmp.Compare(last.vx, 0) == cmp.Compare(f.vx, 0) && cmp.Compare(last.vy, 0) == cmp.Compare(f.vy, 0) {
			f.duration += last.duration
			w.path = w.path[:n-1]
		}
	}
	w.path = append(w.path, f)

	var key uint8
	if x < 0 {
		key |= 8
	}
	if x > 0 {
		key |= 4
	}
	if y < 0 {
		key |= 2
	}
	if y > 0 {
		key |= 1
	}
	w.keys = append(w.keys, key)

	w.elapsed += moveTick
	if w.elapsed >= moveFlush {
		w.flush()
	}
}

// flush sends the gathered path as one move packet starting where the last one ended.
func (w *walker) flush() {
	pkt := &request.MovePlayer{Position: w.start, KeyStates: w.keys, Bounds: types.Rect[int16]{Left: w.start.X, Top: w.start.Y, Right: w.start.X, Bottom: w.start.Y}}
	for _, f := range w.path {
		basic := &dto.BasicMovement{Command: f.attr, Stance: f.stance}
		pos := types.Vector2[int16]{X: f.x, Y: f.y}
		vel := types.Vector2[int16]{X: f.vx, Y: f.vy}
		if f.attr == attrJumpDown {
			pkt.Fragments = append(pkt.Fragments, &dto.JumpDownMovement{BasicMovement: basic, Position: pos, Velocity: vel, Foothold: f.fh, FallFoothold: f.fallFrom, Duration: f.duration})
		} else {
			pkt.Fragments = append(pkt.Fragments, &dto.AbsoluteLifeMovement{BasicMovement: basic, Position: pos, Velocity: vel, Foothold: f.fh, Duration: f.duration})
		}
		pkt.Bounds.Left, pkt.Bounds.Right = min(pkt.Bounds.Left, f.x), max(pkt.Bounds.Right, f.x)
		pkt.Bounds.Top, pkt.Bounds.Bottom = min(pkt.Bounds.Top, f.y), max(pkt.Bounds.Bottom, f.y)
	}
	w.packets = append(w.packets, pkt)
	last := w.path[len(w.path)-1]
	w.start = types.Point[int16]{X: last.x, Y: last.y}
	w.path, w.keys, w.elapsed = nil, nil, 0
}
