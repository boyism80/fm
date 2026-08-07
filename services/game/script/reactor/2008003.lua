-- Reactor name (Reactor.wz/2008003.img.xml): 여신상 조각

return {
	on_reactor = function(reactor)
		local map = reactor:map()
		if map == nil then
			return
		end
		local player = reactor:trigger()
		if player == nil then
			return
		end
		local sm = player:state_machine()
		if sm == nil then
			return
		end
		local status = tonumber(sm:get_property("status")) or 0
		sm:set_property("status", tostring(status + 1))
		local minerva = map:reactor_by_name("minerva")
		if minerva ~= nil then
			minerva:hit(minerva:state() + 1)
		end
	end
}
