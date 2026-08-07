-- Portal name: elevator (핼리오스탑 2층/99층)

return {
	on_enter = function(me)
		local group = state_machine("elevator")
		if group == nil then
			me:notice("엘리베이터가 고장났습니다.", Msg.PinkText)
			return
		end
		local map = me:map()
		if map == nil then
			return
		end
		local reactor = map:reactor_by_name("elevator")
		if reactor ~= nil and reactor:state() == 0 then
			me:map(map:wz():id() + 10)
			return
		end
		me:notice("엘리베이터 문이 닫혀있습니다.", Msg.PinkText)
	end,
}
