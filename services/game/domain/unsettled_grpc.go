package domain

import (
	"maps"
	"slices"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func NewUnsettledRequestsFromInternalProto(owner *Character, pbs []*internal.UnsettledRequest) *UnsettledRequests {
	u := &UnsettledRequests{owner: owner, requests: make(map[uint64]*UnsettledRequest, len(pbs))}
	for _, pb := range pbs {
		u.requests[pb.GetId()] = &UnsettledRequest{ID: pb.GetId(), Type: pb.GetType(), Target: pb.GetTarget(), Detail: pb.GetDetail()}
	}
	return u
}

func (u *UnsettledRequests) BuildRPC() []*internal.UnsettledRequest {
	pbs := make([]*internal.UnsettledRequest, 0, len(u.requests))
	for _, id := range slices.Sorted(maps.Keys(u.requests)) {
		request := u.requests[id]
		pbs = append(pbs, &internal.UnsettledRequest{Id: request.ID, Type: request.Type, Target: request.Target, Detail: request.Detail})
	}
	return pbs
}
