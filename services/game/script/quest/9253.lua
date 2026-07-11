-- 메이플 핸디월드 (Quest.wz/QuestData/9253.img): 메이플 핸디월드

local quest_id = 9253

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	q:start(npc, true)
	me:show_quest_completion(9253)
end
