-- Item name (String.wz/Cash.img.xml): 퀵배송 이용권

return {
	on_cash = function(me, item_id, text, ear)
		me:open_quick_delivery()
		return false
	end,
}
