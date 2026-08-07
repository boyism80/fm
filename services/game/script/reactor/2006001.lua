-- Reactor name (Reactor.wz/2006001.img.xml): 여신의 석상

return {
	on_reactor = function(reactor)
		local map = reactor:map()
		if map == nil then
			return
		end
		local x, y = reactor:position()
		map:spawn_npc(2013002, x, y - 10)
		local player = reactor:trigger()
		if player == nil then
			return
		end
		local sm = player:state_machine()
		if sm ~= nil then
			sm:set_property("status", "7")
		end
	end
}
