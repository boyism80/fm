local check = require("script/integration/lib/script_check")

local BOTS = 3
local START = 1
local SCRIPTED_START = 4
local SCRIPTED_END = 5
local MAX_STEPS = 40
local entries, missing = wz.quest_scripts()

local REPLIES = { resp.update_quest }
for _, name in ipairs(check.REPLIES) do
	REPLIES[#REPLIES + 1] = name
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

local function prepare_start(bot, quest)
	local p, name = bot:request({ resp.warp, resp.notice }, req.normal_chat { message = "/퀘스트시작준비 " .. quest }, function(p, name)
		if name == resp.warp then
			return true
		end
		return p.message:find("NPC", 1, true) ~= nil
			or p.message:find("실패", 1, true) ~= nil
			or p.message:find("없음", 1, true) ~= nil
			or p.message:find("이미", 1, true) ~= nil
	end)
	if p == false then
		return false, "응답 없음"
	end
	if name == resp.warp or p.message:find("NPC", 1, true) ~= nil then
		return true
	end
	return false, p.message
end

local function started(bot, quest)
	local m = command(bot, "/퀘스트상태 " .. quest, function(m)
		return m:find("퀘스트 " .. quest .. "[ :]") ~= nil
	end)
	return m ~= nil and m:sub(-#"진행") == "진행"
end

local function visit(bot, spot)
	if spot == nil then
		return true
	end
	if bot:instance_move(spot.map) == false then
		return false
	end
	if spot.x ~= nil then
		bot:move(spot.x, spot.y, spot.foothold)
	end
	return true
end

local function answer(p, name)
	if name == resp.dialog_yes_no then
		return req.dialog { dialog_type = 1, next = true }
	elseif name == resp.dialog_accept then
		return req.dialog { dialog_type = p.enable_escape and 11 or 12, next = true }
	elseif name == resp.dialog_list then
		return req.dialog { dialog_type = 4, next = true, selected = 0 }
	elseif name == resp.dialog_input then
		return req.dialog { dialog_type = 2, next = true, text = "1" }
	elseif name == resp.dialog_style then
		return req.dialog { dialog_type = 7, next = true, selected = 0 }
	end
	return req.dialog { dialog_type = 0, next = true }
end

local function talk(ctx, bot, where, mode, quest, spot)
	local npc = 0
	if spot ~= nil then
		npc = spot.npc
	end
	local p, name = bot:request(REPLIES, req.quest_action { mode = mode, quest_id = quest, npcid = npc }, check.replied, 5000)
	if p == false then
		ctx:fail(where .. ": 응답 없음")
		return nil
	end
	if name == resp.update_stats then
		log("info", where .. ": 대화 없이 끝남")
		return nil
	end

	local status = nil
	for _ = 1, MAX_STEPS do
		if name == resp.warp or name == resp.open_npc_shop then
			return status
		end

		local next = nil
		local wait = 5000
		if name == resp.update_quest then
			if p.quest_status.quest_id == quest then
				status = p.quest_status.status
			end
		elseif name == resp.update_stats then
			wait = 500
		elseif check.is_dialog(name) then
			next = answer(p, name)
		end

		local more, more_name = bot:request(REPLIES, next, check.replied, wait)
		if more == false then
			if name == resp.update_stats then
				return status
			end
			ctx:fail(where .. ": 대화 도중 응답 없음")
			return nil
		end
		p, name = more, more_name
	end
	bot:dialog(false)
	log("info", where .. ": 대화가 " .. MAX_STEPS .. "단계 안에 끝나지 않음")
	return status
end

local function run(ctx, bot, entry)
	local quest = entry.quest
	local where = "quest " .. quest

	command(bot, "/인벤토리초기화", function(m)
		return m:find("슬롯 비움", 1, true) ~= nil
	end)
	command(bot, "/퀘스트초기화 " .. quest, function(m)
		return m:find("퀘스트 " .. quest .. " ", 1, true) ~= nil
	end)
	local ok, reason = prepare_start(bot, quest)
	if ok == false then
		log("info", where .. ": 시작 준비 실패 - " .. reason)
		return
	end
	if visit(bot, entry.start) == false then
		return
	end

	local status = nil
	if entry.start_script then
		status = talk(ctx, bot, where .. " on_start", SCRIPTED_START, quest, entry.start)
		if entry.end_script and status ~= 1 then
			log("info", where .. ": on_start 뒤 진행 중이 아니어서 on_end 생략")
			return
		end
	elseif entry.start ~= nil then
		if started(bot, quest) == false then
			local p = bot:request(resp.update_quest, req.quest_action { mode = START, quest_id = quest, npcid = entry.start.npc }, function(p)
				return p.quest_status.quest_id == quest and p.quest_status.status == 1
			end, 5000)
			if p == false then
				log("info", where .. ": 일반 시작 실패로 on_end 생략")
				return
			end
		end
	else
		log("info", where .. ": 시작 NPC가 없어 on_end 생략")
		return
	end
	if entry.end_script == false then
		return
	end

	local prepared = command(bot, "/퀘스트완료준비 " .. quest, function(m)
		return m:find("퀘스트 " .. quest .. " ", 1, true) ~= nil
	end)
	if prepared == nil or prepared:find("준비 완료", 1, true) == nil then
		log("info", where .. ": 완료 준비 실패 - " .. tostring(prepared))
		return
	end
	if visit(bot, entry.finish or entry.start) == false then
		return
	end
	if talk(ctx, bot, where .. " on_end", SCRIPTED_END, quest, entry.finish or entry.start) ~= 2 then
		log("info", where .. ": on_end 뒤 완료되지 않음")
	end
end

test_suite {
	name = "Quest scripts",
	bot_count = BOTS,

	on_initialize = function(ctx)
		log("info", string.format("%d quest scripts run by WZ, %d not run", #entries, #missing))
		return true
	end,

	scenarios = {
		{ parallel = check.queues(BOTS, entries, run) },
	},
}
