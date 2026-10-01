local KEY = "pinkbean_first_enter"

return {
	on_map_enter = function(me, map)
		if map:property(KEY) then
			return
		end

		map:respawn({ include_one_time = true })
		map:spawn_npc(2141000, -190, -42)
		map:property(KEY, true)
	end
}
