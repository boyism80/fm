return {
	on_enter = function(me)
		local sm = me:state_machine()
		if sm == nil then
			return
		end
		local key = "stage6_2"
		if sm:get_property(key) == "0" then
			me:play_portal_sound()
			me:map(926110303)
			sm:set_property(key, "1")
		else
			me:notice("이미 누군가가 이 포탈 안에 들어가 있습니다.", Msg.PinkText)
		end
	end
}
