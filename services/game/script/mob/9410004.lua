-- Mob name (String.wz/Mob.img.xml): 폭주족 원숭이

local DROPS = {
	{ item = 4000202, min = 1, max = 1, chance = 500000 },
	{ item = 4031296, min = 1, max = 1, chance = 50000, quest = 4010 },
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
