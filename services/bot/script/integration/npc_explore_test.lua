local check = require("script/integration/lib/script_check")

local BOTS = 4
local MAX_PATHS = 16
local MAX_STEPS = 40
local MAX_GRANTS = 24
local STACK = 50
local entries, missing = wz.npc_scripts()

local PROFILES = {
	{ name = "초보자", level = 10, job = 0, meso = 0, grant = false },
	{ name = "숙련자", level = 70, job = 110, meso = 100000000, grant = true },
	{ name = "퀘스트진행", level = 70, job = 110, meso = 100000000, grant = true, quest = 1 },
	{ name = "퀘스트완료", level = 70, job = 110, meso = 100000000, grant = true, quest = 2 },
}

local REPLIES = { resp.notice, resp.update_quest, resp.guild_message }
for _, name in ipairs(check.REPLIES) do
	REPLIES[#REPLIES + 1] = name
end

local summary = { npcs = 0, paths = 0, suspects = 0 }

local function replied(p, name)
	if name == resp.notice then
		return p.message:find("스크립트 오류", 1, true) == 1
	end
	return check.replied(p, name)
end

local function command(bot, text, done)
	local p = bot:request(resp.notice, req.normal_chat { message = text }, function(p)
		return done(p.message)
	end)
	if p == false then
		return nil
	end
	return p.message
end

local function snapshot(bot)
	local m = command(bot, "/봇상태", function(m)
		return m:find("봇상태 ", 1, true) == 1
	end)
	if m == nil then
		return nil
	end
	local s = { items = {} }
	for key, value in m:gmatch("(%a+)=(%-?%d+)") do
		if key ~= "items" then
			s[key] = tonumber(value)
		end
	end
	for id, count in (m:match("items=(%S*)") or ""):gmatch("(%d+):(%d+)") do
		s.items[tonumber(id)] = tonumber(count)
	end
	return s
end

local function reset(bot, profile, grants, quests)
	local text = string.format("/봇초기화 %d %d %d", profile.level, profile.job, profile.meso)
	if profile.grant and grants ~= "" then
		text = text .. " " .. grants
	else
		text = text .. " -"
	end
	if profile.quest ~= nil then
		local states = {}
		for i, id in ipairs(quests) do
			states[i] = id .. ":" .. profile.quest
		end
		text = text .. " " .. table.concat(states, ",")
	end
	return command(bot, text, function(m)
		return m == "봇초기화 완료" or m:find("사용법", 1, true) == 1
	end) == "봇초기화 완료"
end

local function options(p, name)
	if name == resp.dialog_yes_no then
		return {
			{ label = "예", pkt = req.dialog { dialog_type = 1, next = true } },
			{ label = "아니오", pkt = req.dialog { dialog_type = 1, next = false } },
		}
	elseif name == resp.dialog_accept then
		local kind = p.enable_escape and 11 or 12
		return {
			{ label = "수락", pkt = req.dialog { dialog_type = kind, next = true } },
			{ label = "거절", pkt = req.dialog { dialog_type = kind, next = false } },
		}
	elseif name == resp.dialog_list then
		local list = {}
		for i = 1, #p.selections do
			list[i] = { label = "목록" .. i, pkt = req.dialog { dialog_type = 4, next = true, selected = i - 1 } }
		end
		return list
	elseif name == resp.dialog_input then
		return { { label = "입력", pkt = req.dialog { dialog_type = 2, next = true, text = "1" } } }
	elseif name == resp.dialog_style then
		return { { label = "스타일", pkt = req.dialog { dialog_type = 7, next = true, selected = 0 } } }
	end
	return { { pkt = req.dialog { dialog_type = 0, next = true } } }
end

local function walk(bot, npc, oid, prefix, queue)
	local trace = { picks = {}, labels = {}, quests = {}, dialogs = 0 }
	local p, name = bot:request(REPLIES, req.npc_click { oid = oid }, replied, 5000)
	if p == false then
		trace.ended = "응답 없음"
		return trace
	end
	if check.is_dialog(name) and p.npc ~= npc then
		local own, own_name = bot:request(REPLIES, nil, replied, 500)
		if own ~= false then
			p, name = own, own_name
		end
	end

	for _ = 1, MAX_STEPS do
		local next = nil
		local wait = 5000
		if name == resp.warp then
			trace.ended = "맵 이동"
			return trace
		elseif name == resp.open_npc_shop then
			trace.ended = "상점"
			return trace
		elseif name == resp.guild_message then
			trace.ended = "길드 창"
			return trace
		elseif name == resp.dialog_list and #(p.selections or {}) == 0 then
			trace.ended = "빈 목록"
			trace.empty = true
			bot:send(req.dialog { dialog_type = 4, next = false })
			return trace
		elseif name == resp.notice then
			trace.error = p.message
			wait = 1000
		elseif name == resp.update_quest then
			trace.quests[#trace.quests + 1] = p.quest_status.quest_id .. ":" .. p.quest_status.status
		elseif name == resp.update_stats then
			wait = 500
		else
			trace.dialogs = trace.dialogs + 1
			local opts = options(p, name)
			local at = #trace.picks + 1
			local pick = prefix[at] or 1
			if prefix[at] == nil then
				for i = 2, #opts do
					local child = {}
					for k = 1, at - 1 do
						child[k] = trace.picks[k]
					end
					child[at] = i
					queue[#queue + 1] = child
				end
			end
			trace.picks[at] = pick
			if opts[pick] == nil then
				trace.ended = "선택지 사라짐"
				return trace
			end
			if opts[pick].label ~= nil then
				trace.labels[#trace.labels + 1] = opts[pick].label
			end
			next = opts[pick].pkt
		end

		local more, more_name = bot:request(REPLIES, next, replied, wait)
		if more == false then
			if name == resp.update_stats or name == resp.notice then
				trace.ended = "종료"
			else
				trace.ended = "응답 없음"
			end
			return trace
		end
		p, name = more, more_name
	end
	trace.ended = "단계 초과"
	return trace
end

local function diff(before, after, seen)
	local effects = {}
	if before == nil or after == nil then
		return { "상태 조회 실패" }
	end
	if after.level ~= before.level then
		effects[#effects + 1] = string.format("레벨 %d→%d", before.level, after.level)
	elseif after.exp ~= before.exp then
		effects[#effects + 1] = string.format("경험치 %+d", after.exp - before.exp)
		seen.exp = true
	end
	if after.job ~= before.job then
		effects[#effects + 1] = string.format("직업 %d→%d", before.job, after.job)
	end
	if after.meso ~= before.meso then
		effects[#effects + 1] = string.format("메소 %+d", after.meso - before.meso)
		seen.meso = true
	end
	if after.fame ~= before.fame then
		effects[#effects + 1] = string.format("인기도 %+d", after.fame - before.fame)
	end
	if after.map ~= before.map then
		effects[#effects + 1] = string.format("맵 %d", after.map)
		seen.maps[after.map] = true
	end
	local ids = {}
	for id in pairs(before.items) do
		ids[#ids + 1] = id
	end
	for id in pairs(after.items) do
		if before.items[id] == nil then
			ids[#ids + 1] = id
		end
	end
	table.sort(ids)
	for _, id in ipairs(ids) do
		local delta = (after.items[id] or 0) - (before.items[id] or 0)
		if delta ~= 0 then
			effects[#effects + 1] = string.format("아이템 %d %+d", id, delta)
			seen.items[id] = true
		end
	end
	return effects
end

local function suspects(entry, refs, seen)
	local list = {}
	for _, id in ipairs(refs.maps) do
		if id ~= entry.map and seen.maps[id] == nil then
			list[#list + 1] = string.format("스크립트의 맵 %d로 가는 갈래 없음", id)
		end
	end
	for _, id in ipairs(refs.changes) do
		if seen.items[id] == nil then
			list[#list + 1] = string.format("스크립트의 아이템 %d 변화 없음", id)
		end
	end
	if refs.meso and seen.meso == false then
		list[#list + 1] = "스크립트에 메소 변화가 있으나 일어나지 않음"
	end
	if refs.exp and seen.exp == false then
		list[#list + 1] = "스크립트에 경험치 변화가 있으나 일어나지 않음"
	end
	for _, path in ipairs(seen.empty) do
		list[#list + 1] = "선택지가 없는 목록: " .. path
	end
	return list
end

local function explore(ctx, bot, entry)
	local where = "npc " .. entry.npc .. " (map " .. entry.map .. ")"
	local refs = wz.script_refs("npc", tostring(entry.npc))
	if refs == nil then
		ctx:fail(where .. ": 스크립트 파일을 읽지 못함")
		return
	end
	local grants = {}
	for _, id in ipairs(refs.items) do
		if #grants >= MAX_GRANTS then
			break
		end
		local count = STACK
		if id < 2000000 then
			count = 1
		end
		grants[#grants + 1] = id .. ":" .. count
	end
	grants = table.concat(grants, ",")

	local seen = { maps = {}, items = {}, meso = false, exp = false, empty = {} }
	local lines = { where }
	summary.npcs = summary.npcs + 1
	for _, profile in ipairs(PROFILES) do
		if profile.quest ~= nil and #refs.quests == 0 then
			break
		end
		local queue = { {} }
		local runs = 0
		local away = true
		while #queue > 0 and runs < MAX_PATHS do
			local prefix = table.remove(queue, 1)
			runs = runs + 1
			summary.paths = summary.paths + 1

			if reset(bot, profile, grants, refs.quests) == false then
				ctx:fail(where .. ": /봇초기화 실패")
				return
			end
			if away then
				if bot:instance_move(entry.map) == false then
					ctx:fail(where .. ": 맵 이동 실패")
					return
				end
				bot:move(entry.x, entry.y, entry.foothold)
			end
			local oid = bot:npc(entry.npc, 3000)
			if oid == nil then
				ctx:fail(where .. ": 맵에 NPC가 없음")
				return
			end

			local before = snapshot(bot)
			local trace = walk(bot, entry.npc, oid, prefix, queue)
			local after = snapshot(bot)
			local effects = diff(before, after, seen)
			away = after == nil or after.map ~= entry.map

			local path = table.concat(trace.labels, " > ")
			if path == "" then
				path = "-"
			end
			local label = string.format("%s [%s] %s", where, profile.name, path)
			if trace.error ~= nil then
				ctx:fail(label .. ": " .. trace.error)
			end
			if trace.ended == "응답 없음" then
				ctx:fail(label .. ": 대화 도중 응답 없음")
			end
			if trace.empty then
				seen.empty[#seen.empty + 1] = string.format("[%s] %s", profile.name, path)
			end
			if #trace.quests > 0 then
				effects[#effects + 1] = "퀘스트 " .. table.concat(trace.quests, ",")
			end
			lines[#lines + 1] = string.format("  [%s] %s => %s (대화 %d, %s)", profile.name, path, table.concat(effects, ", "), trace.dialogs, trace.ended)
			if trace.ended == "상점" then
				break
			end
		end
		if #queue > 0 then
			lines[#lines + 1] = string.format("  [%s] 경로 상한 %d 도달, 남은 분기 %d", profile.name, MAX_PATHS, #queue)
		end
	end

	local list = suspects(entry, refs, seen)
	for _, s in ipairs(list) do
		lines[#lines + 1] = "  의심: " .. s
	end
	if #list > 0 then
		summary.suspects = summary.suspects + 1
	end
	ctx:report(table.concat(lines, "\n"))
end

test_suite {
	name = "NPC explore",
	bot_count = BOTS,

	on_initialize = function(ctx)
		log("info", string.format("%d npc scripts placed in WZ, %d not placed", #entries, #missing))
		return true
	end,

	scenarios = {
		{ parallel = check.queues(BOTS, entries, explore) },
	},

	on_finished = function(ctx)
		local text = string.format("NPC %d개, 경로 %d개, 의심 NPC %d개", summary.npcs, summary.paths, summary.suspects)
		ctx:report(text)
		log("info", text)
	end,
}
