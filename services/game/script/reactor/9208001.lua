-- Reactor name (Reactor.wz/9208001.img.xml): 악마형 석상

local gq = require("script/lib/guild_quest")

return {
	on_reactor = function(reactor)
		local map = reactor:map()
		if map == nil then
			return
		end
		if map:wz():id() ~= 910210000 then
			gq.record_statue(reactor)
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
