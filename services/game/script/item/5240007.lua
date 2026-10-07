-- Item name (String.wz/Cash.img.xml): 아이스 후르츠

return {
	on_cash = function(me, item_id, text, ear, pet_sn)
		return me:feed_pet_cash(item_id)
	end,
}
