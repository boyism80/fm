package game

import internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"

// FamilyMember is one person in a family. The tree copies it from internal and does not keep that message.
type FamilyMember struct {
	characterID     uint32
	seniorID        uint32
	name            string
	level           uint32
	classID         uint32
	reputation      uint32
	totalReputation uint32
	channelID       int32
	buffUntil       int64
	buffFeature     internal.FamilyFeature
	usedFeatures    []internal.FamilyFeature
}

func (m *FamilyMember) CharacterID() uint32 {
	if m == nil {
		return 0
	}
	return m.characterID
}

func (m *FamilyMember) SeniorID() uint32 {
	if m == nil {
		return 0
	}
	return m.seniorID
}

func (m *FamilyMember) Name() string {
	if m == nil {
		return ""
	}
	return m.name
}

func (m *FamilyMember) Level() uint32 {
	if m == nil {
		return 0
	}
	return m.level
}

func (m *FamilyMember) ClassID() uint32 {
	if m == nil {
		return 0
	}
	return m.classID
}

func (m *FamilyMember) Reputation() uint32 {
	if m == nil {
		return 0
	}
	return m.reputation
}

func (m *FamilyMember) TotalReputation() uint32 {
	if m == nil {
		return 0
	}
	return m.totalReputation
}

func (m *FamilyMember) ChannelID() int32 {
	if m == nil {
		return 0
	}
	return m.channelID
}

func (m *FamilyMember) BuffUntil() int64 {
	if m == nil {
		return 0
	}
	return m.buffUntil
}

func (m *FamilyMember) BuffFeature() internal.FamilyFeature {
	if m == nil {
		return 0
	}
	return m.buffFeature
}

func (m *FamilyMember) UsedFeatures() []internal.FamilyFeature {
	if m == nil {
		return nil
	}
	return m.usedFeatures
}

func familyMemberFromProto(pb *internal.FamilyMember) *FamilyMember {
	if pb == nil {
		return nil
	}
	return &FamilyMember{
		characterID:     pb.GetCharacterId(),
		seniorID:        pb.GetSeniorId(),
		name:            pb.GetName(),
		level:           pb.GetLevel(),
		classID:         pb.GetClassId(),
		reputation:      pb.GetReputation(),
		totalReputation: pb.GetTotalReputation(),
		channelID:       pb.GetChannelId(),
		buffUntil:       pb.GetBuffUntil(),
		buffFeature:     pb.GetBuffFeature(),
		usedFeatures:    append([]internal.FamilyFeature(nil), pb.GetUsedFeatures()...),
	}
}

// FamilyTree is one family. Members are copied when the tree is built.
type FamilyTree struct {
	inFamily bool
	leaderID uint32
	precept  string
	members  []*FamilyMember
}

func FamilyTreeFromProto(pb *internal.Family) FamilyTree {
	if pb == nil {
		return FamilyTree{}
	}
	members := make([]*FamilyMember, 0, len(pb.GetMembers()))
	for _, m := range pb.GetMembers() {
		members = append(members, familyMemberFromProto(m))
	}
	tree := FamilyTree{precept: pb.GetPrecept(), members: members}
	if pb.LeaderId != nil {
		tree.inFamily = true
		tree.leaderID = pb.GetLeaderId()
	}
	return tree
}

func (t FamilyTree) InFamily() bool {
	return t.inFamily
}

func (t FamilyTree) Precept() string {
	return t.precept
}

func (t FamilyTree) Member(id uint32) *FamilyMember {
	if id == 0 {
		return nil
	}
	for _, m := range t.members {
		if m.characterID == id {
			return m
		}
	}
	return nil
}

func (t FamilyTree) Juniors(id uint32) []*FamilyMember {
	var juniors []*FamilyMember
	for _, m := range t.members {
		if m.seniorID == id {
			juniors = append(juniors, m)
		}
	}
	return juniors
}

// Leader is nil outside a family.
func (t FamilyTree) Leader() *FamilyMember {
	if !t.inFamily {
		return nil
	}
	return t.Member(t.leaderID)
}

func (t FamilyTree) Descendants(id uint32) uint32 {
	count := uint32(0)
	pending := t.Juniors(id)
	for len(pending) > 0 {
		junior := pending[len(pending)-1]
		pending = append(pending[:len(pending)-1], t.Juniors(junior.CharacterID())...)
		count++
	}
	return count
}

// Size is 0 outside a family.
func (t FamilyTree) Size() uint32 {
	if !t.InFamily() {
		return 0
	}
	return uint32(len(t.members))
}

// Pedigree is what the client draws for id: the leader as its own node, then two
// generations above and below id. The leader repeats when it is one of them, and
// a character outside a family is its own leader node.
func (t FamilyTree) Pedigree(id uint32) []*FamilyMember {
	me := t.Member(id)
	if me == nil {
		return nil
	}
	leader := t.Leader()
	if leader == nil {
		leader = me
	}
	visible := []*FamilyMember{leader}
	senior := t.Member(me.SeniorID())
	if senior != nil {
		if grand := t.Member(senior.SeniorID()); grand != nil {
			visible = append(visible, grand)
		}
		visible = append(visible, senior)
	}
	visible = append(visible, me)
	if senior != nil {
		for _, sibling := range t.Juniors(senior.CharacterID()) {
			if sibling != me {
				visible = append(visible, sibling)
			}
		}
	}
	juniors := t.Juniors(id)
	visible = append(visible, juniors...)
	for _, junior := range juniors {
		visible = append(visible, t.Juniors(junior.CharacterID())...)
	}
	return visible
}
