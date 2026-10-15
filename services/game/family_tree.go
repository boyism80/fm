package game

import internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"

// FamilyTree is a whole family as internal sent it. Every member's character
// shares the same message, so it is never modified.
type FamilyTree struct {
	family *internal.Family
}

func (t FamilyTree) InFamily() bool {
	return t.family.LeaderId != nil
}

func (t FamilyTree) Precept() string {
	return t.family.GetPrecept()
}

func (t FamilyTree) Member(id uint32) *internal.FamilyMember {
	if id == 0 {
		return nil
	}
	for _, m := range t.family.GetMembers() {
		if m.GetCharacterId() == id {
			return m
		}
	}
	return nil
}

func (t FamilyTree) Juniors(id uint32) []*internal.FamilyMember {
	var juniors []*internal.FamilyMember
	for _, m := range t.family.GetMembers() {
		if m.GetSeniorId() == id {
			juniors = append(juniors, m)
		}
	}
	return juniors
}

// Leader is nil outside a family.
func (t FamilyTree) Leader() *internal.FamilyMember {
	return t.Member(t.family.GetLeaderId())
}

func (t FamilyTree) Descendants(id uint32) uint32 {
	count := uint32(0)
	pending := t.Juniors(id)
	for len(pending) > 0 {
		junior := pending[len(pending)-1]
		pending = append(pending[:len(pending)-1], t.Juniors(junior.GetCharacterId())...)
		count++
	}
	return count
}

// Size is 0 outside a family.
func (t FamilyTree) Size() uint32 {
	if !t.InFamily() {
		return 0
	}
	return uint32(len(t.family.GetMembers()))
}

// Pedigree is what the client draws for id: the leader as its own node, then two
// generations above and below id. The leader repeats when it is one of them, and
// a character outside a family is its own leader node.
func (t FamilyTree) Pedigree(id uint32) []*internal.FamilyMember {
	me := t.Member(id)
	leader := t.Leader()
	if leader == nil {
		leader = me
	}
	visible := []*internal.FamilyMember{leader}
	senior := t.Member(me.GetSeniorId())
	if grand := t.Member(senior.GetSeniorId()); grand != nil {
		visible = append(visible, grand)
	}
	if senior != nil {
		visible = append(visible, senior)
	}
	visible = append(visible, me)
	if senior != nil {
		for _, sibling := range t.Juniors(senior.GetCharacterId()) {
			if sibling != me {
				visible = append(visible, sibling)
			}
		}
	}
	juniors := t.Juniors(id)
	visible = append(visible, juniors...)
	for _, junior := range juniors {
		visible = append(visible, t.Juniors(junior.GetCharacterId())...)
	}
	return visible
}
