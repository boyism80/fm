-- Reactor name (Reactor.wz/2201003.img.xml): 루디브리엄 범용리엑터

function on_reactor_2201003(reactor)
	local map = reactor:map()
	if map == nil then
		return
	end
	map:message('알리샤르가 나타났습니다!')
	map:spawn_mob(9300012, 970, 150)
end
