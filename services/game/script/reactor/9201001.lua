-- Reactor name (Reactor.wz/9201001.img.xml): 샤렌3세의 유골:NPC소환

function on_reactor_9201001(reactor)
	local map = reactor:map()
	if map == nil then
		return
	end
	map:message('밝은 빛과 함께, 누군가가 나타났습니다.')
	local x, y = reactor:position()
	map:spawn_npc(9040003, x, y - 10)
end
