-- Item name (String.wz/Cash.img.xml): 펫작명하기

return {
	on_cash = function(me, item_id, text, ear, pet_sn)
		return me:rename_pet(text)
	end,
}
