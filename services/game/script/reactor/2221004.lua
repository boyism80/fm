-- Reactor name (Reactor.wz/2221004.img.xml): 놀부네 지붕:놀부의 박씨를 떨어뜨려서 박을 소환한다.

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
