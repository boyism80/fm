local dojo = require("script/lib/dojo")

return {
	on_enter = function(me)
		local sm = me:state_machine()
		if sm == nil or sm:get_property("cleared") ~= "1" then
			me:message("아직 몬스터가 남아있습니다.", Msg.PinkText)
			return
		end
		me:play_portal_sound()
		sm:finish(dojo.LOBBY)
	end
}
