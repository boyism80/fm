function on_reactor_2516000(reactor)
	local map = reactor:map()
	if map == nil then
		return
	end
	map:message('해적왕이 죽고 도라지 왕이 풀려났습니다!')
	local x, y = reactor:position()
	map:spawn_npc(2094001, x, y - 10)
end
