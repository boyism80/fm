-- Item name (String.wz/Consume.img.xml): 귀 장식 행운 주문서 70%

return {
	on_scroll = function(me, scroll_slot, target_slot)
		return me:enhance(scroll_slot, target_slot, EquipmentPart.Eye, nil, function(target, scroll)
		    target:add_bonus_stats(scroll:wz():bonus_stats())
		end)
	end
}
