-- Item name (String.wz/Cash.img.xml): 칠판 : 7일권

return {
	on_cash = function(me, item_id, text)
		if me:forbidden("chalkboard") then
			me:message("이곳에서는 사용할 수 없습니다.", Msg.PinkText)
			return false
		end
		me:set_chalkboard(text)
		return false
	end,
}
