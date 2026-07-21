local M = {}

local seeded = false

local function ensure_seed()
	if seeded then
		return
	end
	math.randomseed(now())
	seeded = true
end

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
	ensure_seed()
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

function M.party_warp(sm, map_id)
	if sm == nil or map_id == nil then
		return
	end
	for _, p in ipairs(sm:players()) do
		if p ~= nil then
			p:map(map_id)
		end
	end
end

function M.is_gm(me)
	if me == nil then
		return false
	end
	return me:role() == ROLE.Admin
end

function M.rand(a, b)
	ensure_seed()
	return math.random(a, b)
end

return M
