-- Item name (String.wz/Consume.img.xml): 루돌프의 뿔 전용주문서 60%

return {
	on_scroll = function(me, scroll_slot, target_slot)
		return me:enhance(scroll_slot, target_slot, EquipmentPart.Cap, nil, function(target, scroll)
		    target:add_bonus_stats(scroll:wz():bonus_stats())
		end)
	end
}
