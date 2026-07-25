-- 핸드폰 메이플스토리 도적편 (Quest.wz/QuestData/9252.img): 핸드폰 메이플스토리 도적편

local quest_id = 9252

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		q:start(npc, true)
		me:show_quest_completion(9252)
	end
}
