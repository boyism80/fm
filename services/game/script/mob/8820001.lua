-- Mob name (String.wz/Mob.img.xml): 핑크빈

local pinkbean = require("script/lib/pinkbean")

return {
	on_mob_die = function(mob, attacker, map)
		pinkbean.clear_map_after(map, 3000)
	end
}
