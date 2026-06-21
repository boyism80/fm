-- Mob name (String.wz/Mob.img.xml): 카오스 혼테일

local sponge_revive = require("script/lib/sponge_revive")

function on_revive_8810120(mob, map, x, y, link_oid, revives)
	sponge_revive.migrate_with_sponge(map, mob, revives, sponge_revive.HT_MIGRATE_SPONGE_IDS)
end
