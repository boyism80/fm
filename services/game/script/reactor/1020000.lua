function on_reactor_1020000(reactor)
	local trigger = reactor:trigger()
	if trigger == nil then
		return
	end
	trigger:map(910200000, 1)
end
