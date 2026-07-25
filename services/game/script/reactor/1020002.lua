-- Reactor name (Reactor.wz/1020002.img.xml): 91020000

return {
	on_reactor = function(reactor)
		local trigger = reactor:trigger()
		if trigger == nil then
			return
		end
		trigger:map(910200000, 3)
	end
}
