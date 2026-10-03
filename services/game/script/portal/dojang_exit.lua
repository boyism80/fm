local DOJO_ENTRANCE = 250000100

return {
	on_enter = function(me)
		local dest = me:saved_location("MULUNG_TC") or DOJO_ENTRANCE
		me:clear_saved_location("MULUNG_TC")
		me:play_portal_sound()
		me:map(dest, 0)
	end
}
