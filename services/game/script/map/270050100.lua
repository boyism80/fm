local KEY = "pinkbean_first_enter"

function on_map_enter_270050100(me, map)
	if map:property(KEY) then
		return
	end

	map:respawn(true)
	map:spawn_npc(2141000, -190, -42)
	map:property(KEY, true)
end
