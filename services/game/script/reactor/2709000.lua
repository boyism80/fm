function on_reactor_hit_2709000(reactor)
	local map = reactor:map()
	if map == nil then
		return
	end
	map:spawn_mob(8820008, 8, -53)
end
