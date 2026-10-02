return {
	on_enter = function(me)
		if me:level() < 50 then
			me:message("50 레벨 이상만 입장 가능합니다.", Msg.PinkText)
			return
		end
		me:map(950101000, 0)
	end
}
