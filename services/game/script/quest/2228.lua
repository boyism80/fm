local quest_id = 2228

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	me:dialog(npc, "저주의 기운이 끊어진 것을 느꼈어... 파우스트를 퇴치해 준 건 역시 너겠지? 정말 고마워... 이젠 정말 마음이 가벼워졌어. 모두 네 덕분이야.", false, true)
	me:dialog(npc, "더 이상 과거에 연연해하지 않겠어. 유령에 불과하지만, 그래도 아직 뭔가 더 좋은 일을 할 수 있을 거라 믿어. 예를 들어... 이 숲을 정화할 방법을 찾아보는 것도 괜찮지 않을까?", false, true)

	me:population(me:population() + 3)
	q:start(npc, true)
	q:force_complete(npc)
	me:show_effect(EffectType.QuestCompletion)
end
