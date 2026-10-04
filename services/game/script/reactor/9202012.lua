-- Reactor name (Reactor.wz/9202012.img.xml): 보물상자&lt;보너스>

local DROPS = {
	{ item = 2290113, min = 1, max = 1, chance = 3600 },
	{ item = 2290118, min = 1, max = 1, chance = 4400 },
	{ item = 2290120, min = 1, max = 1, chance = 5400 },
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
