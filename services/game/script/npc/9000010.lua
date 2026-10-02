-- NPC name (String.wz/Npc.img.xml): 피에트라

local EVENT_TICKET = 4031018
local REWARD = 4000038
local HOLY_SYMBOL = 9001002
local FALLBACK_MAP = 100000000

return {
	on_click = function(me, npc)
		if next(me:item(EVENT_TICKET)) ~= nil then
			local sel = me:dialog_list(npc, "#b#t4031018##k를 #p9000006#에게서 보상으로 교환할 수 있답니다. 알고 싶으신것이 있으세요?\r\n", {
				"#b#p9000006#은 누구죠?",
				"이전에 있던곳으로 데려다 주세요.",
			})
			if sel == nil then
				return
			end
			if sel == 1 then
				me:dialog(npc, "#b#p9000006##k에게 #t4031018#를 건네고, 이벤트 보상으로 교환할 수 있습니다. 제 옆에 바로 있는 엔피시입니다.")
				return
			end
		end

		if me:dialog(npc, "이전에 있던 곳으로 데려다 드리겠습니다.", false, true) == false then
			return
		end

		if me:exchange(nil, { item = { [REWARD] = 1 } }) ~= ExchangeResult.OK then
			me:dialog(npc, "인벤토리 공간이 충분한지 확인해 주세요.")
			return
		end

		local back = me:saved_location("EVENT")
		if back == nil then
			back = FALLBACK_MAP
		end
		me:clear_saved_location("EVENT")
		local skill = me:add_skill(HOLY_SYMBOL)
		if skill ~= nil then
			me:buff(skill, BuffFlag.HolySymbol, 50, { time = 1200 })
		end
		me:map(back)
		if skill == nil then
			return
		end
		me:message("20분 동안 경험치 1.5배 버프가 적용되었습니다. 게임 재접속, 혹은 홀리심볼 버프를 받게되면 해제되니 주의해주세요.", Msg.PinkText)
		me:message("20분 동안 경험치 1.5배 버프가 적용되었습니다. \r\n\r\n게임 재접속, 혹은 홀리심볼 버프를 받게되면 해제되니 주의해주세요.", Msg.Popup)
	end
}
