local pq = require("script/integration/lib/party_quest")

local DIALOGS = { resp.dialog, resp.dialog_yes_no, resp.dialog_list, resp.warp }
local DARK_MARBLE = 4031013
local PROOF = 4031012

local function talk(ctx, bot, template, picks)
	local oid = pq.npc(ctx, bot, template)
	if oid == false then
		return false
	end
	local p, name = bot:request(DIALOGS, req.npc_click { oid = oid })
	for _ = 1, 20 do
		if p == false then
			return ctx:fail(string.format("NPC %d 대화 중 응답 없음 (맵 %d)", template, bot:map()))
		end
		if name == resp.warp then
			return true
		end
		if name == resp.dialog and p.next == false then
			bot:dialog(false)
			return p.text
		end
		local reply = req.dialog { dialog_type = 0, next = true }
		if name == resp.dialog_yes_no then
			reply = req.dialog { dialog_type = 1, next = true }
		elseif name == resp.dialog_list then
			reply = req.dialog { dialog_type = 4, next = true, selected = table.remove(picks, 1) }
		end
		p, name = bot:request(DIALOGS, reply)
	end
	return ctx:fail(string.format("NPC %d 대화가 끝나지 않음", template))
end

local function advance(ctx, path)
	local bot = ctx:bot(0)
	if pq.command(bot, "/봇초기화 30 " .. path.class .. " 0 - -", "봇초기화 완료") == false then
		return ctx:fail("봇 초기화 실패")
	end
	if bot:map_move(path.master_map) == false then
		return ctx:fail("전직관에게 이동 실패: " .. path.master_map)
	end
	if talk(ctx, bot, path.master, {}) == false then
		return false
	end
	if (bot:items()[path.letter] or 0) ~= 1 then
		return ctx:fail("추천서를 받지 못함: " .. path.letter)
	end

	if bot:map_move(path.instructor_map) == false then
		return ctx:fail("전직 교관에게 이동 실패: " .. path.instructor_map)
	end
	if talk(ctx, bot, path.instructor, {}) == false then
		return false
	end
	if bot:map() < path.test_map or bot:map() > path.test_map + 2 then
		return ctx:fail("시험장으로 이동하지 않음: " .. bot:map())
	end
	if pq.collect(ctx, bot, nil, DARK_MARBLE, 30) == false then
		return false
	end
	if talk(ctx, bot, path.examiner, {}) == false then
		return false
	end
	if pq.wait_map(ctx, bot, path.exit or path.instructor_map, 5000) == false then
		return false
	end
	if (bot:items()[PROOF] or 0) ~= 1 or (bot:items()[DARK_MARBLE] or 0) ~= 0 then
		return ctx:fail("영웅의 증거와 검은 구슬 교환 실패")
	end

	if bot:map_move(path.master_map) == false then
		return ctx:fail("전직관에게 돌아가지 못함")
	end
	local said = talk(ctx, bot, path.master, { path.choose, 0 })
	if said == false then
		return false
	end
	if bot:class() ~= path.second then
		return ctx:fail(string.format("2차 전직 실패: %d (기대 %d) %s", bot:class(), path.second, tostring(said)))
	end
	if (bot:items()[PROOF] or 0) ~= 0 then
		return ctx:fail("영웅의 증거가 회수되지 않음")
	end
	return true
end

test_suite {
	name = "Job: 2차 전직",
	bot_count = 1,

	on_initialize = function(ctx)
		local bot = ctx:bot(0)
		if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
			return ctx:fail("플레이어 모드 설정 실패")
		end
		if pq.command(bot, "/무적", "무적 상태: enabled") == false then
			return ctx:fail("무적 설정 실패")
		end
		if pq.command(bot, "/즉사", "즉사 상태: enabled") == false then
			return ctx:fail("즉사 설정 실패")
		end
		return true
	end,

	scenarios = {
		function(ctx)
			return advance(ctx, {
				class = 100,
				second = 110,
				master = 1022000,
				master_map = 102000003,
				letter = 4031008,
				instructor = 1072000,
				instructor_map = 102020300,
				test_map = 108000300,
				examiner = 1072004,
				choose = 3,
			})
		end,
		function(ctx)
			return advance(ctx, {
				class = 200,
				second = 210,
				master = 1032001,
				master_map = 101000003,
				letter = 4031009,
				instructor = 1072001,
				instructor_map = 101020000,
				test_map = 108000200,
				examiner = 1072005,
				exit = 101010000,
				choose = 3,
			})
		end,
		function(ctx)
			return advance(ctx, {
				class = 300,
				second = 310,
				master = 1012100,
				master_map = 100000201,
				letter = 4031010,
				instructor = 1072002,
				instructor_map = 106010000,
				test_map = 108000100,
				examiner = 1072006,
				choose = 2,
			})
		end,
		function(ctx)
			return advance(ctx, {
				class = 400,
				second = 410,
				master = 1052001,
				master_map = 103000003,
				letter = 4031011,
				instructor = 1072003,
				instructor_map = 102040000,
				test_map = 108000400,
				examiner = 1072007,
				choose = 2,
			})
		end,
	},
}
