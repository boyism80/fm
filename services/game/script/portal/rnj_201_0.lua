local pq = require("script/lib/party_quest")

return {
	on_enter = function(me)
		me:play_portal_sound()
		pq.warp(me, 926100200)
	end
}
