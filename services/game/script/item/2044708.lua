-- Item name (String.wz/Consume.img.xml): [5주년]아대 공격력 주문서 40%

return {
	on_scroll = function(me, scroll_slot, target_slot)
		return me:enhance(scroll_slot, target_slot, EquipmentPart.Weapon, WeaponType.Claw, function(target, scroll)
		    target:add_bonus_stats(scroll:wz():bonus_stats())
		end)
	end
}
