-- Mob name (String.wz/Mob.img.xml): 혼테일의 머리A

local sponge_revive = require("script/lib/sponge_revive")

function on_revive_8810002(mob, map, x, y, link_oid, revives)
	sponge_revive.spawn_revives_linked(map, x, y, revives, link_oid)
end
