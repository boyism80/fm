-- Mob name (String.wz/Mob.img.xml): 애기보스 소환용 투명몹

function on_revive_8820008(mob, map, x, y, revives)
	mob:revive({ 8820009, 8820000 }, x, y)
end
