-- NPC name (String.wz/Npc.img.xml): 누리스

local GQ_ITEMS = { 1032033, 4001024, 4001025, 4001026, 4001027, 4001028, 4001029, 4001030, 4001031, 4001032, 4001033, 4001034, 4001035, 4001037 }

return {
	on_click = function(me, npc)
		local ear = me:equipped(EquipmentPart.Ear)
		if ear ~= nil and ear:wz():id() == 1032033 then
			me:dialog(npc, "#b#t1032033##k을 장비 해제해주세요.")
			return
		end

		local cost = {}
		local has_cost = false
		for _, item_id in ipairs(GQ_ITEMS) do
			local count = 0
			for _, it in pairs(me:item(item_id)) do
				count = count + it:count()
			end
			if count > 0 then
				cost[item_id] = count
				has_cost = true
			end
		end
		if has_cost then
			me:exchange({ item = cost }, nil)
		end
		me:map(101030104, 0)
	end
}
