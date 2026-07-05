local quest_id = 1050

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	if not me:dialog_yes_no(npc, '마법사로 전직할 레벨이 되셨군요. 마법사는 엘리니아에서 전직하실 수 있으며, 특별히 지금은 제가 마법사로 전직하실 수 있도록 엘리니아로 이동시켜드릴 수 있습니다. 지금 이동시켜 드릴까요?') then
		return
	end

	me:map(101000000)
	q:force_start(npc)
	q:force_complete(npc)
end
