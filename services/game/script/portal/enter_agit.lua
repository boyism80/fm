return {
	on_enter = function(me)
		if me:quest(31001):completed() == false then
			me:message("들어가기엔 자격이 부족하다.")
			return
		end
		me:play_portal_sound()
		me:map(200100010, 1)
	end
}
