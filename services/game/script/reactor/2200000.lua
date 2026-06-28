function on_reactor_2200000(reactor)
	local trigger = reactor:trigger()
	if trigger == nil then
		return
	end
	trigger:notice(5, '알 수 없는 힘에 의해 바깥으로 쫓겨났습니다.')
	trigger:map(221024400, 4)
end
