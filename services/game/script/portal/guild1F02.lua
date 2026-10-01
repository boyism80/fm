-- Portal (old/scripts/portal/guild1F02.js): 수로의 미로

local pq = require("script/lib/party_quest")
local gq = require("script/lib/guild_quest")

return {
	on_enter = function(me)
		gq.set_maze(me, "1")
		pq.warp_portal(me, 990000700, "st00")
		local sm = me:state_machine()
		if sm ~= nil then
			gq.gain_gp_once(me, sm, "gainGP02", 5)
		end
	end
}
