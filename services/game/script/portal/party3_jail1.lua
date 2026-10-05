local pq = require("script/lib/party_quest")
local function map_count(sm, map_id)
	local map = sm:map(map_id)
	if map == nil then
		return 0
	end
	local n = 0
	for _, _ in pairs(map:characters()) do
		n = n + 1
	end
	return n
end

return {
	on_enter = function(me)
		local sm = me:state_machine()
		if sm == nil then
			return
		end
		if map_count(sm, 920010910) > 0
			or map_count(sm, 920010911) > 0
			or map_count(sm, 920010912) > 0 then
			me:message("이미 감옥에 누군가가 들어가 있습니다.")
			return
		end
		pq.warp(me, 920010910)
	end
}
