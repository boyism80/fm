-- Mob name (String.wz/Mob.img.xml): 핑크빈

local sponge_revive = require("script/lib/sponge_revive")

function on_revive_8820012(mob, map, x, y, link_oid, revives)
	sponge_revive.spawn_with_sponge(map, x, y, revives, sponge_revive.PB_SPAWN_SPONGE_IDS)
end
