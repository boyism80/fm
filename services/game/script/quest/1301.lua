local quest_id = 1301

local RANK_UP = { A = "S", B = "A", C = "B", D = "C", F = "D" }

local function rank_check_met(mode, value, expected)
	if mode == "less" then
		return value < expected
	end
	if mode == "more" then
		return value > expected
	end
	if mode == "equal" then
		return value == expected
	end
	return false
end

local function recalc_rank(quest)
	if quest == nil then
		return
	end
	local old_rank = quest:record_ex("rank")
	if old_rank == nil or old_rank == "" or old_rank == "S" then
		return
	end
	local new_rank = RANK_UP[old_rank]
	if new_rank == nil then
		return
	end
	local wz = quest:wz()
	if wz == nil or wz.party_ranks == nil then
		return
	end
	local checks = wz.party_ranks[new_rank]
	if checks == nil then
		return
	end
	for _, check in ipairs(checks) do
		if check == nil or check.property == nil then
			return
		end
		local raw = quest:record_ex(check.property)
		if raw == nil then
			return
		end
		local parsed = tonumber(raw)
		if parsed == nil then
			return
		end
		if not rank_check_met(check.mode, parsed, check.value) then
			return
		end
	end
	quest:record_ex("rank", new_rank)
end

local function party_members_on_map(me)
	local map = me:map()
	local party = me:party()
	if map == nil or party == nil then
		return { me }
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
		return { me }
	end
	return out
end

function on_quest_start_1301(me)
	local quest = me:quest(quest_id)
	if quest == nil then
		return
	end
	if not quest:started() then
		if not quest:start(0, true) then
			return
		end
		quest = me:quest(quest_id)
		if quest == nil then
			return
		end
	end
	if quest:record_ex("rank") ~= nil then
		return
	end
	quest:record_ex("min", "0")
	quest:record_ex("sec", "0")
	quest:record_ex("date", "0000-00-00")
	quest:record_ex("have", "0")
	quest:record_ex("rank", "F")
	quest:record_ex("try", "0")
	quest:record_ex("cmp", "0")
	quest:record_ex("CR", "0")
	quest:record_ex("VR", "0")
	quest:record_ex("gvup", "0")
	quest:record_ex("vic", "0")
	quest:record_ex("lose", "0")
	quest:record_ex("draw", "0")
end

function on_quest_try_1301(me)
	on_quest_start_1301(me)
	local quest = me:quest(quest_id)
	if quest == nil then
		return
	end
	quest:start_time(now())
	local try_val = tonumber(quest:record_ex("try")) or 0
	quest:record_ex("try", tostring(try_val + 1))
end

local function finish_one_1301(me)
	on_quest_start_1301(me)
	local quest = me:quest(quest_id)
	if quest == nil then
		return
	end
	local started = quest:start_time()
	if started == nil or started <= 0 then
		return
	end
	local elapsed = now() - started
	if elapsed < 0 then
		elapsed = 0
	end
	local mins = math.floor(elapsed / 60)
	local secs = elapsed % 60
	local best_min = tonumber(quest:record_ex("min"))
	if best_min ~= nil and (best_min <= 0 or mins < best_min) then
		quest:record_ex("min", tostring(mins))
		quest:record_ex("sec", tostring(secs))
		local dt = datetime()
		quest:record_ex("date", string.format("%04d-%02d-%02d", dt.year, dt.month, dt.day))
	end
	local cmp = tonumber(quest:record_ex("cmp"))
	if cmp == nil then
		quest:start_time(0)
		return
	end
	local new_cmp = cmp + 1
	quest:record_ex("cmp", tostring(new_cmp))
	local try_count = tonumber(quest:record_ex("try"))
	if try_count ~= nil and try_count > 0 then
		local cr = math.floor((new_cmp * 100) / try_count)
		if cr < 0 then
			cr = 0
		end
		quest:record_ex("CR", tostring(cr))
	end
	recalc_rank(quest)
	quest:start_time(0)
end

function on_quest_end_1301(me)
	local members = party_members_on_map(me)
	if #members <= 1 then
		finish_one_1301(me)
		return
	end
	for _, member in ipairs(members) do
		finish_one_1301(member)
	end
end
