-- Reactor name (Reactor.wz/9208001.img.xml): 샤레니안 성문 석상

return {
	on_reactor = function(reactor)
		local map = reactor:map()
		if map == nil then
			return
		end
		local wz = map:wz()
		if wz == nil or wz.id ~= 910210000 then
			return
		end
		local trigger = reactor:trigger()
		if trigger == nil then
			return
		end
		local sm = trigger:state_machine()
		if sm == nil then
			return
		end
		local guess = sm:get_property("guess") or ""
		sm:set_property("guess", guess .. reactor:name())
	end
}
