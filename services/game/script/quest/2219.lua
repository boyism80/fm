local quest_id = 2219

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		local name = me:name()
		me:dialog(npc, "어! 안녕하세요! 지하철로 들어가실 거면 먼저 표를 구입해 주세요! 표를 구입하지 않고서는 지하철로 들어가실 수 없습니다... 응? 왜 그러고 계시죠? 무슨 볼 일이라도...? 그러고보니 왠지 익숙한 이름인데... " .. name .. "... " .. name .. "...? 아!", false, true)
		me:dialog(npc, "당신이 제이엠이 말하던 바로 그 사람이로군요! 반갑습니다! 커닝시티를 위험에서 구하기 위해 애쓰고 계신다면서요! 정말 대단해요! 수고하시라는 의미에서 지하철을 공짜로 태워드리고 싶지만... 에, 그건 불가능하고... 아무튼 정말 훌륭하십니다!", false, true)
		me:dialog(npc, "그래서, 음모는 어디까지 밝혀내셨나요? 커닝시티를 위험에 빠뜨리다니... 도대체 누가 그런 짓을 하고 있는 건지 알아내셨나요? 다른 지역의 주민? 오시리아 대륙에서 온 괴인? 그것도 아니면 커닝시티의 누군가가? 아아! 정말 궁금해요!", false, true)
		if not me:dialog_accept(npc, "앗... 아직 당신께 정보를 안 드렸군요. 그럼 지금까지 모은 정보를 알려드릴게요. 자, 정보를 들을 준비는 되셨죠? 하긴 당연히 되셨겠죠! 커닝시티의 음모도 밝히실 분인데!") then
			me:dialog(npc, "정보를 들을 준비가 되시면 다시 말을 걸어주세요.")
			return
		end

		me:dialog(npc, "혹시 #b성장의 마법약#k이라는 걸 아세요? 먹으면 누구든지간에 무조건 덩치가 몇 배로 커져버린다는 약! 그게 바로 성장의 마법약이래요! 소문에 의하면, 슬라임이 이 약을 먹고 킹슬라임이 된 거라고 주장하는 사람도 있다고 해요!", false, true)
		me:dialog(npc, "넷...? 정보는 이게 다냐고요? 네... 이게 다인데요. 죄송해요. 정보가 좀 빈약하죠? 인적 드문 지하철에서 간간히 다니는 모험가들을 붙잡고 수다 떨다 얻은 정보라서 좀 그래요. 대충 이해해 주세요. 그럼 꼭 음모 밝혀내세요!", false, true)

		q:start(npc, true)
		q:force_complete(npc)
		me:show_effect(EffectType.QuestCompletion)
	end
}
