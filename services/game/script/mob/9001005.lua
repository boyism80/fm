-- Mob name (String.wz/Mob.img.xml): 파이렛 옥토

local QUEST_ID = 2192
local DROP_ID = 4031857

return {
	on_mob_die = function(mob, attacker, map)
		if attacker == nil or map == nil then
			return
		end
		local quest = attacker:quest(QUEST_ID)
		if quest == nil or not quest:started() then
			return
		end
		if math.random(1000000) >= 400000 then
			return
		end
		local x, y = mob:position()
		map:spawn_item(DROP_ID, 1, { x, y }, attacker)
	end
}
