-- Mob name (String.wz/Mob.img.xml): 혼테일의 머리A

function on_revive_8810002(mob, map, x, y, revives)
	for i = 1, #revives do
		local id = revives[i]
		if id ~= nil and id ~= 0 then
			map:spawn_mob(id, x, y)
		end
	end
end
