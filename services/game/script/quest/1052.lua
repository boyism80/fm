local quest_id = 1052

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	if not me:dialog_yes_no(npc, '도적으로 전직할 레벨이 되셨군요. 도적은 커닝시티에서 전직하실 수 있으며, 특별히 지금은 제가 도적으로 전직하실 수 있도록 커닝시티로 이동시켜 드릴까요?') then
		return
	end

	me:map(103000000)
	q:start(npc, true)
	q:force_complete(npc)
end
