-- Channel donation progress (old/src/server/maps/AramiaFireWorks.java): the NPC takes items, a full goal starts script/state_machine/fireworks.lua.

local M = {}

local SUN_MAP = 970010000
local SCRIPT = "script/state_machine/fireworks.lua"
local TREE_STAGES = 4

local ITEMS = {
	keg = 4001128,
	sun = 4001246,
}

-- Takes the donation from the character and returns false when they do not have it.
function M.give(me, kind, count)
	local item = ITEMS[kind]
	if me:exchange({ item = { [item] = count } }, nil) ~= ExchangeResult.OK then
		return false
	end

	if fireworks_give(kind, count) then
		world_notice(string.format("<%d채널> 필요한 만큼의 %s이(가) 모두 모였습니다!", channel_id() + 1, item2name(item)))
		local sm, err = state_machine("fireworks"):start_persistent(kind, { kind = kind })
		if sm == nil then
			log("fireworks start:", err)
		end
		return true
	end

	if kind == "sun" then
		-- One stage per sixth of the goal, and the first sixth is where the progress starts.
		local stage = math.floor(fireworks_percentage(kind) * 6 / 10000) - 1
		run_on_map(SUN_MAP, SCRIPT, "set_tree_stage", math.max(0, math.min(TREE_STAGES, stage)))
	end
	return true
end

-- Hundredths of a percent become whole percent.
function M.percentage(kind)
	return fireworks_percentage(kind) / 100
end

return M
