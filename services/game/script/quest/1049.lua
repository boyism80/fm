local quest_id = 1049

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	if not me:dialog_yes_no(npc, '전사로 전직할 레벨이 되셨군요. 전사는 페리온에서 전직하실 수 있으며, 특별히 지금은 제가 전사로 전직하실 수 있도록 페리온으로 이동시켜드릴 수 있습니다. 지금 이동시켜 드릴까요?') then
		return
	end

	me:map(102000000)
	q:start(npc, true)
	q:force_complete(npc)
end
