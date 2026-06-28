-- Reactor name (Reactor.wz/2406000.img.xml): 나인스피릿의둥지

function on_reactor_2406000(reactor)
	local map = reactor:map()
	if map == nil then
		return
	end
	map:message('알에서 아기용이 태어났습니다.')
	local x, y = reactor:position()
	map:spawn_npc(2081008, x, y - 10)
end
