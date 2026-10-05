-- Portal name: elevator (핼리오스탑 2층/99층)

return {
	on_enter = function(me)
		local group = state_machine("elevator")
		if group == nil then
			me:message("엘리베이터가 고장났습니다.", Msg.PinkText)
			return
		end
		local map = me:map()
		if map == nil then
			return
		end
		local reactor = map:find_reactor_name("elevator")
		local sm = group:get("persistent")
		if reactor ~= nil and reactor:state() == 0 and sm ~= nil then
			sm:enter_player(me, map:wz():id() + 10)
			return
		end
		me:message("엘리베이터 문이 닫혀있습니다.", Msg.PinkText)
	end,
}
