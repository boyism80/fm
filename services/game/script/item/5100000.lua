-- Item name (String.wz/Cash.img.xml): 축하곡 오르골

return {
	on_cash = function(me, item_id)
		if me:map():play_song(item_id, me:name()) == false then
			me:message("이미 다른 곡이 흘러나오고 있습니다.", Msg.PinkText)
			return false
		end
		return true
	end,
}
