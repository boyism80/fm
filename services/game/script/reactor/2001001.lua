-- Reactor name (Reactor.wz/2001001.img.xml): 파파픽시의 화분

return {
	on_reactor = function(reactor)
		local map = reactor:map()
		if map == nil then
			return
		end
		local x, y = reactor:position()
		map:spawn_mob(9300049, x, y - 10)
	end
}
