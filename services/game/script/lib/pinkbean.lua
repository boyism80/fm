local M = {}

local SPONGE_DIE_MS = 4380

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

function M.advance_sponge(map, x, y, next_sponge_id, stages)
	sleep(SPONGE_DIE_MS)
	local mobs = M.mobs_by_id(map)
	local markers = {}
	for _, stage in ipairs(stages) do
		local marker = mobs[stage.marker]
		if marker == nil then
			return
		end
		markers[#markers + 1] = marker
	end

	local sponge = map:spawn_mob(next_sponge_id, x, y)
	if sponge == nil then
		return
	end

	for i, stage in ipairs(stages) do
		local part = map:spawn_mob(stage.part, x, y, MobSpawnType.Revive, markers[i]:oid())
		if part == nil then
			return
		end
		part:parent(sponge)
	end

	for _, marker in ipairs(markers) do
		marker:kill(MobDieAnimation.FadeOut)
	end
end

function M.clear_map_after(map, delay_ms)
	sleep(delay_ms or 3000)
	if map == nil then
		return
	end
	map:kill_all_mobs(MobDieAnimation.FadeOut)
end

return M
