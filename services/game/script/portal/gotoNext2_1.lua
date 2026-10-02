return {
	on_enter = function(me)
		local q3140 = me:quest(3140)
		local q3182 = me:quest(3182)
		if (q3140:started() == false and q3140:completed() == false) and (q3182:started() == false and q3182:completed() == false) then
			me:message("길이 막혀있습니다.", Msg.PinkText)
			return
		end
		me:map(211060500, 1)
	end
}
