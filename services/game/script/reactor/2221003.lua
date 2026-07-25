-- Reactor name (Reactor.wz/2221003.img.xml): 흥부네 지붕:흥부의 박씨를 떨어뜨려 박을 소환한다.

return {
	on_reactor = function(reactor)
		local map = reactor:map()
		if map == nil then
			return
		end
		local x, y = reactor:position()
		map:spawn_mob(9500400, x, y - 10)
	end
}
