-- Item name (String.wz/Cash.img.xml): 바나나

return {
	on_cash = function(me, item_id, text, ear, pet_sn)
		return me:feed_pet_cash(item_id)
	end,
}
