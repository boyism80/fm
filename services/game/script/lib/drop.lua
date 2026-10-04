local M = {}

local function roll(drops, owner)
	local spawned = {}
	for _, drop in ipairs(drops) do
		local quest_ok = true
		if drop.quest ~= nil then
			quest_ok = false
			if owner ~= nil then
				local quest = owner:quest(drop.quest)
				quest_ok = quest ~= nil and quest:started()
			end
		end
		if quest_ok and math.random(1000000) <= drop.chance then
			local count = drop.min
			if drop.max > drop.min then
				count = math.random(drop.min, drop.max)
			end
			if drop.meso then
				spawned[#spawned + 1] = { meso = count }
			else
				spawned[#spawned + 1] = { item = drop.item, count = count }
			end
		end
	end
	return spawned
end

function M.from_mob(drops, mob, attacker, map)
	if attacker == nil or map == nil then
		return
	end
	local spawned = roll(drops, attacker)
	if #spawned == 0 then
		return
	end
	mob:drop(spawned, attacker)
end

function M.from_reactor(drops, reactor)
	local map = reactor:map()
	if map == nil then
		return
	end
	local trigger = reactor:trigger()
	local spawned = roll(drops, trigger)
	if #spawned == 0 then
		return
	end
	local x, y = reactor:position()
	map:drop(spawned, { x, y }, trigger)
end

return M
