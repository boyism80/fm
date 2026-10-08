local dojo = require("script/lib/dojo")

return {
	on_enter = function(me)
		local sm = me:state_machine()
		local floor = dojo.current_floor(me)
		if sm == nil or floor == nil then
			return
		end
		if dojo.is_rest(floor) then
			me:message("소공에게 말을 걸어 도전을 계속하세요.", Msg.PinkText)
			return
		end
		if not dojo.cleared(sm, floor) then
			me:message("아직 몬스터가 남아있습니다.", Msg.PinkText)
			return
		end
		me:play_portal_sound()
		dojo.advance(sm, floor)
	end
}
