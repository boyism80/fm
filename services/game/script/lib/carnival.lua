local class = require("script/lib/class")

local M = {}

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

return M
