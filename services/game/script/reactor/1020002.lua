function on_reactor_1020002(reactor)
	local trigger = reactor:trigger()
	if trigger == nil then
		return
	end
	trigger:map(910200000, 3)
end
