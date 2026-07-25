-- Reactor name (Reactor.wz/2401002.img.xml): 아이스이글의 알

return {
	on_reactor = function(reactor)
		local map = reactor:map()
		if map == nil then
			return
		end
		local x, y = reactor:position()
		map:spawn_mob(9300090, x, y - 10)
	end
}
