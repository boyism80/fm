-- Item name (String.wz/Cash.img.xml): (+)소유권 없는 아이템&메소 획득 스킬

return {
	on_cash = function(me, item_id, text, ear, pet_sn)
		return me:change_pet_skill(pet_sn, item_id)
	end,
}
