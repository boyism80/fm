-- Mob name (String.wz/Mob.img.xml): 핑크빈

local pinkbean = require("script/lib/pinkbean")

function on_revive_8820010(mob, map, x, y, revives)
	local mobs = pinkbean.mobs_by_id(map)
	local marker1 = mobs[8820024]
	if marker1 == nil then
		return
	end
	local marker2 = mobs[8820021]
	if marker2 == nil then
		return
	end

	local sponge = map:spawn_mob(8820011, x, y)
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

	marker1:kill(MobDieAnimation.FadeOut)
	marker2:kill(MobDieAnimation.FadeOut)
end
