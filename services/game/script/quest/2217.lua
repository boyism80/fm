local quest_id = 2217

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if not me:dialog_accept(npc, "음... " .. me:name() .. "... 네가 제이엠이 말한 그 사람인 모양이네. 그럼 얼마전까지 모은 정보를 알려 줄게. 도움이 될지는 모르겠지만...") then
			me:dialog(npc, "정보를 들을 준비가 되면 다시 말을 걸어줘.")
			return
		end

		me:dialog(npc, "크리스가 그러는데, 얼마 전부터 하수도에서 이상한 냄새가 나기 시작했대. ...하수도야 원래 냄새가 이상하지만, 평소와 다른 #r뭔가가 섞여든 것 같다#k는데... 다일하고는 상관 없으려나? 아무튼 정보는 이게 다야.", false, true)

		q:start(npc, true)
		q:force_complete(npc)
		me:show_effect(EffectType.QuestCompletion)
	end
}
