-- Reactor name (Reactor.wz/2229009.img.xml): 선비귀신 퇴치용

return {
	on_reactor = function(reactor)
		local map = reactor:map()
		if map == nil then
			return
		end
		for _, mob in pairs(map:mobs(6090003, 1)) do
			mob:kill()
		end
	end
}
