-- Mob name (String.wz/Mob.img.xml): 핑크빈

local pinkbean = require("script/lib/pinkbean")

function on_revive_8820014(mob, map, x, y, revives)
	local mobs = pinkbean.mobs_by_id(map)
	local chair = mobs[8820000]
	if chair == nil then
		return
	end

	local pkb = map:spawn_mob(8820001, x, y, MobSpawnType.Revive, chair:oid())
	if pkb == nil then
		return
	end

	chair:kill(MobDieAnimation.FadeOut)
end
