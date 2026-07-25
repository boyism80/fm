-- Reactor name (Reactor.wz/1020000.img.xml): 910200000

return {
	on_reactor = function(reactor)
		local trigger = reactor:trigger()
		if trigger == nil then
			return
		end
		trigger:map(910200000, 1)
	end
}
