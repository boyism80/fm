-- 페티트의 메이플 메이트 (Quest.wz/QuestData/9251.img): 페티트의 메이플 메이트

local quest_id = 9251

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	q:start(npc, true)
	me:show_quest_completion(9251)
end

function on_end(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	q:force_complete(npc)
	me:show_effect(EffectType.QuestCompletion)
end
