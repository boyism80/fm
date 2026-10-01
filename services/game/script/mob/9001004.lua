-- Mob name (String.wz/Mob.img.xml): 카이린의 분신

local BLACK_CHARM = 4031059

return {
	on_mob_die = function(mob, attacker, map)
		if attacker == nil or map == nil then
			return
		end
		local x, y = mob:position()
		map:spawn_item(BLACK_CHARM, 1, { x, y }, attacker)
	end
}
