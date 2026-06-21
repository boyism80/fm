-- Mob name (String.wz/Mob.img.xml): 혼테일의 왼손

local sponge_revive = require("script/lib/sponge_revive")

function on_revive_8810005(mob, map, x, y, link_oid, revives)
	sponge_revive.spawn_revives_linked(map, x, y, revives, link_oid)
end
