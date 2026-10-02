return {
	on_enter = function(me)
		local q4317 = me:quest(4317)
		if q4317:started() == false and q4317:completed() == false then
			me:message("길이 막혀있습니다.", Msg.PinkText)
			return
		end
		me:map(231030000, 1)
	end
}
