-- Portal (old/scripts/portal/guild1F00.js): 샤렌 3세의 무덤

local pq = require("script/lib/party_quest")
local gq = require("script/lib/guild_quest")

local MAZE_ENDS = {
	["1"] = 990000620,
	["2"] = 990000631,
	["3"] = 990000641,
}

return {
	on_enter = function(me)
		pq.warp_portal(me, MAZE_ENDS[gq.maze(me)] or 990000611, "st00")
	end
}
