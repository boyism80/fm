return {
	on_enter = function(me)
		local dest = me:saved_location("MULUNG_TC") or me:map():wz():return_map_id()
		me:clear_saved_location("MULUNG_TC")
		me:play_portal_sound()
		me:map(dest, 0)
	end
}
