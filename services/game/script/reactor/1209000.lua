-- Reactor name (Reactor.wz/1209000.img.xml): 진짜바트

return {
	on_reactor = function(reactor)
		local trigger = reactor:trigger()
		if trigger == nil then
			return
		end
		trigger:notice("진짜 바트를 찾았습니다.", Msg.PinkText)
		local quest = trigger:quest(116400)
		if quest ~= nil then
			if quest:started() then
				quest:record("q22")
			else
				quest:start("q22")
			end
		end
		trigger:map(120000104)
	end
}
