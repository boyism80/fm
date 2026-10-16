package game

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

// Account is shared by every character of the account in this world. Only
// one of them is online at a time, so the online character owns it.
type Account struct {
	storage *Storage
	records *Records
}

func NewAccountFromInternalProto(pb *internal.Account, owner *Character) *Account {
	return &Account{
		storage: NewStorageFromInternalProto(pb.GetStorage(), owner),
		records: NewRecordsFromInternalProto(pb.GetRecords()),
	}
}

func (a *Account) BuildRPC(worldID uint32) *internal.Account {
	return &internal.Account{
		Storage: a.storage.BuildRPC(worldID),
		Records: a.records.BuildRPC(),
	}
}
