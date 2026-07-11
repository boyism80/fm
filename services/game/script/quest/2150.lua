local quest_id = 2150

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	local sel = me:dialog_list(npc, "안녕하세요 여행자님 오늘은 무슨 일로 오셨나요?", {
		"귀신나무에 대해 알고 있니?",
	})
	if sel == nil then
		return
	end

	sel = me:dialog_list(npc, "어머! 그 소문을 들으신거에요? 얼마 전에 헤네시스의 카밀라가 엄마 심부름으로 페리온에 왔다가 돌아가는 길에 귀신을 봤대요.", {
		"정말이니?",
	})
	if sel == nil then
		return
	end

	me:dialog(npc, "밤 늦게 헤네시스로 돌아가는 길이었는데 어둠속에서 나무줄기를 밟은 것 같아서 주위를 둘러보는데 희번덕거리는 눈이 카밀라를 잡아먹을 것처럼 쳐다봤다고 하더라구요.", false, true)
	me:dialog(npc, "카밀라는 너무 무서워서 그대로 기절하고 말았대요. 날이 밝은 뒤에 어른들이 그 자리에 다시 가봤는데 아무 것도 없었대요. 귀신이 분명한 것 같아요. 어쩌죠? 이제 무서워서 마을 밖에 나갈 수가 없을 것 같아요.", false, true)

	q:start(npc, true)
	q:force_complete(npc)
	me:show_effect(EffectType.QuestCompletion)
	me:exp(me:exp() + 100)
end
