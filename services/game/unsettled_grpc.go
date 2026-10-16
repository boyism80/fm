package game

import (
	"maps"
	"slices"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func NewUnsettledRequestsFromInternalProto(owner *Character, pbs []*internal.UnsettledRequest) *UnsettledRequests {
	u := &UnsettledRequests{owner: owner, requests: make(map[uint64]*UnsettledRequest, len(pbs))}
	for _, pb := range pbs {
		u.requests[pb.GetId()] = &UnsettledRequest{id: pb.GetId(), onType: pb.GetType(), target: pb.GetTarget(), detail: pb.GetDetail()}
	}
	return u
}

func (u *UnsettledRequests) BuildRPC() []*internal.UnsettledRequest {
	pbs := make([]*internal.UnsettledRequest, 0, len(u.requests))
	for _, id := range slices.Sorted(maps.Keys(u.requests)) {
		request := u.requests[id]
		pbs = append(pbs, &internal.UnsettledRequest{Id: request.id, Type: request.onType, Target: request.target, Detail: request.detail})
	}
	return pbs
}
