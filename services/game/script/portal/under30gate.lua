-- Portal (old/scripts/portal/under30gate.js): 지하수로

return {
	on_enter = function(me)
		if me:level() > 30 then
			me:message("이곳은 지나갈 수 없습니다.", Msg.PinkText)
			return
		end
		me:play_portal_sound()
		me:map(990000640, 1)
	end
}
