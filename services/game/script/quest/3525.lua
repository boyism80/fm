-- 잃어버린 추억을 찾아 (Quest.wz/QuestData/3525.img): 잃어버린 추억을 찾아

local quest_id = 3525

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		me:dialog(npc, "... 이 기운은.. 날카로운 매의 눈빛이로군요. 어서오세요, #h0#. 활을 다룰줄도 모르고 화살도 무서워 하던 풋풋한 초보자였던 당신이 여기까지 성장할 줄이야.", false, true)
		me:dialog(npc, "당신이 이렇게 강력한 궁수가 될 줄은 이미 느끼고 있었답니다.", false, true)
		me:dialog(npc, "계속 더 정진하세요. 자네를 궁수로 만들어준 사람으로써 확신해요. 더 강한 궁수가 될 거란 걸...", false, true)

		q:start(npc, true)
		q:force_complete(npc)
		me:show_effect(EffectType.QuestCompletion)
		me:show_quest_completion(3507)
		local q7081 = me:quest(7081)
		if q7081 ~= nil then
			q7081:start(npc, "1")
		end
	end
}
