local M = {}

function M.mobs_by_id(map)
	local by_id = {}
	if map == nil then
		return by_id
	end
	for _, mob in pairs(map:mobs()) do
		local id = mob:id()
		if by_id[id] == nil then
			by_id[id] = mob
		end
	end
	return by_id
end

function M.clear_map_after(map, delay_ms)
	sleep(delay_ms or 3000)
	if map == nil then
		return
	end
	map:kill_all_mobs(MobDieAnimation.FadeOut)
end

return M
