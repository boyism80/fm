package game

import (
	"log"
	"time"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
)

func (ch *Character) luaEffectFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"chat": func(L *lua.LState) int {
			argc := L.GetTop()
			ch := luax.Check[Character](L, 1)

			message := L.CheckString(2)
			highlight := false
			dontRecordHistory := false
			switch argc {
			case 4:
				dontRecordHistory = L.CheckBool(4)
				fallthrough
			case 3:
				highlight = L.CheckBool(3)
			}

			ch.listener.OnChat(ch, message, highlight, dontRecordHistory)
			return 0
		},
		"show_skill_effect": func(L *lua.LState) int {
			argc := L.GetTop()
			ch := luax.Check[Character](L, 1)

			skillEntry := luax.Check[SkillEntry](L, 2)

			skillLevelInt := skillEntry.Level()
			if skillLevelInt < 0 {
				skillLevelInt = 0
			}
			skillLevel := uint8(skillLevelInt)
			if argc >= 3 {
				effectType := L.CheckInt(3)
				if effectType == int(pconst.SkillEffectTypeCast) || effectType == int(pconst.SkillEffectTypeAffected) {
					ch.listener.OnShowSelfSkillEffect(ch, pconst.SkillEffectType(effectType), skillEntry.wz.ID, skillLevel, nil)
					return 0
				}
				L.ArgError(3, "show_skill_effect() supports SkillEffectType.Cast or SkillEffectType.Affected")
				return 0
			}

			ch.listener.OnShowSelfSkillEffect(ch, pconst.SkillEffectTypeCast, skillEntry.wz.ID, skillLevel, nil)
			return 0
		},
		"show_effect": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			effectType := response.EffectType(L.CheckInt(2))
			switch effectType {
			case response.EffectTypeLevelUp,
				response.EffectTypeClassChange,
				response.EffectTypeQuestCompletion,
				response.EffectTypeRegisterCard,
				response.EffectTypeItemLevelUp:
				ch.listener.OnShowSelfEffect(ch, effectType)
				return 0
			default:
				L.ArgError(2, "show_effect() supports EffectType.LevelUp, ClassChange, QuestCompletion, RegisterCard, or ItemLevelUp")
				return 0
			}
		},
		"show_quest_completion": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() != 2 {
				L.ArgError(2, "show_quest_completion(quest_id) requires quest id")
				return 0
			}
			questID := uint32(L.CheckInt(2))
			if ch.listener != nil {
				ch.listener.OnShowQuestCompletion(ch, questID)
			}
			return 0
		},
		"play_sound": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			if argc < 2 || argc > 3 {
				L.ArgError(2, "play_sound(name [, broadcast])")
				return 0
			}
			sound := L.CheckString(2)
			broadcast := false
			if argc >= 3 {
				broadcast = L.CheckBool(3)
			}
			if ch.listener != nil {
				ch.listener.OnPlaySound(ch, sound, broadcast)
			}
			return 0
		},
		"play_portal_sound": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "play_portal_sound() takes no arguments")
				return 0
			}
			if ch.listener != nil {
				ch.listener.OnPlayPortalSound(ch)
			}
			return 0
		},
		"show_dragon_blood_effect": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			skillEntry := luax.Check[SkillEntry](L, 2)
			if skillEntry.wz == nil {
				L.ArgError(2, "SkillEntry with Wz expected")
				return 0
			}
			skillLevelInt := skillEntry.Level()
			if skillLevelInt < 0 {
				skillLevelInt = 0
			}
			ch.listener.OnShowSelfDragonBloodEffect(ch, skillEntry.wz.ID, uint8(skillLevelInt))
			return 0
		},
		"show_hp_healed_effect": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			amount := int32(L.CheckInt(2))
			ch.listener.OnShowSelfHPHealedEffect(ch, amount)
			return 0
		},
		"show_incubator_result": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			ch.listener.OnIncubatorResult(ch, uint32(L.OptInt(2, 0)), int16(L.OptInt(3, 0)))
			return 0
		},
		"show_reward_item_animation": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			itemID := uint32(L.CheckInt(2))
			effectText := L.CheckString(3)
			ch.listener.OnShowSelfRewardItemAnimation(ch, itemID, effectText)
			return 0
		},
		"show_item_maker_success_effect": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			ch.listener.OnShowSelfItemMakerSuccessEffect(ch)
			return 0
		},
		"show_crafting_effect": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			effectText := L.CheckString(2)
			timeValue := int32(L.CheckInt(3))
			modeValue := int32(L.CheckInt(4))
			ch.listener.OnShowSelfCraftingEffect(ch, effectText, timeValue, modeValue)
			return 0
		},
		"show_dice_effect": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			effectID := int32(L.CheckInt(2))
			skillEntry := luax.Check[SkillEntry](L, 3)
			if skillEntry.wz == nil {
				L.ArgError(3, "SkillEntry with Wz expected")
				return 0
			}
			skillLevelInt := skillEntry.Level()
			if skillLevelInt < 0 {
				skillLevelInt = 0
			}
			ch.listener.OnShowSelfDiceEffect(ch, effectID, skillEntry.wz.ID, uint8(skillLevelInt))
			return 0
		},

		"message": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			text := L.CheckString(2)
			msgType := constant.MsgLightBlueText
			scope := constant.MessageScopeSelf
			ear := false
			top := L.GetTop()
			if top >= 3 && L.Get(3) != lua.LNil {
				msgType = constant.ServerMessageType(L.CheckInt(3))
			}
			if top >= 4 && L.Get(4) != lua.LNil {
				scope = constant.MessageScope(L.CheckInt(4))
			}
			if top >= 5 {
				ear = L.ToBool(5)
			}
			switch scope {
			case constant.MessageScopeMap:
				if m := ch.GetMap(); m != nil {
					m.MapMessage(msgType, text)
				}
			case constant.MessageScopeChannel:
				if ch.gameWorld != nil {
					ch.gameWorld.Systems().Get[*NoticeSystem]().BroadcastNotice(msgType, text, 0, false)
				}
			case constant.MessageScopeWorld:
				cfg, ok := luax.GetConfiguration(L)
				if !ok || cfg.ActorContext == nil || ch.gameWorld == nil {
					L.Push(lua.LFalse)
					return 1
				}
				promise := ch.gameWorld.Systems().Get[*InternalSystem]().BroadcastNoticeAsync(cfg.ActorContext, ch, msgType, text, ear)
				if promise == nil {
					L.Push(lua.LFalse)
					return 1
				}
				return LuaYieldPromise(L, ch.gameWorld, promise, func(_ *internal.BroadcastNoticeReply, err error) []lua.LValue {
					if err != nil {
						return []lua.LValue{lua.LFalse}
					}
					return []lua.LValue{lua.LTrue}
				})
			default:
				ch.listener.OnMessage(ch, msgType, text)
			}
			L.Push(lua.LTrue)
			return 1
		},
		"megaphone_lines": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			tbl := L.CheckTable(2)
			ear := L.ToBool(3)
			lines := make([]string, 0, 3)
			for i := 1; i <= tbl.Len(); i++ {
				lines = append(lines, lua.LVAsString(tbl.RawGetInt(i)))
			}
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil || len(lines) < 1 || len(lines) > 3 {
				L.Push(lua.LFalse)
				return 1
			}
			promise := ch.gameWorld.Systems().Get[*InternalSystem]().BroadcastMegaphoneAsync(cfg.ActorContext, ch, &response.Megaphone{
				Type:  constant.MsgThreeMegaphoneLines,
				Lines: lines,
				Ear:   ear,
			})
			return LuaYieldPromise(L, ch.gameWorld, promise, func(_ *internal.BroadcastNoticeReply, err error) []lua.LValue {
				return []lua.LValue{lua.LBool(err == nil)}
			})
		},
		"item_megaphone": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			megaphone := &response.Megaphone{
				Type:  constant.MsgItemMegaphone,
				Lines: []string{L.CheckString(2)},
				Ear:   L.ToBool(3),
			}
			if L.Get(4) != lua.LNil {
				invType := constant.InventoryType(L.CheckInt(4))
				slot := L.CheckInt(5)
				if invType < constant.InventoryTypeEquipment || invType > constant.InventoryTypeCash || slot < 1 || slot > 255 {
					L.Push(lua.LFalse)
					return 1
				}
				item := ch.inventory.tabs[invType].Get(uint8(slot))
				if item == nil {
					L.Push(lua.LFalse)
					return 1
				}
				megaphone.Item = item.BuildS2C()
			}
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.Push(lua.LFalse)
				return 1
			}
			promise := ch.gameWorld.Systems().Get[*InternalSystem]().BroadcastMegaphoneAsync(cfg.ActorContext, ch, megaphone)
			return LuaYieldPromise(L, ch.gameWorld, promise, func(_ *internal.BroadcastNoticeReply, err error) []lua.LValue {
				return []lua.LValue{lua.LBool(err == nil)}
			})
		},
		"set_chalkboard": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			L.Push(lua.LBool(ch.SetChalkboard(L.CheckString(2))))
			return 1
		},
		"close_chalkboard": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			ch.CloseChalkboard()
			return 0
		},
		"place_message_box": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			box := NewMessageBox(ch, uint32(L.CheckInt(2)), L.CheckString(3))
			L.Push(lua.LBool(box != nil && box.Place(ch)))
			return 1
		},
		"snowball_knockback": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			ch.listener.OnSnowballKnockback(ch)
			return 0
		},
		"map_effect": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			ch.listener.OnMapEffect(ch, L.CheckString(2))
			return 0
		},
		"show_instruction": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() < 2 || L.GetTop() > 4 {
				L.ArgError(2, "show_instruction(text[, width[, height]])")
				return 0
			}
			text := L.CheckString(2)
			width := uint16(0)
			height := uint16(0)
			if L.GetTop() >= 3 {
				width = uint16(L.CheckInt(3))
			}
			if L.GetTop() >= 4 {
				height = uint16(L.CheckInt(4))
			}
			ch.Send(&response.Hint{
				Text:   text,
				Width:  width,
				Height: height,
			}, types.SEND_POLICY_ENCRYPT)
			return 0
		},
		"generic_guild_message": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			ch.listener.OnGuildMessage(ch, pconst.GuildResponseCode(L.CheckInt(2)))
			return 0
		},
		"event_instructions": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			ch.listener.OnEventInstructions(ch)
			return 0
		},
		"stop_clock": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			ch.RemoveTimer(clockTimer)
			ch.listener.OnStopClock(ch)
			return 0
		},
		"clock": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			seconds := L.CheckInt(2)
			fn := L.OptFunction(3, nil)
			ch.RemoveTimer(clockTimer)
			ch.listener.OnClock(ch, int32(seconds))
			if fn == nil {
				return 0
			}

			// The client drops the clock with the field, so the timer never outlives this map or its Lua root.
			added := ch.AddTimer(clockTimer, time.Duration(seconds)*time.Second, false, func() {
				m := ch.GetMap()
				if m == nil {
					return
				}
				root := m.GetLuaRoot()
				if root == nil {
					return
				}
				if _, err := luax.CallFunction(root, fn, ch); err != nil {
					log.Printf("clock: %v", err)
				}
			})
			if !added {
				L.RaiseError("clock: timer not added")
			}
			return 0
		},
	}
}
