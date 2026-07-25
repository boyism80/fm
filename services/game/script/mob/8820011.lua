-- Mob name (String.wz/Mob.img.xml): 핑크빈

local pinkbean = require("script/lib/pinkbean")

return {
	on_revive = function(mob, map, x, y, revives)
		local mobs = pinkbean.mobs_by_id(map)
		local marker1 = mobs[8820024]
		if marker1 == nil then
			return
		end
		local marker2 = mobs[8820025]
		if marker2 == nil then
			return
		end
		local marker3 = mobs[8820022]
		if marker3 == nil then
			return
		end

		local sponge = map:spawn_mob(8820012, x, y)
		if sponge == nil then
			return
		end

		local part1 = map:spawn_mob(8820003, x, y, MobSpawnType.Revive, marker1:oid())
		if part1 == nil then
			return
		end
		part1:parent(sponge)

		local part2 = map:spawn_mob(8820004, x, y, MobSpawnType.Revive, marker2:oid())
		if part2 == nil then
			return
		end
		part2:parent(sponge)

		local part3 = map:spawn_mob(8820005, x, y, MobSpawnType.Revive, marker3:oid())
		if part3 == nil then
			return
		end
		part3:parent(sponge)

		marker1:kill(MobDieAnimation.FadeOut)
		marker2:kill(MobDieAnimation.FadeOut)
		marker3:kill(MobDieAnimation.FadeOut)
	end
}
