return {
	on_enter = function(me)
		local map = me:map()
		if map == nil then
			return
		end
		local door = map:find_reactor_name("rnj3_out1")
		if door ~= nil and door:state() > 0 then
			me:play_portal_sound()
			me:map(926100201)
		else
			me:message("지금은 포탈이 닫혀있습니다.", Msg.PinkText)
		end
	end
}
