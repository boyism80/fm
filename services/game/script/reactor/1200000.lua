-- Reactor name (Reactor.wz/1200000.img.xml): 가짜바트

function on_reactor_1200000(reactor)
	local trigger = reactor:trigger()
	if trigger == nil then
		return
	end
	trigger:notice('가짜 바트를 때렸습니다.', Msg.PinkText)
	trigger:map(120000104)
end
