package entity

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func (s *SavedLocations) Load(pbs []*internal.SavedLocation) {
	for _, pb := range pbs {
		if pb == nil {
			continue
		}
		name := pb.GetName()
		if name == "" {
			continue
		}
		s.entries[name] = pb.GetMapId()
	}
}

func (s *SavedLocations) ToProto() []*internal.SavedLocation {
	if len(s.entries) == 0 {
		return nil
	}
	out := make([]*internal.SavedLocation, 0, len(s.entries))
	for name, mapID := range s.entries {
		out = append(out, &internal.SavedLocation{
			CharacterId: s.owner.GetID(),
			Name:        name,
			MapId:       mapID,
		})
	}
	return out
}
