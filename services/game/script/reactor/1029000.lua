function on_reactor_1029000(reactor)
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
