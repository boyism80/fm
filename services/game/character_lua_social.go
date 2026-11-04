package game

import (
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	lua "github.com/yuin/gopher-lua"
)

func (ch *Character) luaSocialFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"send_note": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			recipient := L.CheckString(2)
			message := L.CheckString(3)
			if !ClientText(message).Valid(200) {
				L.Push(lua.LFalse)
				return 1
			}
			if !ClientText(recipient).Valid(13) {
				ch.listener.OnNoteFailed(ch, internal.NoteResult_NOTE_RESULT_WRONG_NAME)
				L.Push(lua.LFalse)
				return 1
			}
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.Push(lua.LFalse)
				return 1
			}
			promise := ch.gameWorld.Systems().Get[*InternalSystem]().SendNoteAsync(cfg.ActorContext, ch, recipient, message)
			return LuaYieldPromise(L, ch.gameWorld, promise, func(reply *internal.NotesReply, err error) []lua.LValue {
				if err != nil {
					ch.listener.OnNoteFailed(ch, internal.NoteResult_NOTE_RESULT_INBOX_FULL)
					return []lua.LValue{lua.LFalse}
				}
				if reply.Status != internal.NoteResult_NOTE_RESULT_OK {
					ch.listener.OnNoteFailed(ch, reply.Status)
					return []lua.LValue{lua.LFalse}
				}
				ch.listener.OnNoteSent(ch)
				return []lua.LValue{lua.LTrue}
			})
		},
		"marriage": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if ch.wedding.marriage == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, ch.wedding.marriage))
			return 1
		},
		"reserve_wedding": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.Push(lua.LFalse)
				return 1
			}
			promise, err := ch.wedding.Reserve(cfg.ActorContext, uint32(L.CheckNumber(2)))
			if err != nil {
				L.Push(lua.LFalse)
				return 1
			}
			return LuaYieldPromise(L, ch.gameWorld, promise, func(result bool, err error) []lua.LValue {
				if err != nil {
					return []lua.LValue{lua.LFalse}
				}
				return []lua.LValue{lua.LBool(result)}
			})
		},
		"break_engagement": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.Push(lua.LFalse)
				return 1
			}
			promise, err := ch.wedding.BreakEngagement(cfg.ActorContext)
			if err != nil {
				L.Push(lua.LFalse)
				return 1
			}
			return LuaYieldPromise(L, ch.gameWorld, promise, func(result bool, err error) []lua.LValue {
				if err != nil {
					return []lua.LValue{lua.LFalse}
				}
				return []lua.LValue{lua.LBool(result)}
			})
		},
		"request_divorce": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.Push(lua.LNil)
				return 1
			}
			promise, err := ch.wedding.RequestDivorce(cfg.ActorContext)
			if err != nil {
				L.Push(lua.LNil)
				return 1
			}
			return LuaYieldPromise(L, ch.gameWorld, promise, func(result string, err error) []lua.LValue {
				if err != nil {
					return []lua.LValue{lua.LNil}
				}
				return []lua.LValue{lua.LString(result)}
			})
		},
		"cancel_divorce": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.Push(lua.LFalse)
				return 1
			}
			promise, err := ch.wedding.CancelDivorce(cfg.ActorContext)
			if err != nil {
				L.Push(lua.LFalse)
				return 1
			}
			return LuaYieldPromise(L, ch.gameWorld, promise, func(canceled bool, err error) []lua.LValue {
				return []lua.LValue{lua.LBool(err == nil && canceled)}
			})
		},
		"invited_to": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			L.Push(lua.LBool(ch.wedding.InvitedTo(uint32(L.CheckNumber(2)))))
			return 1
		},
		"open_wedding_wishlist": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			ch.listener.OnWeddingWishlistInput(ch)
			return 0
		},
		"open_wedding_gift": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			receiverID := uint32(L.CheckNumber(2))
			wishes := make([]string, 0)
			if tbl, ok := L.Get(3).(*lua.LTable); ok {
				tbl.ForEach(func(_, value lua.LValue) {
					wishes = append(wishes, value.String())
				})
			}
			ch.wedding.OpenGift(receiverID, wishes)
			return 0
		},
		"open_wedding_gift_box": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.Push(lua.LFalse)
				return 1
			}
			return LuaYieldPromise(L, ch.gameWorld, ch.wedding.OpenGiftBox(cfg.ActorContext), func(result bool, err error) []lua.LValue {
				if err != nil {
					return []lua.LValue{lua.LFalse}
				}
				return []lua.LValue{lua.LBool(result)}
			})
		},
		"field_team": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			L.Push(lua.LNumber(ch.fieldTeam))
			return 1
		},
		"party": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "party() takes no arguments")
				return 0
			}
			pid := ch.party.ID()
			if pid == nil || ch.gameWorld == nil {
				L.Push(lua.LNil)
				return 1
			}
			p := ch.gameWorld.Systems().Get[*PartySystem]().Get(*pid)
			if p == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, p))
			return 1
		},
		"state_machine": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if ch.stateMachine == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, ch.stateMachine))
			return 1
		},
		"try_party_quest": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			questID := uint32(L.CheckInt(2))
			ch.quests.RunHook(questID, "on_quest_try")
			return 0
		},
		"end_party_quest": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			questID := uint32(L.CheckInt(2))
			ch.quests.RunHook(questID, "on_quest_end")
			return 0
		},
		"carnival_team": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			team := ch.carnival.Team()
			if team == nil {
				L.Push(lua.LNil)
			} else {
				L.Push(luax.NewLuable(L, team))
			}
			return 1
		},
		"carnival_summon": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			ch.listener.OnCarnivalSummon(ch, constant.CarnivalTab(L.CheckInt(2)), uint8(L.CheckInt(3)), ch.GetName())
			return 0
		},
		"save_location": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			ch.savedLocations.Save(L.CheckString(2))
			return 0
		},
		"saved_location": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			mapID, ok := ch.savedLocations.Find(L.CheckString(2))
			if !ok {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(lua.LNumber(mapID))
			return 1
		},
		"clear_saved_location": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			ch.savedLocations.Clear(L.CheckString(2))
			return 0
		},
		"records": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			L.Push(luax.NewLuable(L, ch.records))
			return 1
		},
		"account_records": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			L.Push(luax.NewLuable(L, ch.account.records))
			return 1
		},
		"quest": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() != 2 {
				L.ArgError(2, "quest(quest_id) requires quest id")
				return 0
			}
			questID := uint32(L.CheckInt(2))
			qp := ch.quests.LuaQuest(questID)
			if qp == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, qp))
			return 1
		},
		"quests": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "quests() takes no arguments")
				return 0
			}
			tbl := L.NewTable()
			if ch.quests != nil {
				ch.quests.ForEach(func(questID uint32, qp *Quest) {
					if qp == nil || qp.status == QuestStatusNotStarted {
						return
					}
					tbl.RawSet(lua.LNumber(questID), luax.NewLuable(L, qp))
				})
			}
			L.Push(tbl)
			return 1
		},
		"start_quest": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() != 3 {
				L.ArgError(2, "start_quest(quest_id, npc) requires quest id and npc")
				return 0
			}
			questID := uint32(L.CheckInt(2))
			npcID, ok := LuaCheckNpcID(L, 3)
			if !ok {
				return 0
			}
			if ch.gameWorld == nil {
				L.Push(lua.LNil)
				return 1
			}
			_, err := ch.quests.Start(questID, QuestPhaseOpts{NpcID: &npcID})
			if err != nil {
				L.Push(lua.LNil)
				return 1
			}
			qp := ch.quests.Get(questID)
			if qp == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, qp))
			return 1
		},
		"clear_quests": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() == 1 {
				L.Push(lua.LNumber(ch.quests.ClearAll()))
				return 1
			}
			if L.GetTop() == 2 {
				questID := uint32(L.CheckInt(2))
				cleared := 0
				if ch.quests.Clear(questID) {
					cleared = 1
				}
				L.Push(lua.LNumber(cleared))
				return 1
			}
			L.ArgError(2, "clear_quests() or clear_quests(quest_id)")
			return 0
		},
		"completed_quest_count": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "completed_quest_count() takes no arguments")
				return 0
			}
			L.Push(lua.LNumber(ch.quests.CompletedCount()))
			return 1
		},
		"guild": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "guild() takes no arguments")
				return 0
			}
			guildID, inGuild := ch.guild.ID()
			if !inGuild || ch.gameWorld == nil {
				L.Push(lua.LNil)
				return 1
			}
			g := ch.gameWorld.Systems().Get[*GuildSystem]().Get(guildID)
			if g == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, g))
			return 1
		},
		"show_guild_ranking": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			npc := 0
			switch npcArg := L.Get(2).(type) {
			case lua.LNumber:
				npc = int(npcArg)
			case *lua.LUserData:
				if n, ok := npcArg.Value.(*Npc); ok && n.wz != nil && n.wz.BaseSpawn != nil {
					npc = int(n.wz.BaseSpawn.ID)
				}
			}
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil || ch.gameWorld == nil {
				return 0
			}
			ch.gameWorld.Systems().Get[*GuildSystem]().ShowRankingAsync(cfg.ActorContext, ch, uint32(npc))
			return 0
		},
		"overall_ranking": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			count := uint32(L.CheckInt(2))
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.Push(lua.LNil)
				return 1
			}
			promise := ch.gameWorld.Systems().Get[*InternalSystem]().GetRankingAsync(cfg.ActorContext, internal.RankingKind_RANKING_KIND_OVERALL, count)
			return LuaYieldPromise(L, ch.gameWorld, promise, func(reply *internal.GetRankingReply, err error) []lua.LValue {
				if err != nil {
					return []lua.LValue{lua.LNil}
				}
				tbl := L.NewTable()
				for _, e := range reply.GetEntries() {
					entry := L.NewTable()
					entry.RawSetString("rank", lua.LNumber(e.GetRank()))
					entry.RawSetString("name", lua.LString(e.GetName()))
					entry.RawSetString("level", lua.LNumber(e.GetLevel()))
					entry.RawSetString("class", lua.LNumber(e.GetClassId()))
					tbl.Append(entry)
				}
				return []lua.LValue{tbl}
			})
		},
		"gain_guild_gp": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			amount := L.CheckInt(2)
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				return 0
			}
			ch.gameWorld.Systems().Get[*GuildSystem]().GainGPAsync(cfg.ActorContext, ch, int32(amount))
			return 0
		},
		"create_alliance": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			return luaCreateAlliance(L, ch)
		},
		"disband_alliance": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			return luaDisbandAlliance(L, ch)
		},
		"inc_alliance_capacity": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			return luaIncAllianceCapacity(L, ch)
		},
		"family": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			self := ch.family.Self()
			table := L.NewTable()
			table.RawSetString("leader", lua.LNumber(ch.family.tree.Leader().CharacterID()))
			table.RawSetString("senior", lua.LNumber(self.SeniorID()))
			table.RawSetString("juniors", lua.LNumber(len(ch.family.tree.Juniors(ch.GetID()))))
			table.RawSetString("reputation", lua.LNumber(self.Reputation()))
			table.RawSetString("total_reputation", lua.LNumber(self.TotalReputation()))
			L.Push(table)
			return 1
		},
		"buddy_capacity": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			L.Push(lua.LNumber(ch.buddies.Capacity()))
			return 1
		},
		"expand_buddy_list": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.Push(lua.LFalse)
				return 1
			}
			promise, err := ch.buddies.Expand(cfg.ActorContext, int32(L.CheckInt(2)))
			if err != nil {
				L.Push(lua.LFalse)
				return 1
			}
			return LuaYieldPromise(L, ch.gameWorld, promise, func(expanded bool, err error) []lua.LValue {
				if err != nil {
					return []lua.LValue{lua.LNil}
				}
				return []lua.LValue{lua.LBool(expanded)}
			})
		},
	}
}
