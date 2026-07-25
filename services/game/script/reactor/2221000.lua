-- Reactor name (Reactor.wz/2221000.img.xml): 돌무더기:돼지고기 산적을 떨어뜨려 노란왕도깨비를 소환한다.

return {
	on_reactor = function(reactor)
		local map = reactor:map()
		if map == nil then
			return
		end
		local x, y = reactor:position()
		map:spawn_mob(7130400, x, y - 10)
		map:message('노란색 왕도깨비가 나타났습니다.')
	end
}
