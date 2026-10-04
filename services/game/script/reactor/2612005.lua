-- Reactor name (Reactor.wz/2612005.img.xml): 5색 비커

local DROPS = {
	{ item = 4031798, min = 1, max = 1, chance = 999999, quest = 3366 },
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
