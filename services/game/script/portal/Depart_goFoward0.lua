return {
	on_enter = function(me)
		me:play_portal_sound()
		me:map(me:map():wz():id() + 10, "right00")
	end
}
