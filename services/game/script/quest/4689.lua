-- [암벽 거인] 타란튤로스 전투 (Quest.wz/QuestInfo.img.xml): 타란튤로스 전투

local quest_id = 4689

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if not me:dialog_yes_no(npc, "더이상 여기로부터 도망갈 수 없을 것 같아... \r\n어떻하지? 여기서 싸운다면 좀 더 동료가 필요할지도 몰라... 어떻게 할 거야? 싸워?") then
			return
		end

		q:start(npc, true)
		me:map(802000309)
	end,

	on_end = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if not me:dialog(npc, "타란튤로스를 해치운건가?", false, true) then
			return
		end
		if not me:dialog_yes_no(npc, "타란튤로스를 해치우자 암벽 거인의 떨림이 잦아들고 맑은 기운이 사방에서 샘솟기 시작한다. 암벽 거인들 위기로부터 구출해냈다. 이제 당분간은 오염되지 않을 것 같다.") then
			return
		end

		me:map(240092100)
	end
}
