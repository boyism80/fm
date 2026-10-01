-- Portal (old/scripts/portal/guild1F04.js): 미로의 끝

local pq = require("script/lib/party_quest")
local gq = require("script/lib/guild_quest")

return {
	on_enter = function(me)
		gq.set_maze(me, "3")
		pq.warp_portal(me, 990000700, "st00")
		local sm = me:state_machine()
		if sm ~= nil then
			gq.gain_gp_once(me, sm, "gainGP04", 5)
		end
	end
}
