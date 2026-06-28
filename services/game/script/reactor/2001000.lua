-- Reactor name (Reactor.wz/2001000.img.xml): 여신의 화분:여신탑의 네펜데스 소환

function on_reactor_2001000(reactor)
	local map = reactor:map()
	if map == nil then
		return
	end
	local x, y = reactor:position()
	map:spawn_mob(9300048, x, y - 10)
end
