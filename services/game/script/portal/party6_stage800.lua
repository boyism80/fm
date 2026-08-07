local pq = require("script/lib/party_quest")

local PURPLE_STONE = 4001163
local MONSTER_MARBLE = 4001169
local PURIFY_BEAD = 2270004
local ENTRY_MAP = 300030100

return {
	on_enter = function(me)
		pq.remove_all(PURPLE_STONE, me)
		pq.remove_all(MONSTER_MARBLE, me)
		pq.remove_all(PURIFY_BEAD, me)
		me:map(ENTRY_MAP)
	end
}
