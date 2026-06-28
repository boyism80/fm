function on_reactor_2229009(reactor)
	local map = reactor:map()
	if map == nil then
		return
	end
	for _, mob in pairs(map:mobs(6090003, 1)) do
		mob:kill()
	end
end
