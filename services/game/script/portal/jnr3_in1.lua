return {
	on_enter = function(me)
		local map = me:map()
		if map == nil then
			return
		end
		local door = map:reactor_by_name("jnr32_out")
		if door ~= nil and door:state() > 0 then
			me:play_portal_sound()
			me:map(926110202)
		else
			me:notice("지금은 포탈이 닫혀있습니다.", Msg.PinkText)
		end
	end
}
