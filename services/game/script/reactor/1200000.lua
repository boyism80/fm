function on_reactor_1200000(reactor)
	local trigger = reactor:trigger()
	if trigger == nil then
		return
	end
	trigger:notice(5, '가짜 바트를 때렸습니다.')
	trigger:map(120000104)
end
