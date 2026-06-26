-- Mob name (String.wz/Mob.img.xml): 핑크빈

local pinkbean = require("script/lib/pinkbean")

function on_revive_8820013(mob, map, x, y, revives)
	local mobs = pinkbean.mobs_by_id(map)
	local marker1 = mobs[8820024]
	if marker1 == nil then
		return
	end
	local marker2 = mobs[8820025]
	if marker2 == nil then
		return
	end
	local marker3 = mobs[8820026]
	if marker3 == nil then
		return
	end
	local marker4 = mobs[8820027]
	if marker4 == nil then
		return
	end
	local marker5 = mobs[8820019]
	if marker5 == nil then
		return
	end

	local sponge = map:spawn_mob(8820014, x, y)
	if sponge == nil then
		return
	end

	local part1 = map:spawn_mob(8820015, x, y, MobSpawnType.Revive, marker1:oid())
	if part1 == nil then
		return
	end
	part1:parent(sponge)

	local part2 = map:spawn_mob(8820016, x, y, MobSpawnType.Revive, marker2:oid())
	if part2 == nil then
		return
	end
	part2:parent(sponge)

	local part3 = map:spawn_mob(8820017, x, y, MobSpawnType.Revive, marker3:oid())
	if part3 == nil then
		return
	end
	part3:parent(sponge)

	local part4 = map:spawn_mob(8820018, x, y, MobSpawnType.Revive, marker4:oid())
	if part4 == nil then
		return
	end
	part4:parent(sponge)

	local part5 = map:spawn_mob(8820002, x, y, MobSpawnType.Revive, marker5:oid())
	if part5 == nil then
		return
	end
	part5:parent(sponge)

	marker1:kill(MobDieAnimation.FadeOut)
	marker2:kill(MobDieAnimation.FadeOut)
	marker3:kill(MobDieAnimation.FadeOut)
	marker4:kill(MobDieAnimation.FadeOut)
	marker5:kill(MobDieAnimation.FadeOut)
end
