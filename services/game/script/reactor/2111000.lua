-- Reactor name (Reactor.wz/2111000.img.xml): 몬스터소환3

return {
	on_reactor = function(reactor)
		local map = reactor:map()
		if map == nil then
			return
		end
		local trigger = reactor:trigger()
		if trigger ~= nil then
			trigger:notice('몬스터가 소환되었습니다!', Msg.PinkText)
		end
		local x, y = reactor:position()
		y = y - 10
		for _ = 1, 3 do
			map:spawn_mob(9300004, x, y)
		end
	end
}
