-- Item name (String.wz/Cash.img.xml): 자동 줍기 스킬

return {
	on_cash = function(me, item_id, text, ear, pet_sn)
		return me:change_pet_skill(pet_sn, item_id)
	end,
}
