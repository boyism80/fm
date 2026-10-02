return {
	on_enter = function(me)
		if me:quest(2238):completed() == false then
			me:message("아직 들어갈 수 없습니다.", Msg.PinkText)
			return
		end
		me:play_portal_sound()
		me:map(105100101, 1)
	end
}
