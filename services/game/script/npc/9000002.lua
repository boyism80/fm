-- NPC name (String.wz/Npc.img.xml): 피에트로

local ANCIENT_SCROLL = 4031019
local ANCIENT_SCROLL_MINUTES = 3 * 24 * 60
local FALLBACK_MAP = 100000000

return {
	on_click = function(me, npc)
		if me:dialog(npc, "빰빠바빰~ 축하합니다! #b게임 이벤트#k를 완수하셨습니다.", false, true) == false then
			return
		end
		if me:dialog(npc, "이벤트 상품으로 고대문자로 작성된 비밀 정보가 들어있는 스크롤, #b#t4031019##k를 드립니다. #r이 아이템은 버리면 회수가 불가능 하니 절대 버리시면 안돼요~#k", true, true) == false then
			return
		end
		if me:dialog(npc, "#t4031019#는 #r#p9000007##k, 또는 루디브리엄 성에 있는 #r지니#k에게 가져가면 해독하실 수 있습니다. 받기 전에 인벤토리 공간이 충분한지 확인해 주세요. 그럼 행운을 빕니다~", true, true) == false then
			return
		end

		if next(me:item(ANCIENT_SCROLL)) == nil then
			local code = me:exchange(nil, {
				item = { [ANCIENT_SCROLL] = 1 },
				period = { [ANCIENT_SCROLL] = ANCIENT_SCROLL_MINUTES },
			})
			if code ~= ExchangeResult.OK then
				me:dialog(npc, "인벤토리 공간이 부족하신 건 아닌가요? 다시 한번 확인해 주세요~")
				return
			end
		end

		local back = me:saved_location("EVENT")
		if back == nil then
			back = FALLBACK_MAP
		end
		me:clear_saved_location("EVENT")
		me:map(back)
	end
}
