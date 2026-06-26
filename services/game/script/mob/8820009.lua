-- Mob name (String.wz/Mob.img.xml): set0 투명몹

local pinkbean = require("script/lib/pinkbean")

function on_revive_8820009(mob, map, x, y, revives)
	local mobs = pinkbean.mobs_by_id(map)
	local marker = mobs[8820020]
	if marker == nil then
		return
	end

	local sponge = map:spawn_mob(8820010, x, y)
	if sponge == nil then
		return
	end

	local part = map:spawn_mob(8820003, x, y, MobSpawnType.Revive, marker:oid())
	if part == nil then
		return
	end

	part:parent(sponge)
	marker:kill(MobDieAnimation.FadeOut)
end
