return {
	on_enter = function(me)
		if me:quest(3368):started() == false then
			me:message("지금은 이 연구실에 볼 일이 없다.", Msg.PinkText)
			return
		end
		me:play_portal_sound()
		local sm, err = state_machine("yurete_lab3"):start_solo(me)
		if sm == nil then
			log("yurete_lab3 start_solo:", err)
			me:message("이 안에 이미 다른 누군가가 들어가 있는 것 같다.", Msg.PinkText)
		end
	end
}
