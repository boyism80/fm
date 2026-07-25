-- Reactor name (Reactor.wz/2221001.img.xml): 석등:곡차를 떨어뜨려서 파란왕도깨비를 소환한다.

return {
	on_reactor = function(reactor)
		local map = reactor:map()
		if map == nil then
			return
		end
		local x, y = reactor:position()
		map:spawn_mob(7130401, x, y - 10)
		map:message('파란색 왕도깨비가 나타났습니다.')
	end
}
