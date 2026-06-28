-- Reactor name (Reactor.wz/2201000.img.xml): 퀘스트용장난감목마를 최대 10마리까지 소환

function on_reactor_2201000(reactor)
	local map = reactor:map()
	if map == nil then
		return
	end
	local x, y = reactor:position()
	y = y - 10
	local count = math.random(3, 10)
	for _ = 1, count do
		map:spawn_mob(9300011, x, y)
	end
end
