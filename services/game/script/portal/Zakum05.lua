local zakum_battle = require("script/state_machine/zakum_battle")

return {
	on_enter = function(me)
		if zakum_battle.gate_closed() then
			me:message("이미 자쿰과의 전투가 시작되어 입장하실 수 없습니다.", Msg.PinkText)
			return
		end
		me:play_portal_sound()
		me:map(211042400, "west00")
	end
}
