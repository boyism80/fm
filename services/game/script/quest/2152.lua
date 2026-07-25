local quest_id = 2152

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		local sel = me:dialog_list(npc, "어서와요. 용건이 뭐죠?", {
			"귀신나무에 대해 아시는 것이 있나요?",
		})
		if sel == nil then
			return
		end

		sel = me:dialog_list(npc, "윈스턴 박사님의 연구를 도와 주고 있나 보군요? 글쎄요. 저도 박사님의 부탁을 받고 조사를 좀 해봤는데 알아낸 것이 없어요. 단지 요즘 페리온과의 접경지대에 있는 엘리니아의 숲이 급속도로 메말라가기 시작했다는 것을 알아냈죠. 진행속도는 느리지만 경계해야 할 일이에요.", {
			"네. 시간을 내주셔서 감사합니다.",
		})
		if sel == nil then
			return
		end

		me:dialog(npc, "많은 도움이 되지 못한 것 같아서 미안하군요.", false, true)

		q:start(npc, true)
		q:force_complete(npc)
		me:show_effect(EffectType.QuestCompletion)
		me:exp(me:exp() + 200)
	end
}
