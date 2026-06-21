local M = {}

function M.relink_sponge(map, dying_sponge, new_sponge)
	if map == nil or dying_sponge == nil or new_sponge == nil then
		return
	end
	local dying_oid = dying_sponge:oid()
	for _, mob in pairs(map:mobs()) do
		if mob:oid() ~= new_sponge:oid() then
			local linked = mob:linked_sponge()
			if linked == dying_sponge or mob:spawn_link() == dying_oid then
				mob:set_sponge(new_sponge)
			end
		end
	end
end

function M.clear_map_after(map, delay_ms)
	sleep(delay_ms or 3000)
	if map == nil then
		return
	end
	map:kill_all_mobs(1)
end

return M
