-- Mob name (String.wz/Mob.img.xml): 카오스 혼테일

local horntail = require("script/lib/horntail")

function on_mob_die_8810118(mob, attacker, map)
	local new_sponge_id = 8810119
	if map:mob_by_template(new_sponge_id) ~= nil then
		return
	end
	local x, y = mob:position()
	local new_sponge = map:spawn_mob(new_sponge_id, x, y)
	if new_sponge == nil then
		return
	end
	horntail.relink_sponge(map, mob, new_sponge)
end

function on_revive_8810118(mob, map, x, y, revives)
end
