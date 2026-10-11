-- Item name (String.wz/Cash.img.xml): 칠판 : 1일권

return {
	on_cash = function(me, item_id, text)
		me:set_chalkboard(text)
		return false
	end,
}
