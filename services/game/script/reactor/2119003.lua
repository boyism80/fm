-- Reactor name (Reactor.wz/2119003.img.xml): 좀비의 무덤+상자4

function on_reactor_2119003(reactor)
	local map = reactor:map()
	if map == nil then
		return
	end
	for _, mob in pairs(map:mobs(6090000, 1)) do
		mob:kill()
	end
end
