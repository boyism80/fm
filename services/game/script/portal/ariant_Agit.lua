return {
	on_enter = function(me)
		if me:quest(3930):completed() == false or me:quest(3933):completed() == false or me:quest(3936):completed() == false then
			me:message("이 문은 잠겨 있는 것 같다.")
			return
		end
		me:play_portal_sound()
		me:message("문 안쪽에서 잠금쇠가 열리는 소리가 들린다. 문이 조금 열렸다.")
		me:map(260000201, 1)
	end
}
