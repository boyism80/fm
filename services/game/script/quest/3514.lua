-- Quest name (Quest.wz/Quest.img.xml): 감정을 파는 마제사

local quest_id = 3514

local function item_count(me, item_id)
	local slots = me:item(item_id)
	local count = 0
	for _, it in pairs(slots) do
		count = count + it:count()
	end
	return count
end

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

		if item_count(me, 2022337) >= 1 then
			return
		end

		if not me:dialog(npc, "호오~ 약은 다 먹은 모양이군. 어때? 말 그대로 최고의 효과 아니야? 역시 이 몸의 약이란 완벽해!", false, true) then
			return
		end
		if not me:dialog(npc, "뭐? ...그냥 체력이 왕창 떨어지면 되는 거 아니냐고? 흠흠. 누구야? 그런 헛소리를 하는 게... 그럴 리가 없잖아? 하하하하!", false, true) then
			return
		end
		q:force_complete(npc)
		me:show_effect(EffectType.QuestCompletion)
	end
}
