local quest_id = 1051

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if not me:dialog_yes_no(npc, '궁수로 전직할 레벨이 되셨군요. 궁수는 헤네시스에서 전직하실 수 있으며, 특별히 지금은 제가 궁수로 전직하실 수 있도록 헤네시스로 이동시켜드릴 수 있습니다. 지금 이동시켜 드릴까요?') then
			return
		end

		me:map(100000000)
		q:start(npc, true)
		q:force_complete(npc)
	end
}
