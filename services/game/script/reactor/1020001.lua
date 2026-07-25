-- Reactor name (Reactor.wz/1020001.img.xml): 91020000

return {
	on_reactor = function(reactor)
		local trigger = reactor:trigger()
		if trigger == nil then
			return
		end
		trigger:map(910200000, 2)
	end
}
