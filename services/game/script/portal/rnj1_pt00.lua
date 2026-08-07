local DEST = 926100001
local PROP = "stage1"
local NEED = "1"

return {
	on_enter = function(me)
		local sm = me:state_machine()
		if sm ~= nil and sm:get_property(PROP) == NEED then
			me:play_portal_sound()
			me:map(DEST)
		else
			me:notice("지금은 포탈이 닫혀있습니다.", Msg.PinkText)
		end
	end
}
