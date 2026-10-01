-- NPC name (String.wz/Npc.img.xml): 휴머노이드 A

local QUEST = 3335
local SNOW_ROSE = 4031695

return {
	on_click = function(me, npc)
		if me:quest(QUEST):started() == false then
			me:dialog(npc, "인간이 되고 싶습니다. 따뜻한 심장을 가진 인간이... 인간이 된다면 그녀의 손을 잡아줄 수도 있겠지요. 하지만 지금은 그럴 수 없지요...")
			return
		end
		if next(me:item(SNOW_ROSE)) ~= nil then
			me:dialog(npc, "이미 설원 장미를 가져 오셨군요. 그 설원 장미를 필리아씨에게 가져다 주세요.")
			return
		end
		if me:dialog_yes_no(npc, "와주셨군요... 설원 장미를 피울 준비는 되셨나요? 5월의 이슬이 있어야만 장미를 피울 수 있다는 건 알고 계시죠?") == false then
			me:dialog(npc, "준비가 되시면 저를 다시 찾아와 주세요.")
			return
		end
		if me:dialog(npc, "그럼 설원 장미를 피울 부화기가 마련된 곳으로 당신을 안내하겠습니다...", false, true) == false then
			return
		end
		local sm, err = state_machine("snow_rose"):start_solo(me)
		if sm == nil then
			log("snow_rose start_solo:", err)
			me:dialog(npc, "이미 다른 누군가가 이 안에서 설원 장미를 피우고 있는 것 같습니다. 다음에 다시 찾아와 주세요.")
		end
	end
}
