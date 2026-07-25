-- Mob name (String.wz/Mob.img.xml): 혼테일의 오른손

return {
	on_revive = function(mob, map, x, y, revives)
		for i = 1, #revives do
			local id = revives[i]
			if id ~= nil and id ~= 0 then
				map:spawn_mob(id, x, y)
			end
		end
	end
}
