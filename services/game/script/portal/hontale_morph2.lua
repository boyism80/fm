local DRAGON_MORPH = 4

return {
	on_enter = function(me)
		me:play_portal_sound()
		if me:morph() == DRAGON_MORPH then
			me:morph(false)
		end
		me:map(240040600, "east00")
	end
}
