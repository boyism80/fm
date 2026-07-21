-- Mob name (String.wz/Mob.img.xml): 호위용 멧돼지

function on_mob_period_drop_9300102(mob, count)
	local map = mob:map()
	if map ~= nil then
		map:message(string.format("멧돼지가 %d번째 페로몬 샘플을 떨어뜨렸습니다.", count))
	end
	return 4031507
end
