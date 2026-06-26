-- Mob name (String.wz/Mob.img.xml): 핑크빈

local pinkbean = require("script/lib/pinkbean")

function on_mob_die_8820001(mob, attacker, map)
	pinkbean.clear_map_after(map, 3000)
end
