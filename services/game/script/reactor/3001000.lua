-- Reactor name (Reactor.wz/3001000.img.xml): 괴인의 제단

function on_reactor_3001000(reactor)
	local map = reactor:map()
	if map == nil then
		return
	end
	map:message('포이즌 골렘이 나타났습니다.')
	local x, y = reactor:position()
	map:spawn_mob(9300180, x, y - 10)
end
