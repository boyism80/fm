-- 하늘위의 섬 크리세 (Quest.wz/QuestInfo.img.xml): 하늘위의 섬 크리세

local quest_id = 31000

function on_end(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	if not me:dialog(npc, "이제 오는거야? 얼마나 기다렸는지 모른다고.", false, true) then
		return
	end
	if not me:dialog(npc, "오르비스의 하늘 위에는 크리세라고 불리우는 천상의 섬이 있어. 거기엔 덩치는 큰 거인들이지만 착하고 여린 심성을 가진 거인족이 살고 있고. 그런데 얼마전 부터 크리세가 점점 멀어지기 시작하더니 연락이 되질 않고 있어. 무슨 일이 생긴게 분명한데... 마음 같아서는 당장이라도 찾아가 보고 싶지만 너도 알다시피 난 여기를 비울 수 없는 입장이라...", true, true) then
		return
	end
	if not me:dialog(npc, "그래서 하는 말인데 네가 크리세에 무슨일이 생긴건 아닌지 확인해줄래? 그럼 내가 크리세로 갈 수 있도록 도와줄게. 다녀와서 무슨일인지 꼭 나에게 전달해줘.", true, true) then
		return
	end
	if not me:dialog_accept(npc, "출발할 준비는 모두 된거야? 먼 길이 될 테니 착실히 준비해 가는게 좋을거야. 지금 보내 줄게.") then
		return
	end
	if not me:dialog(npc, "좋았어. 지금 바로 보내줄게. 쉽지 않은 여정이 될 수 있으니 마음 단단히 먹으라고.", false, true) then
		return
	end

	q:force_complete(npc)
	me:show_effect(EffectType.QuestCompletion)
	me:map(200100001, 0)
	me:notice("점프 키를 누르면 하늘을 날아 크리세에 도착 할 수 있습니다.", -1)
end
