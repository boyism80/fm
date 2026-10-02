return {
	on_enter = function(me)
		if me:quest(3139):completed() == false then
			me:message("길이 막혀있습니다.", Msg.PinkText)
			return
		end
		me:map(211060300, 2)
	end
}
