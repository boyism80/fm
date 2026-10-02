return {
	on_enter = function(me)
		local q31303 = me:quest(31303)
		if q31303:started() == false and q31303:completed() == false then
			me:message("길이 막혀있습니다.", Msg.PinkText)
			return
		end
		me:map(223010110, 1)
	end
}
