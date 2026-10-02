return {
	on_enter = function(me)
		me:play_portal_sound()
		me:map(103000301, 0)
		me:warp_later(103000100, 20)
	end
}
