package entity

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

// Account is shared by every character of the account in this world. Only
// one of them is online at a time, so the online character owns it.
type Account struct {
	Storage *Storage
	Records *Records
}

func NewAccountFromInternalProto(pb *internal.Account, owner *Character) *Account {
	return &Account{
		Storage: NewStorageFromInternalProto(pb.GetStorage(), owner),
		Records: NewRecordsFromInternalProto(pb.GetRecords()),
	}
}

func (a *Account) ToProto(worldID uint32) *internal.Account {
	return &internal.Account{
		Storage: a.Storage.ToProto(worldID),
		Records: a.Records.ToProto(),
	}
}
