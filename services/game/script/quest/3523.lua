-- 잃어버린 추억을 찾아 (Quest.wz/QuestData/3523.img): 잃어버린 추억을 찾아

local quest_id = 3523

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	me:dialog(npc, "...이 강력한 위압감을 보니 자네는 대단히 강력한 전사로군. 이름이... #h0#? 하하! 그러고 보니 기억에 있는 이름이군. 아주 오래 전에 전사가 되겠다고 찾아왔던 미숙한 초보자들 중에, 그런 이름이 있었지.", false, true)
	me:dialog(npc, "자네가 이렇게 강력한 전사가 될 줄이야! 이제는 이 주먹펴고일어서라고 해도 승부를 장담하기 힘들겠는걸? 정말 대단하네... 그래, 자네라면 훌륭한 전사로 성장할 줄 알았네.", false, true)
	me:dialog(npc, "계속 더 정진하게. 자네를 전사로 만들어준 사람으로써 확신한다네. 더 강한 전사가 될 거란 걸...", false, true)

	q:start(npc, true)
	q:force_complete(npc)
	me:show_effect(EffectType.QuestCompletion)
	me:show_quest_completion(3507)
	local q7081 = me:quest(7081)
	if q7081 ~= nil then
		q7081:start(npc, "1")
	end
end
