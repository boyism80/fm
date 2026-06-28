function on_reactor_2619003(reactor)
	local map = reactor:map()
	if map == nil then
		return
	end
	for _, mob in pairs(map:mobs(6090004, 1)) do
		mob:kill()
	end
end
