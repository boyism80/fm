-- Mob name (String.wz/Mob.img.xml): 월묘

function on_mob_period_drop_9300061(mob, count)
	local map = mob:map()
	if map ~= nil then
		map:message(string.format("월묘가 %d번째 떡을 만들었습니다.", count))
	end
	return 4001101
end
