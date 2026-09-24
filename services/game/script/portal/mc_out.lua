local FALLBACK_MAP = 103000000

return {
	on_enter = function(me)
		local map_id = me:saved_location("MONSTERCARNIVAL")
		if map_id == nil then
			map_id = FALLBACK_MAP
		end
		me:clear_saved_location("MONSTERCARNIVAL")
		me:map(map_id, 0)
	end
}
