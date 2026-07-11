local quest_id = 2216

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	if not me:dialog_accept(npc, "무슨 일인가해? 자물쇠 따기라면 이 몽땅따에게 모두 맡기라해~ 허어? 뒷골목의 제이엠이 보내서 왔다는 말인가해? " .. me:name() .. "... 흠. 정말인 모양이군. 좋아. 그럼 지금까지 모은 정보를 전해 주도록 하지.") then
		me:dialog(npc, "정보를 들을 준비가 되시면 다시 말을 걸게.")
		return
	end

	me:dialog(npc, "늪지대를 여행하던 모험가의 말에 의하면, 다일은 보통의 리게이터와 꼭 닮았다고 하더군. 마치 리게이터를 #r부풀려 놓은 것처럼#k 거의 똑같이 생겼다는 거야. 하지만 리게이터와 달리 마법까지 사용할 수 있어 훨씬 무섭다더군.", false, true)

	q:start(npc, true)
	q:force_complete(npc)
	me:show_effect(EffectType.QuestCompletion)
end
