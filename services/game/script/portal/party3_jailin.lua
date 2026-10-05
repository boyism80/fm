local pq = require("script/lib/party_quest")

return {
	on_enter = function(me)
		local map = me:map()
		if map == nil or map:wz() == nil then
			return
		end
		pq.warp(me, map:wz():id() + 2)
	end
}
