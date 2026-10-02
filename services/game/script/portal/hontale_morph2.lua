local DRAGON = 2210003

return {
	on_enter = function(me)
		me:play_portal_sound()
		me:map(240040600, "east00")
		if me:morph() == DRAGON then
			me:morph(false)
		end
	end
}
