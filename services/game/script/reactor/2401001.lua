-- Reactor name (Reactor.wz/2401001.img.xml): 파이어호크의 알

function on_reactor_2401001(reactor)
	local map = reactor:map()
	if map == nil then
		return
	end
	local x, y = reactor:position()
	map:spawn_mob(9300089, x, y - 10)
end
