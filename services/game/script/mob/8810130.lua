-- Mob name (String.wz/Mob.img.xml): 카오스 혼테일 소환

function on_revive_8810130(mob, map, x, y, link_oid, revives)
	local sponge_id = 8810118
	local sponge = nil
	local parts = {}
	for i = 1, #revives do
		local id = revives[i]
		if id ~= nil and id ~= 0 then
			local spawned_mob = map:spawn_mob(id, x, y, -2)
			if spawned_mob ~= nil then
				if id == sponge_id then
					sponge = spawned_mob
					sponge:parent(sponge)
				else
					parts[#parts + 1] = spawned_mob
				end
			end
		end
	end
	if sponge == nil then
		return
	end
	for _, part in ipairs(parts) do
		part:parent(sponge)
	end
end
