local class = require("script/lib/class")

local M = {}

local HUB_MAP = 980000000

function M.challenge_info(members, size)
	if members == nil then
		members = {}
	end
	if size == nil then
		size = #members
	end
	if #members == 0 then
		return "파티 인원: " .. tostring(size)
	end
	local lines = { "#b" }
	for _, ch in ipairs(members) do
		lines[#lines + 1] = string.format(
			"%s / 레벨 : %d / 직업 : %s\r\n",
			ch:name(),
			ch:level(),
			class.name(ch:class())
		)
	end
	lines[#lines + 1] = "#k\r\n"
	return table.concat(lines)
end

function M.revive(me)
	local map = me:map()
	if map == nil or map:wz() == nil then
		return
	end
	local match = carnival.map_match(map:wz():id())
	if match == nil then
		return
	end
	local portal = "sp"
	local team = me:carnival_team()
	if team ~= nil then
		local team_id = team:id()
		if team_id == CARNIVAL_TEAM.RED then
			portal = "red_revive"
		elseif team_id == CARNIVAL_TEAM.BLUE then
			portal = "blue_revive"
		end
	end
	local sm = me:state_machine()
	if sm == nil then
		return
	end
	local field = sm:group():map(match:field_map_id())
	if field == nil then
		return
	end
	local target = field:portal(portal)
	if target == nil then
		return
	end
	me:map(field, target:id())
end

function M.leave(me)
	local sm = me:state_machine()
	if sm == nil then
		return
	end
	me:message("누군가가 나가기 엔피시를 클릭하여 모두 나가집니다.", Msg.Notice)
	local match = carnival.map_match(tonumber(sm:id()))
	sm:finish(HUB_MAP)
	if match ~= nil then
		match:finish()
	end
end

return M
