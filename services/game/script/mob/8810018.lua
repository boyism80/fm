-- Mob name (String.wz/Mob.img.xml): 혼테일

local sponge_revive = require("script/lib/sponge_revive")

function on_mob_die_8810018(mob, attacker, map)
	sponge_revive.clear_horntail_map_after(map, 3000)
end
