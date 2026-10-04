local pq = require("script/integration/lib/party_quest")

local ENTRY = 240050000
local HUB = 240050100
local CHOICE = 240050200
local DARK_CAVE = 240050310
local CLEAR = 240050400
local SIGN = 2083001
local GROOT = 2083003
local GATE_KEY = 4001087
local GATE_MOB = 9300065
local MAZES = {
	{ map = 240050101, mob = 9300067, key = 4001088, next = 240050102 },
	{ map = 240050102, mob = 9300069, key = 4001089, next = 240050103 },
	{ map = 240050103, mob = 9300071, key = 4001090, next = 240050104 },
	{ map = 240050104, mob = 9300073, key = 4001091, next = 240050105 },
}
local LAST_MOB = 9300075
local LAST_KEY = 4001092
local KEY_HOLE = 2408002
local PATH_STONE = 2408001
local DARK = 3
local DARK_BOX = 2402006
local DARK_KEY = 4001093
local DIALOG_LIST = 4

local function click_warp(ctx, bot, npc, map_id)
	local oid = pq.npc(ctx, bot, npc)
	if oid == false then
		return false
	end
	local warp = bot:request(resp.warp, req.npc_click { oid = oid }, function(p)
		return p.character.map == map_id
	end, 15000)
	if warp == false then
		return ctx:fail(bot:name() .. " 맵 이동 실패: " .. map_id)
	end
	return true
end

test_suite {
	name = "Party Quest: 혼테일 PQ",
	bot_count = 6,

	on_initialize = function(ctx)
		for i = 0, ctx:bot_count() - 1 do
			local bot = ctx:bot(i)
			if pq.command(bot, "/레벨바꾸기 80", "레벨 설정") == false then
				return ctx:fail(bot:name() .. " 레벨 설정 실패")
			end
			if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
				return ctx:fail(bot:name() .. " 플레이어 모드 설정 실패")
			end
			if bot:map_move(ENTRY) == false then
				return false
			end
		end
		return true
	end,

	on_finished = function(ctx)
		ctx:bot(0):request(resp.party_update_disband, req.party_operation { operation = PARTY.Leave }, nil, 5000)
	end,

	scenarios = {
		pq.form_party,
		function(ctx)
			local leader = ctx:bot(0)
			local oid = pq.npc(ctx, leader, SIGN)
			if oid == false then
				return false
			end
			local dlg = leader:npc_click(oid)
			if dlg == false or dlg.selections == nil then
				return ctx:fail("혼테일의 이정표 선택지가 오지 않음")
			end
			local warp = leader:request(resp.warp, req.dialog { dialog_type = DIALOG_LIST, next = true, selected = 0 }, function(p)
				return p.character.map == HUB
			end, 15000)
			if warp == false then
				return ctx:fail("생명의 동굴 입장 실패")
			end
			return pq.wait_all(ctx, HUB, 15000)
		end,
		function(ctx)
			local leader = ctx:bot(0)
			local runner = ctx:bot(1)
			if pq.collect(ctx, runner, GATE_MOB, GATE_KEY, 1) == false then
				return false
			end
			if pq.portal(ctx, runner, "in00", MAZES[1].map) == false then
				return false
			end
			for _, maze in ipairs(MAZES) do
				if pq.collect(ctx, runner, maze.mob, maze.key, 1) == false then
					return false
				end
				if pq.feed(ctx, runner, KEY_HOLE, maze.key, 1) == false then
					return false
				end
				if pq.loot_spawn(ctx, leader, maze.key) == false then
					return false
				end
				if pq.talk(ctx, leader, GROOT, "미로방의 문이 열렸습니다") == false then
					return false
				end
				if pq.portal(ctx, runner, "in00", maze.next) == false then
					return false
				end
			end
			if pq.collect(ctx, runner, LAST_MOB, LAST_KEY, 1) == false then
				return false
			end
			if runner:warp("in00") == false or runner:map() ~= HUB then
				return ctx:fail("다섯 번째 미로방에서 돌아오지 못함 (현재 " .. runner:map() .. ")")
			end
			return pq.wait_all(ctx, HUB, 15000)
		end,
		function(ctx)
			if click_warp(ctx, ctx:bot(0), SIGN, CHOICE) == false then
				return false
			end
			return pq.wait_all(ctx, CHOICE, 15000)
		end,
		function(ctx)
			local leader = ctx:bot(0)
			local stone = pq.seek_reactor(ctx, leader, pq.reactor_by_id(PATH_STONE))
			if stone == false then
				return false
			end
			for _ = 1, 4 do
				local current = pq.find_reactor(leader, stone.oid)
				if current ~= nil and current.state == DARK then
					break
				end
				if leader:hit_reactor(stone.oid) == false then
					return ctx:fail("선택의 돌 타격 응답 없음")
				end
			end
			local current = pq.find_reactor(leader, stone.oid)
			if current == nil or current.state ~= DARK then
				return ctx:fail("어둠의 동굴을 선택하지 못함")
			end
			if leader:warp("east00") == false or leader:map() ~= DARK_CAVE then
				return ctx:fail("어둠의 동굴로 이동하지 못함 (현재 " .. leader:map() .. ")")
			end
			return pq.wait_all(ctx, DARK_CAVE, 15000)
		end,
		function(ctx)
			local leader = ctx:bot(0)
			local opened = pq.visit_reactors(ctx, leader, pq.reactor_by_id(DARK_BOX), function(box)
				if pq.break_reactor(ctx, leader, box) == false then
					return false
				end
				return pq.loot_spawn(ctx, leader, DARK_KEY)
			end)
			if opened == false then
				return false
			end
			if (leader:items()[DARK_KEY] or 0) < 6 then
				return ctx:fail("어둠의 열쇠 부족: " .. (leader:items()[DARK_KEY] or 0))
			end
			if click_warp(ctx, leader, SIGN, CLEAR) == false then
				return false
			end
			return pq.wait_all(ctx, CLEAR, 15000)
		end,
	},
}
