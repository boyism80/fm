-- 스탠의 실기 수업 (Quest.wz/QuestData/6032.img.xml): 스탠의 실기 수업

local quest_id = 6032

return {
	on_end = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if not q:started() and not q:completed() then
			q:start(npc, true)
			return
		end

		if not me:dialog(npc, "자! 실기수업을 시작하겠다! 오늘의 실습 주제는 추귀고리 만들기야! 추귀고리를 만들기 위해서는 강철 4개, 하급 몬스터 결정 1개, 청동 1개가 필요하지. 자 이제부터 잘 보게. 추귀고리를 만들기 위해서는 중력의 연성법칙이 필요하거든.", false, true) then
			return
		end
		if not me:dialog(npc, "(스탠이 연성진을 모두 완성하자, 연성진에서 빛이 나면서 엄청난 섬광이 눈앞을 가렸다.)", false, true) then
			return
		end
		if not me:dialog_yes_no(npc, "자 모두 이해했나? 여기서 수업을 마치도록 하겠다.") then
			me:dialog(npc, "흠, 이해하지 못한 것 같으니 한번만 다시 설명해 주도록 하지. 잠시 후에 다시 말을 걸게나.", false, false)
			return
		end

		local qr = me:quest(6029)
		if qr ~= nil then
			local info = qr:record()
			if info == nil or info == "" then
				info = "000"
			end
			qr:record(info:sub(1, 2) .. "1")
		end
		q:force_complete(npc)
		me:show_effect(EffectType.QuestCompletion)
	end
}
