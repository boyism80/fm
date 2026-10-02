return {
	on_enter = function(me)
		local q3141 = me:quest(3141)
		if q3141:started() == false and q3141:completed() == false then
			me:message("길이 막혀있습니다.", Msg.PinkText)
			return
		end
		me:map(211060700, 1)
	end
}
