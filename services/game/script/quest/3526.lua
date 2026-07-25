-- 잃어버린 추억을 찾아 (Quest.wz/QuestData/3526.img): 잃어버린 추억을 찾아

local quest_id = 3526

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		me:dialog(npc, "존재하지 않는 것 처럼 극도로 적은 기척... 하지만 사실 너 같은 사람이야말로 강력한 도적이란 걸 이 다크로님은 알고 있지. 오랜만이다, #h0#.", false, true)
		me:dialog(npc, "제법 성장했군. 이 다크로드님과 비교해도 뒤지지 않겠는걸. 예전에는 기척도 숨길 줄 모르는 초보자 꼬마였는데... 훗. 하긴 시간이 많이 흘렀으니까. 이렇게 강해진 모습을 보니 기분 묘하군... 이런 게 자랑스럽다는 거겠지.", false, true)
		me:dialog(npc, "더 정진하도록 해. 널 도적으로 만든 사람으로써 확신하는데, 넌 여기서 그칠 사람이 아니야. 더 강력한 도적이 될 수 있어. 그때까지 더 노력하도록.", false, true)

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
