-- Reactor name (Reactor.wz/2709000.img.xml): 핑크빈 파워업

return {
	on_reactor = function(reactor)
		local map = reactor:map()
		if map == nil then
			return
		end
		map:spawn_mob(8820008, 8, -53)
	end
}
