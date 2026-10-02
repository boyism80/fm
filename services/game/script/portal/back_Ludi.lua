return {
	on_enter = function(me)
		if me:level() < 10 then
			me:message("10 레벨 이상만 입장 가능합니다.", Msg.PinkText)
			return
		end
		me:map(123456789, 0)
	end
}
