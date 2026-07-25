-- Reactor name (Reactor.wz/2119006.img.xml): 사냥꾼의 제단2

return {
	on_reactor = function(reactor)
		local map = reactor:map()
		if map == nil then
			return
		end
		for _, mob in pairs(map:mobs(6090001, 1)) do
			mob:kill()
		end
	end
}
