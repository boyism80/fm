local quest_id = 2218

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if not me:dialog_accept(npc, "뭐야? 길이라도 잃었어? 그럼 왜 그렇게 쳐다보는 건데? ... 헤에. 그러고 보니 제이엠이 이름이 " .. me:name() .. "인 사람을 보면 정보를 건네주라고 했었지? 좋아. 수집한 정보를 너한테 말해줄게.") then
			me:dialog(npc, "정보를 들을 준비가 되면 다시 말을 걸어줘.")
			return
		end

		me:dialog(npc, "혹시 마법사 #r라케리스#k를 알아? 저쪽, 하수구 앞에서 모험가들을 모으는 여자 말이야. 요즘 그녀가 왠지 이상해. 어딘가 불안해 보인다고 할까, 수상해 보인다고 할까... 창백한 얼굴로 하수구 쪽을 볼 때도 많고.", false, true)
		me:dialog(npc, "뭐, 지레짐작한 것일 뿐일지도 모르지만, 어쨌든 그 여자는 #r요주의 인물#k이니까 다시 한 번 확인해 보는 게 좋을 거라고 제이엠한테 전해줘.", false, true)

		q:start(npc, true)
		q:force_complete(npc)
		me:show_effect(EffectType.QuestCompletion)
	end
}
