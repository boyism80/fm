-- Item name (String.wz/Consume.img.xml): 한손검 공격력 주문서 100%

return {
	on_scroll = function(me, scroll_slot, target_slot)
		return me:enhance(scroll_slot, target_slot, EquipmentPart.Weapon, WeaponType.Sword1H, function(target, scroll)
		    target:add_bonus_stats(scroll:wz():bonus_stats())
		end)
	end
}
