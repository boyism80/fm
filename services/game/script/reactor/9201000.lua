-- Reactor name (Reactor.wz/9201000.img.xml): 투명리엑터:서브몬스터 소환

return {
	on_reactor = function(reactor)
		local map = reactor:map()
		if map == nil then
			return
		end
		for _ = 1, 8 do
			map:spawn_mob(9300033, -100, 50)
		end
	end
}
