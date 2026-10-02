return {
	on_enter = function(me)
		local dest = me:saved_location("AMORIA") or me:map():wz():return_map_id()
		me:clear_saved_location("AMORIA")
		me:play_portal_sound()
		me:map(dest, 0)
	end
}
