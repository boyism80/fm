-- 크세르크세스 퇴치하라 (Quest.wz/QuestInfo.img.xml): 크세르크세스 퇴치하라

local quest_id = 31013

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	if not me:dialog_accept(npc, "이제 준비가 완료되고 있는 것 같아요. 여행자님은 어떠세요?\r\n슬슬 때가 오지 않았나 생각하는데... #b크세르크세스#k를 퇴치해 주시겠습니까?") then
		me:dialog(npc, "아직도 준비가 부족한 걸까요?", false, false)
		return
	end

	q:start(npc, true)
end
