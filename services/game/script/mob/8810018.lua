-- Mob name (String.wz/Mob.img.xml): 혼테일

local horntail = require("script/lib/horntail")

function on_mob_die_8810018(mob, attacker, map)
	horntail.clear_map_after(map, 3000)
end
