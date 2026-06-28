-- Reactor name (Reactor.wz/2111001.img.xml): 제단

local zakum_arm_ids = {
	8800003, 8800004, 8800005, 8800006,
	8800007, 8800008, 8800009, 8800010,
}

function on_reactor_2111001(reactor)
	local map = reactor:map()
	if map == nil then
		return
	end

	map:music('Bgm06/FinalFight')
	local zakum = map:spawn_mob(8800000, -10, -215, -2)
	zakum:fake(true)
	for _, mob_id in ipairs(zakum_arm_ids) do
		map:spawn_mob(mob_id, -10, -215, -2)
	end
	map:message('원석의 힘으로 자쿰이 소환됩니다.')
end
