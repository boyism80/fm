-- Reactor name (Reactor.wz/6802001.img.xml): Green Reg treasure box

local DROPS = {
	{ meso = true, min = 10, max = 49, chance = 999999 },
	{ meso = true, min = 10, max = 49, chance = 999999 },
	{ meso = true, min = 10, max = 49, chance = 999999 },
	{ meso = true, min = 10, max = 49, chance = 999999 },
	{ item = 2050004, min = 1, max = 1, chance = 250000 },
}

return {
	on_reactor = function(reactor)
		reactor:drop_items()
		local map = reactor:map()
		if map == nil then
			return
		end
		local trigger = reactor:trigger()
		local spawned = {}
		for _, drop in ipairs(DROPS) do
			local quest_ok = true
			if drop.quest ~= nil then
				quest_ok = false
				if trigger ~= nil then
					local quest = trigger:quest(drop.quest)
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
		if #spawned == 0 then
			return
		end
		local x, y = reactor:position()
		map:drop(spawned, { x, y }, trigger)
	end
}
