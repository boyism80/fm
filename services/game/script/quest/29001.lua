-- 칭호 도전 - 퀘스트 스페셜리스트 (Quest.wz/QuestData/29001.img): 칭호 도전 - 퀘스트 스페셜리스트

local quest_id = 29001

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		me:dialog(npc, "#v1112980# #e#b#t1112980##k\r\n\r\n - 퀘스트 1000개 완료\r\n\r\n#n이 훈장의 주인이 될 자격이 있는지 시험해 보시겠소?", false, true)
		me:dialog(npc, "자, 퀘스트 스페셜 리스트가 될 자격이 있으려면, 1000개의 퀘스트를 완료해야만 하네. 내 인생에 그런 모험가는 본적이 없지만. 그대라면 가능할지도 모른다고 생각하오.", false, true)
		q:start(npc, true)
	end,

	on_end = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if me:completed_quest_count() >= 1000 then
			local code = me:exchange({}, { item = { [1112980] = 1 } })
			if code == ExchangeResult.LackCapacity then
				me:dialog(npc, "장비창을 비워주세요.", false, false)
				return
			end
			if code ~= ExchangeResult.OK then
				return
			end
			q:force_complete(npc)
		else
			me:dialog(npc, "아직 그대는 #b" .. me:completed_quest_count() .. "개#k의 퀘스트를 완료하였군. 아직 퀘스트 스페셜 리스트가 될 자격을 충족하지 못한것같소. 퀘스트 스페셜리스트가 되기위한 자격은 1000개의 퀘스트 완료요.", false, false)
		end
	end
}
