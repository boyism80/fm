return {
	on_enter = function(me)
		local sm = me:state_machine()
		if sm ~= nil and sm:get_property("stage5") == "2" then
			me:play_portal_sound()
			me:map(926100300)
		else
			me:notice("지금은 포탈이 닫혀있습니다.", Msg.PinkText)
		end
	end
}
