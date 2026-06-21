-- Mob name (String.wz/Mob.img.xml): 핑크빈

local sponge_revive = require("script/lib/sponge_revive")

function on_revive_8820014(mob, map, x, y, link_oid, revives)
	sponge_revive.spawn_revives_animate(map, x, y, revives)
end
