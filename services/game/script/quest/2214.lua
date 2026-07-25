local quest_id = 2214

return {
	on_end = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if not q:started() then
			q:start(npc, true)
			return
		end

		local hour = datetime().hour
		if hour < 17 or hour > 20 then
			me:dialog(npc, "#b(먼지 투성이인 쓰레기통 내부에 손을 넣어봤지만, 딱히 쓰레기가 아닌 물건은 없는 것 같습니다.)")
			return
		end

		if not me:dialog_yes_no(npc, "#b(먼지 투성이인 쓰레기통 내부에 손을 넣자 바스락거리는 뭔가가 만져집니다. 뭔가를 꺼내겠습니까?)") then
			return
		end

		me:dialog(npc, "#b(거미줄로 범벅이 된 구겨진 종이조각을 꺼냈다.)", false, true)

		local code = me:exchange({}, { item = { [4031894] = 1 } })
		if code == ExchangeResult.LackCapacity then
			return
		end
		if code ~= ExchangeResult.OK then
			return
		end

		q:force_complete(npc)
		me:show_effect(EffectType.QuestCompletion)
	end
}
