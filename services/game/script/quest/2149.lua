local quest_id = 2149

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		local sel = me:dialog_list(npc, "...무슨 일이지?", {
			"귀신나무에 대해 들어보신 적이 있나요?",
		})
		if sel == nil then
			return
		end

		sel = me:dialog_list(npc, "겁쟁이들의 이야기를 들은 모양이군. 귀신나무라니 그런게 있을리가 없지. 오랫동안 페리온의 바위산을 돌아다니며 수련을 했지만, 그런 나무는 본 적도 들은 적도 없어.", {
			"아 그런가요?...",
		})
		if sel == nil then
			return
		end

		me:dialog(npc, "단지 요즘 동쪽 바위산에서 의문의 습격을 받는 일이 늘어났다고 하는데, 조금 신경이 쓰이는군...", false, true)

		q:start(npc, true)
		q:force_complete(npc)
		me:show_effect(EffectType.QuestCompletion)
		me:exp(me:exp() + 100)
	end
}
