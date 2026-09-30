local pq = require("script/lib/party_quest")

local M = {}

local SCRIPT = "script/lib/boss_entry.lua"
local EXIT_MAP = 123256780
local MIN_LEVEL = 120
local MIN_MEMBERS = 1
local MAX_MEMBERS = 6

function M.character_count(map)
	local n = 0
	for _, _ in pairs(map:characters()) do
		n = n + 1
	end
	return n
end

function M.reset_map(map)
	map:reset()
end

function M.spawn_boss(map, x, y, ...)
	for _, mob_id in ipairs({ ... }) do
		map:spawn_mob(mob_id, x, y)
	end
end

function M.on_click(me, npc, boss)
	local party = me:party()
	if party == nil then
		me:dialog(npc, "파티를 구성하여 주세요.")
		return
	end
	if me:map():wz():id() ~= EXIT_MAP then
		me:dialog(npc, "#m" .. EXIT_MAP .. "#맵에서만 저에게 말을 걸어 입장하실 수 있습니다.")
		return
	end
	if party:leader_id() ~= me:id() then
		me:dialog(npc, "파티장만이 제게 말을 거실 수 있습니다.")
		return
	end
	local text = "#e< #b" .. boss.name .. " 입장 조건#k >#n#b\r\n"
		.. "- #e#r난이도 #k: " .. boss.difficulty .. "\r\n"
		.. "- 120레벨 이상의 파티원\r\n"
		.. "- 파티원 모두 같은맵에 존재\r\n"
		.. "- 최소 1인부터 최대 6인의 파티\r\n\r\n"
		.. "#k"
		.. "위의 조건을 모두 충족하셨다면, 수락하기를 눌러주세요.\r\n"
	if me:dialog_accept(npc, text) == false then
		return
	end

	for _, map_id in ipairs(boss.maps) do
		local ok, count = run_on_map(map_id, SCRIPT, "character_count")
		if ok and count ~= nil and count > 0 then
			me:dialog(npc, "이미 누군가가 이용 중에 있습니다.")
			return
		end
	end
	party = me:party()
	if party == nil then
		me:dialog(npc, "파티를 구성하고 있지 않으십니다.")
		return
	end
	if party:leader_id() ~= me:id() then
		me:dialog(npc, "파티장만이 저에게 말을 걸 수 있습니다.")
		return
	end

	local d = datetime()
	local today = string.format("%04d%02d%02d", d.year, d.month, d.day)
	local characters = me:map():characters()
	local size = 0
	local leveled = 0
	local in_map = 0
	local entries = 0
	for _, mem in pairs(party:members()) do
		size = size + 1
		if mem:level() >= MIN_LEVEL then
			leveled = leveled + 1
		end
		local ch = characters[mem:id()]
		if ch ~= nil then
			in_map = in_map + 1
			local count_q = ch:quest(boss.quest)
			local day_q = ch:quest(boss.day_quest)
			if count_q:started() == false then
				count_q:start(tostring(boss.limit))
			end
			if day_q:started() == false then
				day_q:start(today)
			end
			if day_q:record() ~= today then
				count_q:record(tostring(boss.limit))
				day_q:record(today)
			end
			if (tonumber(count_q:record()) or 0) > 0 then
				entries = entries + 1
			end
		end
	end
	if pq.is_gm(me) == false then
		if size > MAX_MEMBERS or size < MIN_MEMBERS or in_map < size or leveled < size then
			me:dialog(npc, "파티원 중 레벨이 부족한 사람이 있거나, 파티원 수가 안맞는거 같습니다. 입장 조건을 다시 확인하여주세요.")
			return
		end
		if entries < size then
			me:dialog(npc, "파티원 중 오늘 이미 " .. boss.limit .. "번 입장하신 분이 있습니다.")
			return
		end
	end

	local pid = party:id()
	for _, ch in pairs(characters) do
		local p = ch:party()
		if p ~= nil and p:id() == pid then
			local count_q = ch:quest(boss.quest)
			count_q:record(tostring((tonumber(count_q:record()) or 0) - 1))
		end
	end
	for _, map_id in ipairs(boss.maps) do
		run_on_map(map_id, SCRIPT, "reset_map")
	end
	run_on_map(boss.maps[1], SCRIPT, "spawn_boss", boss.x, boss.y, boss.mobs[1], boss.mobs[2])
	for _, ch in pairs(characters) do
		local p = ch:party()
		if p ~= nil and p:id() == pid then
			ch:map(boss.maps[1], 0)
			ch:warp_later(EXIT_MAP, boss.time * 60)
		end
	end
end

return M
