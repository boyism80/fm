return {
	on_enter = function(me)
		local q3143 = me:quest(3143)
		if q3143:started() == false and q3143:completed() == false then
			me:message("길이 막혀있습니다.", Msg.PinkText)
			return
		end
		me:map(211060010, 0)
	end
}
