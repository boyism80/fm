local quest_id = 2151

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		local sel = me:dialog_list(npc, "무슨 일로 나를 찾아온건가?", {
			"귀신나무에 대해 아시는 것이 있나요?",
		})
		if sel == nil then
			return
		end

		sel = me:dialog_list(npc, "귀신나무라... 아마도 스텀피를 말하는 것 같군.", {
			"스텀피가 뭔가요?",
		})
		if sel == nil then
			return
		end

		sel = me:dialog_list(npc, "페리온이 아직 푸른 숲이었을때부터 지금까지 살아남은 아주 오래된 나무지. 하지만 오랜 세월을 지나는 동안 나무는 분노하기 시작했지. 숲을 파괴하는 인간을 보면서 분노했고, 메말라가는 숲을 보면서 분노했지.", {
			"그래서 어떻게 되었나요?",
		})
		if sel == nil then
			return
		end

		me:dialog(npc, "결국 나무의 분노는 나무를 몬스터로 바꾸어 놓고 말았고, 이제는 닥치는 대로 땅의 양분을 갉아먹는 한낱 괴물이 되어버렸지. 너무 깊이 알려고 하지 말게. 자네의 호기심은 이해하지만, 그는 모든 스텀프들의 왕이야. 결코 쉽게 생각하면 안된다네.", false, true)

		q:start(npc, true)
		q:force_complete(npc)
		me:show_effect(EffectType.QuestCompletion)
		me:exp(me:exp() + 100)
	end
}
