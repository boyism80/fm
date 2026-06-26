local M = {}

function M.relink_sponge(map, dying_sponge, new_sponge)
	if map == nil or dying_sponge == nil or new_sponge == nil then
		return
	end
	for _, mob in pairs(map:mobs()) do
		if mob:oid() ~= new_sponge:oid() then
			local linked = mob:parent()
			if linked == dying_sponge then
				mob:parent(new_sponge)
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
