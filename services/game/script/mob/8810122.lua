-- Mob name (String.wz/Mob.img.xml): 카오스 혼테일

local horntail = require("script/lib/horntail")

function on_mob_die_8810122(mob, attacker, map)
	horntail.clear_map_after(map, 3000)
end

function on_revive_8810122(mob, map, x, y, revives)
end
