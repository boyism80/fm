-- Reactor name (Reactor.wz/2201004.img.xml): 균열조각을 떨어뜨리면 시간의 구를 소환한다.

return {
	on_reactor = function(reactor, item)
		if item ~= nil then
			return
		end
		local map = reactor:map()
		if map == nil then
			return
		end
		local group = state_machine("papulatus")
		if group == nil then
			return
		end
		if group:get_property("battle") == "1" then
			return
		end
		local sm = group:get("Battle")
		if sm == nil then
			local player = reactor:trigger()
			if player ~= nil then
				sm = player:state_machine()
			end
		end
		if sm == nil then
			return
		end
		sm:start()
		map:message("시공의 균열이 <차원 균열의 조각> 으로 메꾸어 졌습니다.")
		map:music("Bgm09/TimeAttack")
		map:spawn_mob(8500000, -410, -400)
	end
}
