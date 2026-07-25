-- Reactor name (Reactor.wz/2221002.img.xml): 나무구멍:메밀묵을 떨어뜨려서 초록왕도깨비를 소환한다.

return {
	on_reactor = function(reactor)
		local map = reactor:map()
		if map == nil then
			return
		end
		map:spawn_mob(7130402, -340, 100)
		map:message('초록색 왕도깨비가 나타났습니다.')
	end
}
