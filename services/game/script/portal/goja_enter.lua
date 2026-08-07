local GROUP_NAME = "gojarani"

return {
	on_enter = function(me)
		local group = state_machine(GROUP_NAME)
		if group == nil then
			me:notice("지금은 입장할 수 없습니다.", Msg.PinkText)
			return
		end
		local state = group:get_property("state")
		if state ~= nil and state ~= "" and state ~= "0" then
			me:notice("이미 다른 사람이 안에 있습니다.", Msg.PinkText)
			return
		end
		me:play_portal_sound()
		local sm, err = group:start_solo(me)
		if sm == nil then
			me:notice("지금은 입장할 수 없습니다.", Msg.PinkText)
			if err ~= nil then
				log("gojarani start_solo:", err)
			end
		end
	end,
}
