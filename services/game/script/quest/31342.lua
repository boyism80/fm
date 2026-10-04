-- [암벽 거인] 치노와 함께 (Quest.wz/QuestInfo.img.xml): [암벽 거인] 치노와 함께

local quest_id = 31342

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if not me:dialog_yes_no(npc, "그럼 출발할까? 승강기를 타고 암벽거인의 몸을 타고 올라 가는 거야. 워낙 거대한 몸 위에 올라가는 것이니만큼 시간이 좀 걸려. 준비를 단단히 하도록 해. 대부분 벌떼가 습격을 하는데 가끔식 운이 좋은 확률으로 벌떼들이 습격을 하지 않을수도 있어~") then
			return
		end

		q:start(npc, true)
		me:map(240091600)
		me:message("운이 좋게 벌떼가 나타나지 않았습니다.", 5)
	end
}
