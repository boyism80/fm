-- 잃어버린 추억을 찾아 (Quest.wz/QuestData/3524.img): 잃어버린 추억을 찾아

local quest_id = 3524

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	me:dialog(npc, "훌륭하게 정제된 마력이군. 대마법사의 반열에 오른 자들만이 보일 수 있는... 그러고 보니, 아주 오래 전에 대마법사의 재능이 있는 초보자를 본 적이 있었는데... 그 이름이 #h0#였지.", false, true)
	me:dialog(npc, "아직 매직클로도 쓸 줄 모르던 미숙한 초보자였던 자네가 이리 성장한 모습을 보게 되다니... 정말 기쁘구만. 자네라면 이렇게 될 줄 알았지.", false, true)
	me:dialog(npc, "계속 더 정진하게... 자네를 마법사로 만든 사람으로써, 확신하고 있다네. 자네가 더 강한 마법사가 되리라고...", false, true)

	q:start(npc, true)
	q:force_complete(npc)
	me:show_effect(EffectType.QuestCompletion)
	me:show_quest_completion(3507)
	local q7081 = me:quest(7081)
	if q7081 ~= nil then
		q7081:start(npc, "1")
	end
end
