local pq = require("script/lib/party_quest")

local OUT = {
	[920010200] = { portal = 4, name = "<산책로>" },
	[920010300] = { portal = 12, name = "<창고>" },
	[920010400] = { portal = 5, name = "<휴게실>" },
	[920010500] = { portal = 13, name = "<봉인된 방>" },
	[920010600] = { portal = 15, name = "<라운지>" },
	[920010700] = { portal = 14, name = "<올라가는 길>" },
	[920011000] = { portal = 16, name = "<암흑의 방>" },
}

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
		if not pq.is_leader(me) then
			me:notice("파티장이 먼저 이곳에서 퇴장해야 합니다.", Msg.PinkText)
			return
		end
		local sm = me:state_machine()
		if sm == nil then
			return
		end
		local map = me:map()
		if map == nil or map:wz() == nil then
			return
		end
		local map_id = map:wz().id
		local info = OUT[map_id]
		if info == nil then
			return
		end
		local group = sm:group()
		if map_id == 920010600 and group ~= nil then
			for id = 920010601, 920010604 do
				if map_count(group, id) > 0 then
					me:notice("객실에 파티원이 입장하고 있어 지금은 퇴장할 수 없습니다.")
					return
				end
			end
		end
		pq.party_warp(sm, 920010100, nil, info.portal)
		sm:notice("파티장이 " .. info.name .. " 에서 퇴장하였습니다.")
		me:play_portal_sound()
	end
}
