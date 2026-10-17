local quest_id = 2234

return {
	on_end = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end
		if not q:started() then
			q:start(npc, true)
			return
		end

		local family = me:family()
		if family.total_reputation < 2000 or family.reputation > 500 then
			me:dialog(npc, "흐음, 아직 #b이름난 모험가의 특권#k을 제대로 맛보지 못했는가? 총 명성도를 2,000 이상 만들고, 현재 명성도를 500 이하로 만들면 되네.")
			return
		end
		me:dialog(npc, "드디어 자네도 #b이름난 모험가의 특권#k을 맛보았는가? 후후훗, 남은 명성도를 보니 충분히 맛을 본 모양이군. 그나저나 이렇게 빠른 시간 내에 이 정도까지 해잴 줄은.. 정말 자네의 리더쉽은 타고난 듯 싶네.", false, true)
		me:dialog(npc, "그럼, 앞으로도 좋은 엘더가 되어 주게나. 패밀리에 대해 궁금한 점이 있거나 칭호에 관심이 생긴다면 언제든지 나를 찾아오게.", false, true)

		if me:exchange(nil, { exp = 3000 }) ~= ExchangeResult.OK then
			return
		end
		q:force_complete(npc)
		me:show_effect(EffectType.QuestCompletion)
	end
}
