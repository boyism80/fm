local pq = require("script/lib/party_quest")

local DEST = 920011000

local function map_count(group, map_id)
	local map = group:map(map_id)
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
			me:map(920011200)
			return
		end
		local group = sm:group()
		if group == nil then
			return
		end
		if pq.is_leader(me) then
			pq.warp_portal(me, DEST, "st00")
			sm:notice("파티장이 <암흑의 방> 에 입장하였습니다.")
			return
		end
		if map_count(group, DEST) > 0 then
			pq.warp_portal(me, DEST, "st00")
		else
			me:notice("파티장이 먼저 입장해야 합니다.", Msg.PinkText)
		end
	end
}
