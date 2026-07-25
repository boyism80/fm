-- Reactor name (Reactor.wz/2200000.img.xml): 221024400

return {
	on_reactor = function(reactor)
		local trigger = reactor:trigger()
		if trigger == nil then
			return
		end
		trigger:notice('알 수 없는 힘에 의해 바깥으로 쫓겨났습니다.', Msg.PinkText)
		trigger:map(221024400, 4)
	end
}
