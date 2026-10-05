local pq = require("script/integration/lib/party_quest")

local ENTRY = 300030100
local START = 930000000
local PATH = 930000010
local STAGE1 = 930000100
local STAGE2 = 930000200
local MAZE = 930000300
local STAGE4 = 930000400
local STAGE5 = 930000500
local BOSS = 930000600
local EXIT = 930000800
local ELLIN = 2133000
local ELLA = 2133001
local SPIRIT = 2133004
local DIALOG_LIST = 4
local FOREST_KEY = 4001161
local FOREST_KEY_MOB = 9300173
local SPINE_SEED = 4001162
local SEED_BOX = 3002000
local SPINE = 3009000
local SPINE_FEED = 4
local MAZE_PORTAL = "001E"
local MAZE_GOAL = "16st"
local BEAD = 2270004
local MARBLE = 4001169
local MARBLE_COUNT = 20
local SPRITE = 9300175
local SPRITE_DAMAGE = 4000
local STONE_BOX = 3002001
local PURPLE_STONE = 4001163
local ALTAR = 3001000
local POISON_GOLEM = 9300180
local ALTAIR_FRAGMENT = 4001198

local function clear_mobs(ctx, bot)
	for _ = 1, 200 do
		local mob = bot:mobs()[1]
		if mob == nil then
			return true
		end
		bot:kill(mob.oid)
	end
	return ctx:fail(bot:name() .. " 몹을 모두 처치하지 못함 (맵 " .. bot:map() .. ")")
end

local function sweep(ctx, bot)
	for _, spot in ipairs(bot:mob_spots()) do
		if pq.move(bot, spot.x, spot.y) == false then
			return ctx:fail(bot:name() .. " 몹 위치로 이동 실패")
		end
		if clear_mobs(ctx, bot) == false then
			return false
		end
	end
	return true
end

local function pass_all(ctx, portal, map_id)
	for i = 0, ctx:bot_count() - 1 do
		if pq.portal(ctx, ctx:bot(i), portal, map_id) == false then
			return false
		end
	end
	return true
end

local function click_warp(ctx, bot, oid, map_id)
	local warp = bot:request(resp.warp, req.npc_click { oid = oid }, function(p)
		return p.character.map == map_id
	end, 15000)
	if warp == false then
		return ctx:fail(bot:name() .. " 맵 이동 실패: " .. map_id)
	end
	return pq.wait_all(ctx, map_id, 15000)
end

local function start(ctx)
	local leader = ctx:bot(0)
	local oid = pq.npc(ctx, leader, ELLIN)
	if oid == false then
		return false
	end
	local dlg = leader:npc_click(oid)
	if dlg == false or dlg.selections == nil then
		return ctx:fail("엘린 선택지가 오지 않음")
	end
	local warp = leader:request(resp.warp, req.dialog { dialog_type = DIALOG_LIST, next = true, selected = 2 }, function(p)
		return p.character.map == START
	end, 15000)
	if warp == false then
		return ctx:fail("독안개의 숲 파티 퀘스트 시작 실패")
	end
	if pq.wait_all(ctx, START, 15000) == false then
		return false
	end
	if pass_all(ctx, "east00", PATH) == false then
		return false
	end
	return pass_all(ctx, "east00", STAGE1)
end

local function clear_stage1(ctx)
	local leader = ctx:bot(0)
	for _ = 1, 5 do
		if sweep(ctx, leader) == false then
			return false
		end
		leader:warp("east00")
		if leader:map() == STAGE2 then
			for i = 1, ctx:bot_count() - 1 do
				if pq.portal(ctx, ctx:bot(i), "east00", STAGE2) == false then
					return false
				end
			end
			return true
		end
	end
	return ctx:fail("몹을 모두 처치해도 1단계 포탈이 열리지 않음")
end

local function grow_spine(ctx)
	local leader = ctx:bot(0)
	for _ = 1, SPINE_FEED do
		if pq.collect(ctx, leader, FOREST_KEY_MOB, FOREST_KEY, 1) == false then
			return false
		end
		if pq.feed(ctx, leader, SEED_BOX, FOREST_KEY, 1) == false then
			return false
		end
		if pq.pick_up(ctx, leader, SPINE_SEED) == false then
			return false
		end
		if pq.feed(ctx, leader, SPINE, SPINE_SEED, 1) == false then
			return false
		end
	end
	return pass_all(ctx, "east00", MAZE)
end

