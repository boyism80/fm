-- Portal (old/scripts/portal/tamepig_out2.js): 사육실 통로

local pq = require("script/lib/party_quest")

local PHEROMONE = 4031507
local RESEARCH_REPORT = 4031508

return {
	on_enter = function(me)
		pq.remove_all(PHEROMONE, me)
		pq.remove_all(RESEARCH_REPORT, me)
		me:play_portal_sound()
		me:map(230000003)
	end
}
