local M = {}

M.SPAWN_ANIMATE = -2
M.SPAWN_REVIVE = -3

M.HT_SPAWN_SPONGE_IDS = {8810018, 8810118}
M.HT_MIGRATE_SPONGE_IDS = {8810119, 8810120, 8810121, 8810122}
M.PB_SPAWN_SPONGE_IDS = {8820009, 8820010, 8820011, 8820012, 8820013}

local function revives_to_list(revives)
	local list = {}
	for i = 1, #revives do
		local id = revives[i]
		if id ~= nil and id ~= 0 then
			list[#list + 1] = id
		end
	end
	return list
end

local function pick_sponge_id(revives, candidates)
	local chosen = nil
	for i = 1, #revives do
		local id = revives[i]
		for _, candidate in ipairs(candidates) do
			if id == candidate then
				chosen = id
			end
		end
	end
	return chosen
end

local function spawn_sponge_revives(map, x, y, revives, sponge_id)
	if map == nil or sponge_id == nil or sponge_id == 0 then
		return
	end
	if map:mob_by_template(sponge_id) ~= nil then
		return
	end
	local sponge = map:spawn_mob(sponge_id, x, y, M.SPAWN_ANIMATE)
	if sponge == nil then
		return
	end
	sponge:sponge(true)
	for _, id in ipairs(revives_to_list(revives)) do
		if id ~= sponge_id then
			local part = map:spawn_mob(id, x, y, M.SPAWN_ANIMATE)
			if part ~= nil then
				part:set_sponge(sponge)
			end
		end
	end
end

local function migrate_sponge(map, dying_mob, sponge_id)
	if map == nil or dying_mob == nil or sponge_id == nil or sponge_id == 0 then
		return
	end
	if map:mob_by_template(sponge_id) ~= nil then
		return
	end
	local x, y = dying_mob:position()
	local sponge = map:spawn_mob(sponge_id, x, y, M.SPAWN_ANIMATE)
	if sponge == nil then
		return
	end
	sponge:sponge(true)
	local dying_oid = dying_mob:oid()
	for _, mob in pairs(map:mobs()) do
		if mob:oid() ~= sponge:oid() then
			local linked = mob:linked_sponge()
			if linked == dying_mob or mob:spawn_link() == dying_oid then
				mob:set_sponge(sponge)
			end
		end
	end
end

function M.spawn_with_sponge(map, x, y, revives, sponge_candidates)
	local sponge_id = pick_sponge_id(revives, sponge_candidates)
	if sponge_id == nil then
		return
	end
	spawn_sponge_revives(map, x, y, revives, sponge_id)
end

function M.migrate_with_sponge(map, dying_mob, revives, sponge_candidates)
	local sponge_id = pick_sponge_id(revives, sponge_candidates)
	if sponge_id == nil then
		return
	end
	migrate_sponge(map, dying_mob, sponge_id)
end

function M.spawn_revives_animate(map, x, y, revives)
	if map == nil then
		return
	end
	for _, id in ipairs(revives_to_list(revives)) do
		map:spawn_mob(id, x, y, M.SPAWN_ANIMATE)
	end
end

function M.spawn_revives_linked(map, x, y, revives, link_oid)
	if map == nil then
		return
	end
	for _, id in ipairs(revives_to_list(revives)) do
		map:spawn_mob(id, x, y, M.SPAWN_REVIVE, link_oid)
	end
end

function M.clear_horntail_map_after(map, delay_ms)
	sleep(delay_ms or 3000)
	if map == nil then
		return
	end
	map:kill_all_mobs(1)
end

return M
