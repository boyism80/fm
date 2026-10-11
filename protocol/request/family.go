package request

import internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"

type Family struct {
	Action   internal.FamilyAction
	TargetID uint32
	Feature  internal.FamilyFeature
	Name     string
	Accepted bool
}
