-- Reactor name (Reactor.wz/2201001.img.xml): 차원의 블록퍼스를 소환한다.

return {
	on_reactor = function(reactor)
		local map = reactor:map()
		if map == nil then
			return
		end
		local x, y = reactor:position()
		y = y - 10
		for _ = 1, 3 do
			map:spawn_mob(9300007, x, y)
		end
	end
}
