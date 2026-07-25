local quest_id = 3108

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		me:dialog(npc, "(조각상은 한 눈에 봐도 눈이 부실정도로 아름답다. 얼음으로 만들어진 것 같이 투명하지만 얼음은 아닌 것 같다. 조각상 주위를 돌며 자세히 살펴보았다.)", false, true)
		me:dialog(npc, "(조각상의 한쪽이 부셔져 있다. 주위에는 커다란 발자국도 몇 개 보인다.)", false, true)

		me:exp(me:exp() + 200)
		q:start(npc, true)
		q:force_complete(npc)
		me:show_effect(EffectType.QuestCompletion)
	end
}
