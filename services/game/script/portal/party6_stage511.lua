local pq = require("script/lib/party_quest")

local FAIL_INDEX = 12

return {
	on_enter = function(me)
		if math.random(1, 10) == 1 then
			pq.warp_portal(me, 930000300, "16st")
		else
			pq.warp_portal(me, 930000300, string.format("%02dst", FAIL_INDEX))
		end
	end
}
