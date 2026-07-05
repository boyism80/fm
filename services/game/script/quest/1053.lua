local quest_id = 1053

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	if not me:dialog_yes_no(npc, '해적으로 전직할 레벨이 되셨군요. 해적은 노틸러스선에서 전직하실 수 있으며, 특별히 지금은 제가 해적으로 전직하실 수 있도록 노틸러스선으로 이동시켜드릴까요?') then
		return
	end

	me:map(120000000)
	q:force_start(npc)
	q:force_complete(npc)
end
