-- Mob name (String.wz/Mob.img.xml): 혼테일의 왼손

function on_revive_8810005(mob, map, x, y, link_oid, revives)
	for i = 1, #revives do
		local id = revives[i]
		if id ~= nil and id ~= 0 then
			map:spawn_mob(id, x, y, -3, link_oid)
		end
	end
end
