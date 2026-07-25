-- Item name (String.wz/Consume.img.xml): 아대 공격력 주문서 60%

return {
	on_scroll = function(me, scroll_slot, target_slot)
		return me:enhance(scroll_slot, target_slot, EquipmentPart.Weapon, WeaponType.Claw, function(target, scroll)
		    target:add_bonus_stats(scroll:wz():bonus_stats())
		end)
	end
}
