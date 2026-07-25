-- Mob name (String.wz/Mob.img.xml): 혼테일 소환

return {
	on_revive = function(mob, map, x, y, revives)
		local spawned_mobs = mob:revive(revives, x, y)
		local sponge_mobs = spawned_mobs[8810018]
		if sponge_mobs == nil then
			return
		end
		local sponge = sponge_mobs[1]
		if sponge == nil then
			return
		end
		for _, mobs in pairs(spawned_mobs) do
			for _, part in ipairs(mobs) do
				if part ~= sponge then
					part:parent(sponge)
				end
			end
		end
	end
}
