-- Reactor name (Reactor.wz/2519001.img.xml): 데비존의 문

return {
	on_reactor = function(reactor)
		local map = reactor:map()
		if map == nil then
			return
		end
		map:block_gen(false, 9300121)
		map:message("문이 잠겼습니다. 해적들이 더 이상 나오지 않습니다.")
		local trigger = reactor:trigger()
		if trigger == nil then
			return
		end
		local sm = trigger:state_machine()
		if sm == nil then
			return
		end
		local n = tonumber(sm:get_property("stage4")) or 0
		sm:set_property("stage4", tostring(n + 1))
	end
}
