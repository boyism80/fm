local EXIT_MAP = 120000101
local CLEAR_MAP = 912010200

return {
	on_enter = function(me)
		local sm = me:state_machine()
		if sm == nil then
			me:map(EXIT_MAP)
			return
		end
		if sm:time_left() >= 120000 then
			me:notice("카이린과 싸워 2분 이상 버텨야 합니다.", Msg.PinkText)
			return
		end
		sm:finish(0)
		me:map(CLEAR_MAP)
	end
}
