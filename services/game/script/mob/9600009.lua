-- Mob name (String.wz/Mob.img.xml): 대왕지네

local DROPS = {
	{ meso = true, min = 1553, max = 3156, chance = 1000000 },
	{ item = 4031227, min = 1, max = 1, chance = 1000000, quest = 4103 },
}

return {
	on_mob_die = function(mob, attacker, map)
		if attacker == nil or map == nil then
			return
		end
		local spawned = {}
		for _, drop in ipairs(DROPS) do
			local quest_ok = true
			if drop.quest ~= nil then
				local quest = attacker:quest(drop.quest)
				quest_ok = quest ~= nil and quest:started()
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
		local x, y = mob:position()
		map:drop(spawned, { x, y }, attacker)
	end
}
