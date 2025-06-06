package entity

import "github.com/asynkron/protoactor-go/actor"

type CharacterListener interface {
	GetContext() actor.Context
	OnDialog(message string, prev bool, next bool)
	OnDialogYesNo(message string, prev bool, next bool)
	OnDialogAccept(message string, enableEscape bool)
	OnDialogList(message string, selections []string)
	OnDialogInput(message string)
	OnChat(message string, highlight bool, dontRecordHistory bool)
}
