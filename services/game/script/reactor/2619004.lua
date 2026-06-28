-- Reactor name (Reactor.wz/2619004.img.xml): 호문쿨루스 퇴치용1

function on_reactor_2619004(reactor)
	local map = reactor:map()
	if map == nil then
		return
	end
	for _, mob in pairs(map:mobs(6090004, 1)) do
		mob:kill()
	end
end
