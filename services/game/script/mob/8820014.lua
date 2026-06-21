-- Mob name (String.wz/Mob.img.xml): 핑크빈

function on_revive_8820014(mob, map, x, y, link_oid, revives)
	for i = 1, #revives do
		local id = revives[i]
		if id ~= nil and id ~= 0 then
			map:spawn_mob(id, x, y, -2)
		end
	end
end
