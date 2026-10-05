-- Reactor name (Reactor.wz/9208007.img.xml): 창의 제단

local HALL_MAP = 990000400

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
		local hall = sm:map(HALL_MAP)
		if hall == nil then
			return
		end
		local gate = hall:find_reactor_name("speargate")
		if gate == nil then
			return
		end
		local count = 0
		for _, altar in pairs(map:reactors()) do
			if altar:id() == reactor:id() and altar:state() >= 1 then
				count = count + 1
			end
		end
		gate:hit(count)
	end
}
