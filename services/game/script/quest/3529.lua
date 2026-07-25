-- 잃어버린 추억을 찾아 (Quest.wz/QuestData/3529.img): 잃어버린 추억을 찾아

local quest_id = 3529

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		me:dialog(npc, "안정된 자세와 기운, 고귀한 자태. 완벽한 기사가 되었군요. #h0#, 오랜만입니다..", false, true)
		me:dialog(npc, "처음 봤을 땐 현지도 제대로 적응도 못하던 초보자였는데, 어느새 이렇게 강해졌군요... 하긴. 당신이라면 잘해낼 줄 알았습니다. 하지만, 예상했던 일이라도 기쁜걸요?", false, true)
		me:dialog(npc, "계속 수련에 정진하도록 해주세요. 당신을 기사단에 들인 사람으로써 확신하는데, 당신은 여기서 그칠 사람이 아닙니다. 더 강력한 기사단 으로 거듭할 수 있겠죠.. 그 때를 기다리겠습니다.", false, true)

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
