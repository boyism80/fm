local M = {}

function M.item_count(me, id)
	if me == nil or id == nil then
		return 0
	end
	local n = 0
	for _, it in pairs(me:item(id)) do
		if it ~= nil then
			n = n + it:count()
		end
	end
	return n
end

function M.has_item(me, id, qty)
	if qty == nil then
		qty = 1
	end
	return M.item_count(me, id) >= qty
end

local function remove_one(me, id)
	if me == nil or id == nil then
		return false
	end
	local n = M.item_count(me, id)
	if n <= 0 then
		return true
	end
	return me:rmitem(id, n)
end

local function is_state_machine(target)
	return target ~= nil and target.players ~= nil and target.rmitem == nil
end

local function party_characters_on_map(anchor)
	local party = anchor:party()
	local map = anchor:map()
	if party == nil or map == nil then
		return { anchor }
	end
	local pid = party:id()
	local out = {}
	for _, ch in pairs(map:characters()) do
		if ch ~= nil then
			local p = ch:party()
			if p ~= nil and p:id() == pid then
				out[#out + 1] = ch
			end
		end
	end
	if #out == 0 then
		return { anchor }
	end
	return out
end

function M.remove_all_party(id, anchor)
	if id == nil or anchor == nil then
		return false
	end
	local ok = true
	for _, ch in ipairs(party_characters_on_map(anchor)) do
		if not remove_one(ch, id) then
			ok = false
		end
	end
	return ok
end

function M.remove_all(id, target)
	if id == nil then
		return false
	end
	if target == nil then
		return false
	end
	if is_state_machine(target) then
		local ok = true
		for _, p in ipairs(target:players()) do
			if p ~= nil and not remove_one(p, id) then
				ok = false
			end
		end
		return ok
	end
	return remove_one(target, id)
end

function M.gain_item(me, id, n)
	if n == nil then
		n = 1
	end
	if n <= 0 then
		return nil
	end
	return me:mkitem(id, n)
end

function M.shuffle(str)
	if str == nil or str == "" then
		return ""
	end
	local chars = {}
	for i = 1, #str do
		chars[i] = str:sub(i, i)
	end
	for i = #chars, 2, -1 do
		local j = math.random(i)
		chars[i], chars[j] = chars[j], chars[i]
	end
	return table.concat(chars)
end

function M.party_exp(sm, amount)
	if sm == nil or amount == nil or amount <= 0 then
		return
	end
	for _, p in ipairs(sm:players()) do
		if p ~= nil then
			p:exchange({}, { exp = amount })
		end
	end
end

function M.party_warp(sm, map_id, except_id, spawn)
	if sm == nil or map_id == nil then
		return
	end
	if spawn == nil then
		spawn = 0
	end
	for _, p in ipairs(sm:players()) do
		if p ~= nil and (except_id == nil or p:id() ~= except_id) then
			p:map(map_id, spawn)
		end
	end
end

function M.is_gm(me)
	if me == nil then
		return false
	end
	return me:role() == ROLE.Admin
end

function M.is_leader(me)
	local party = me:party()
	return party ~= nil and party:leader_id() == me:id()
end

function M.mob_count(map, mob_id)
	if map == nil then
		return 0
	end
	local n = 0
	local mobs
	if mob_id ~= nil then
		mobs = map:mobs(mob_id)
	else
		mobs = map:mobs()
	end
	for _, mob in pairs(mobs) do
		if mob ~= nil then
			n = n + 1
		end
	end
	return n
end

function M.shuffle_reactors(map, min_id, max_id, exclude_id)
	if map == nil then
		return
	end
	local reactors = {}
	local positions = {}
	for _, reactor in pairs(map:reactors()) do
		local id = reactor:id()
		if exclude_id == nil or id ~= exclude_id then
			if min_id == nil or max_id == nil or (id >= min_id and id <= max_id) then
				local x, y = reactor:position()
				reactors[#reactors + 1] = reactor
				positions[#positions + 1] = { x, y }
			end
		end
	end
	for i = #positions, 2, -1 do
		local j = math.random(i)
		positions[i], positions[j] = positions[j], positions[i]
	end
	for i, reactor in ipairs(reactors) do
		reactor:position(positions[i][1], positions[i][2])
	end
end

function M.warp_portal(me, map_id, portal_name)
	local sm = me:state_machine()
	local dest = nil
	if sm ~= nil then
		local group = sm:group()
		if group ~= nil then
			dest = group:map(map_id)
		end
	end
	if dest == nil then
		me:map(map_id, 0)
		return
	end
	local spawn = 0
	if portal_name ~= nil and portal_name ~= "" then
		local portal = dest:portal(portal_name)
		if portal ~= nil then
			spawn = portal:id()
		end
	end
	me:play_portal_sound()
	me:map(dest, spawn)
end

function M.party_all_here(me)
	local party = me:party()
	local map = me:map()
	if party == nil or map == nil then
		return false
	end
	local wz = map:wz()
	if wz == nil then
		return false
	end
	local map_id = wz.id
	for _, mem in ipairs(party:members()) do
		if mem == nil or mem:map_id() ~= map_id then
			return false
		end
		if map:characters()[mem:id()] == nil then
			return false
		end
	end
	return true
end

return M
