-- Reactor name (Reactor.wz/2006000.img.xml): 빛덩어리

return {
	on_reactor = function(reactor)
		local map = reactor:map()
		if map == nil then
			return
		end
		local x, y = reactor:position()
		map:spawn_npc(2013001, x, y - 10)
		local player = reactor:trigger()
		if player == nil then
			return
		end
		local sm = player:state_machine()
		if sm ~= nil then
			sm:set_property("prestage", "clear")
		end
	end
}
