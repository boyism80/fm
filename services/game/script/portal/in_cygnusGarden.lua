return {
	on_enter = function(me)
		if me:level() < 170 then
			me:message("170 레벨 이상만 입장 가능합니다.", Msg.PinkText)
			return
		end
		me:map(271040210, 0)
	end
}