local function solve_maze(ctx)
	local leader = ctx:bot(0)
	for _ = 1, 100 do
		if leader:spawn_portal() == MAZE_GOAL then
			break
		end
		if leader:warp(MAZE_PORTAL) == false then
			return ctx:fail("미로 포탈 응답 없음")
		end
	end
	if leader:spawn_portal() ~= MAZE_GOAL then
		return ctx:fail("미로 끝에 도착하지 못함")
	end
	local oid = pq.npc(ctx, leader, ELLA)
	if oid == false then
		return false
	end
	return click_warp(ctx, leader, oid, STAGE4)
end

local function receive_beads(ctx, leader)
	local oid = pq.npc(ctx, leader, ELLA)
	if oid == false then
		return false
	end
	if leader:npc_click(oid) == false then
		return ctx:fail("엘라가 정화의 구슬을 주지 않음")
	end
	leader:dialog(false)
	if (leader:items()[BEAD] or 0) == 0 then
		return ctx:fail("정화의 구슬을 받지 못함")
	end
	return true
end

local function purify_sprites(ctx)
	local leader = ctx:bot(0)
	local spots = leader:mob_spots()
	local visited = 0
	while (leader:items()[MARBLE] or 0) < MARBLE_COUNT do
		if (leader:items()[BEAD] or 0) == 0 and receive_beads(ctx, leader) == false then
			return false
		end
		local mob = leader:mobs(SPRITE)[1]
		if mob ~= nil then
			leader:attack(mob.oid, SPRITE_DAMAGE)
			if leader:catch(mob.oid, BEAD) == false then
				return ctx:fail("정화의 구슬로 몹을 잡지 못함: " .. mob.oid)
			end
			visited = 0
		elseif visited < #spots then
			visited = visited + 1
			local spot = spots[visited]
			if pq.move(leader, spot.x, spot.y) == false then
				return ctx:fail("몹 위치로 이동 실패")
			end
		else
			local spawn = leader:request(resp.spawn_mob, nil, function(p)
				return p.mob.mob_id == SPRITE
			end, 30000)
			if spawn == false then
				return ctx:fail("몹 재생성 대기 시간 초과")
			end
		end
	end
	local oid = pq.npc(ctx, leader, ELLA)
	if oid == false then
		return false
	end
	return click_warp(ctx, leader, oid, STAGE5)
end

local function find_stone(ctx)
	local leader = ctx:bot(0)
	if pq.open_box(ctx, leader, STONE_BOX, PURPLE_STONE) == false then
		return false
	end
	local oid = pq.npc(ctx, leader, SPIRIT)
	if oid == false then
		return false
	end
	return click_warp(ctx, leader, oid, BOSS)
end

local function defeat_golem(ctx)
	local leader = ctx:bot(0)
	if pq.move(leader, 392, 148) == false then
		return ctx:fail("제단 앞으로 이동 실패")
	end
	local altar = leader:reactors(ALTAR)[1]
	if altar == nil then
		return ctx:fail("제단이 보이지 않음")
	end
	if leader:drop(PURPLE_STONE, 1) == nil then
		return ctx:fail("마력석 버리기 실패")
	end
	local golem = leader:request(resp.spawn_mob, nil, function(p)
		return p.mob.mob_id == POISON_GOLEM
	end, 15000)
	if golem == false then
		return ctx:fail("제단에 마력석을 놓아도 골렘이 나타나지 않음")
	end
	if leader:kill(golem.mob.oid) == false then
		return ctx:fail("골렘 처치 실패")
	end
	repeat
		if clear_mobs(ctx, leader) == false then
			return false
		end
	until leader:request(resp.spawn_mob, nil, nil, 3000) == false
	for i = 0, ctx:bot_count() - 1 do
		local bot = ctx:bot(i)
		local before = bot:items()[ALTAIR_FRAGMENT] or 0
		if pq.portal(ctx, bot, "out00", EXIT) == false then
			return false
		end
		if (bot:items()[ALTAIR_FRAGMENT] or 0) ~= before + 1 then
			return ctx:fail(bot:name() .. " 알테어 조각을 받지 못함")
		end
		if pq.portal(ctx, bot, "east00", ENTRY) == false then
			return false
		end
	end
	return true
end

test_suite {
	name = "Party Quest: 엘린 PQ",
	bot_count = 3,

	on_initialize = function(ctx)
		for i = 0, ctx:bot_count() - 1 do
			local bot = ctx:bot(i)
			if pq.command(bot, "/레벨바꾸기 50", "레벨 설정") == false then
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
		start,
		clear_stage1,
		grow_spine,
		solve_maze,
		purify_sprites,
		find_stone,
		defeat_golem,
	},
}
