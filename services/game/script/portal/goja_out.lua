local EXIT_MAP = 200000204

return {
	on_enter = function(me)
		me:play_portal_sound()
		local sm = me:state_machine()
		if sm ~= nil then
			sm:finish(EXIT_MAP)
			return
		end
		me:map(EXIT_MAP)
		me:notice("의사양반이 있는 안락한 병원을 떠나 버섯똘이에게 돌아갑니다.", Msg.PinkText)
	end,
}
