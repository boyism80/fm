local pq = require("script/lib/party_quest")

return {
	on_enter = function(me)
		local sm = me:state_machine()
		if sm ~= nil and sm:get_property("stage5") == "2" then
			me:play_portal_sound()
			pq.warp(me, 926110300)
		else
			me:message("지금은 포탈이 닫혀있습니다.", Msg.PinkText)
		end
	end
}
