package server

import (
	"fmt"

	"github.com/boyism80/fm/core/mq"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/client"
)

func (gs *GameServer) mqQueues() []mq.Queue {
	world, channel := gs.config.WorldId, gs.config.ChannelId
	toWorld := func(topic string) string { return fmt.Sprintf("fm.%d.all.%s", world, topic) }
	toChannel := func(topic string) string { return fmt.Sprintf("fm.%d.%d.%s", world, channel, topic) }
	queue := func(topic string, routingKeys []string, bind func(d *mq.Dispatcher)) mq.Queue {
		return mq.Queue{
			Name:        fmt.Sprintf("fm.game.w%d.c%d.%s.events", world, channel, topic),
			ConsumerTag: fmt.Sprintf("fm-game-w%d-c%d-%s", world, channel, topic),
			ActorName:   fmt.Sprintf("rabbitmq_%s_w%d_c%d", topic, world, channel),
			RoutingKeys: routingKeys,
			Bind:        bind,
		}
	}
	return []mq.Queue{
		queue("party", []string{toWorld("party"), toChannel("party")}, func(d *mq.Dispatcher) {
			d.Bind[partyMqMemberJoined](gs)
			d.Bind[partyMqMemberLeft](gs)
			d.Bind[partyMqLeaderChanged](gs)
			d.Bind[partyMqLogOnOff](gs)
			d.Bind[partyMqDisbanded](gs)
			d.Bind[partyMqPartySync](gs)
			d.Bind[partyMqPartyInvite](gs)
			d.Bind[partyMqChat](gs)
			d.Bind[partyMqPartyInviteDenied](gs)
		}),
		queue("buddy", []string{toWorld("buddy")}, func(d *mq.Dispatcher) {
			d.Bind[buddyMqChannelUpdate](gs)
			d.Bind[buddyMqUpdate](gs)
			d.Bind[buddyMqAddRequest](gs)
			d.Bind[buddyMqChat](gs)
		}),
		queue("guild", []string{toWorld("guild")}, func(d *mq.Dispatcher) {
			d.Bind[guildMqCreated](gs)
			d.Bind[guildMqMemberJoined](gs)
			d.Bind[guildMqMemberLeft](gs)
			d.Bind[guildMqRankTitlesChanged](gs)
			d.Bind[guildMqMemberRankChanged](gs)
			d.Bind[guildMqEmblemChanged](gs)
			d.Bind[guildMqNoticeChanged](gs)
			d.Bind[guildMqCapacityChanged](gs)
			d.Bind[guildMqGPChanged](gs)
			d.Bind[guildMqMemberOnlineChanged](gs)
			d.Bind[guildMqDisbanded](gs)
			d.Bind[guildMqChat](gs)
			d.Bind[guildMqMessage](gs)
		}),
		queue("alliance", []string{toWorld("alliance")}, func(d *mq.Dispatcher) {
			d.Bind[allianceMqCreated](gs)
			d.Bind[allianceMqDisbanded](gs)
			d.Bind[allianceMqGuildLeft](gs)
			d.Bind[allianceMqGuildAdded](gs)
			d.Bind[allianceMqCapacityChanged](gs)
			d.Bind[allianceMqRankTitlesChanged](gs)
			d.Bind[allianceMqMemberRankChanged](gs)
			d.Bind[allianceMqLeaderChanged](gs)
			d.Bind[allianceMqNoticeChanged](gs)
			d.Bind[allianceMqChat](gs)
		}),
		queue("global", []string{toWorld("global")}, func(d *mq.Dispatcher) {
			d.Bind[globalMqServerDatetime](gs)
			d.Bind[globalMqNotice](gs)
			d.Bind[playerNpcMqCreated](gs)
			d.Bind[playerNpcMqRemoved](gs)
		}),
		queue("parcel", []string{toChannel("parcel")}, func(d *mq.Dispatcher) {
			d.Bind[parcelMqArrived](gs)
			d.Bind[whisperMq](gs)
			d.Bind[messengerMq](gs)
			d.Bind[familyMq](gs)
			d.Bind[noteMqArrived](gs)
			d.Bind[cashGiftMqArrived](gs)
			d.Bind[marriageMqChanged](gs)
			d.Bind[marriageMqSpouseMoved](gs)
		}),
	}
}

func (gs *GameServer) ping() *internal.PingRequest {
	online := uint32(0)
	for _, c := range gs.Clients() {
		if gc, ok := c.(*client.GameClient); ok && gc.GetCharacter() != nil {
			online++
		}
	}
	return &internal.PingRequest{
		Role:        internal.ServerRole_SERVER_ROLE_GAME,
		WorldId:     gs.config.WorldId,
		ChannelId:   gs.config.ChannelId,
		OnlineUsers: online,
	}
}
