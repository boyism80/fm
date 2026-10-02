return {
	on_enter = function(me)
		local q31306 = me:quest(31306)
		if q31306:started() == false and q31306:completed() == false then
			me:message("길이 막혀있습니다.", Msg.PinkText)
			return
		end
		me:map(223010000, 1)
	end
}
