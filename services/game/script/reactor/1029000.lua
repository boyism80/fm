-- Reactor name (Reactor.wz/1029000.img.xml): 엘리니아 꽃 리엑터

return {
	on_reactor = function(reactor)
		local map = reactor:map()
		if map == nil then
			return
		end
		for _, mob in pairs(map:mobs(3230300)) do
			mob:kill()
		end
		for _, mob in pairs(map:mobs(3230301)) do
			mob:kill()
		end
	end
}
