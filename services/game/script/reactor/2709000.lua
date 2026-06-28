-- Reactor name (Reactor.wz/2709000.img.xml): 핑크빈 파워업

function on_reactor_2709000(reactor)
	local map = reactor:map()
	if map == nil then
		return
	end
	map:spawn_mob(8820008, 8, -53)
end
