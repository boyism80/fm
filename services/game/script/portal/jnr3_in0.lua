local pq = require("script/lib/party_quest")

return {
	on_enter = function(me)
		local map = me:map()
		if map == nil then
			return
		end
		local door = map:find_reactor_name("jnr3_out1")
		if door ~= nil and door:state() > 0 then
			me:play_portal_sound()
			pq.warp(me, 926110201)
		else
			me:message("지금은 포탈이 닫혀있습니다.", Msg.PinkText)
		end
	end
}
