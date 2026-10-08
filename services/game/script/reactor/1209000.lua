-- Reactor name (Reactor.wz/1209000.img.xml): 진짜바트

return {
	on_reactor = function(reactor)
		local trigger = reactor:trigger()
		if trigger == nil then
			return
		end
		trigger:message("진짜 바트를 찾았습니다.", Msg.PinkText)
		trigger:records():set_text("air_strike.progress", "q22")
		trigger:map(120000104)
	end
}
