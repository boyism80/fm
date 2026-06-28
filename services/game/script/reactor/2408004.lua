function on_reactor_2408004(reactor)
	local map = reactor:map()
	if map == nil then
		return
	end
	map:message('알에서 아기용이 태어났습니다.')
	local x, y = reactor:position()
	map:spawn_npc(2081008, x, y - 10)
end
