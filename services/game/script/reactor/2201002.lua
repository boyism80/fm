-- Reactor name (Reactor.wz/2201002.img.xml): 루디브리엄 범용리엑터

return {
	on_reactor = function(reactor)
		local map = reactor:map()
		if map == nil then
			return
		end
		map:message('어딘가에 차원의 롬바드가 나타났습니다!')
		map:spawn_mob(9300010, 42, -326)
	end
}
