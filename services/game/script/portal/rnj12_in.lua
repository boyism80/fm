local pq = require("script/lib/party_quest")

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
			return
		end
		local group = sm:group()
		if group == nil then
			return
		end
		local here = me:map()
		local n = 0
		if here ~= nil then
			for _, _ in pairs(here:characters()) do
				n = n + 1
			end
		end
		if n == 4 or map_count(group, 926100401) > 0 then
			pq.party_warp(sm, 926100401)
			sm:notice("유레테가 기계장치를 조작하자 거대한 괴물이 나타났다. 유레테는 기분나쁘게 웃으며 사라졌다.")
		else
			me:notice("파티원 전원이 이곳에 모여있지 않습니다.", Msg.PinkText)
		end
	end
}
