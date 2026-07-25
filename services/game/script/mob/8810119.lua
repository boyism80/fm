-- Mob name (String.wz/Mob.img.xml): 카오스 혼테일

local horntail = require("script/lib/horntail")

return {
	on_mob_die = function(mob, attacker, map)
		local new_sponge_id = 8810120
		if map:mob_by_template(new_sponge_id) ~= nil then
			return
		end
		local x, y = mob:position()
		local new_sponge = map:spawn_mob(new_sponge_id, x, y)
		if new_sponge == nil then
			return
		end
		horntail.relink_sponge(map, mob, new_sponge)
	end,

	on_revive = function(mob, map, x, y, revives)
	end
}
