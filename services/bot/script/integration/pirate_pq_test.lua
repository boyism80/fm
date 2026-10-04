local pq = require("script/integration/lib/party_quest")

local ENTRY = 251010404
local STAGE1 = 925100000
local BOW = 925100100
local DECK1 = 925100200
local DECK2 = 925100300
local DOORS = 925100400
local BOSS = 925100500
local REWARD = 925100600
local EXIT = 925100700
local GUIDE = 2094000
local GUON = 2094002
local WYANG = 2094001
local SEALS = {
	{ mob = 9300114, item = 4001120 },
	{ mob = 9300115, item = 4001121 },
	{ mob = 9300116, item = 4001122 },
}
local KEY = 4001117
local DOOR = 2519000
local PIRATE_KING = 9300119
local DIALOG_LIST = 4

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

local function pass(ctx, bot, x, y, map_id)
	for _ = 1, 5 do
		if clear_mobs(ctx, bot) == false then
			return false
		end
		if pq.move(bot, x, y) == false then
			return ctx:fail(bot:name() .. " 포탈 앞으로 이동 실패")
		end
		bot:warp("next00")
		if bot:map() == map_id then
			return true
		end
	end
	return ctx:fail(bot:name() .. " 포탈이 열리지 않음: " .. bot:map() .. " -> " .. map_id)
end

local function pass_all(ctx, x, y, map_id)
	for i = 0, ctx:bot_count() - 1 do
		if pass(ctx, ctx:bot(i), x, y, map_id) == false then
			return false
		end
	end
	return true
end

local function select_warp(ctx, bot, npc, selected, map_id)
	local oid = pq.npc(ctx, bot, npc)
	if oid == false then
		return false
	end
	local dlg = bot:npc_click(oid)
	if dlg == false or dlg.selections == nil then
		return ctx:fail(bot:name() .. " 선택지가 오지 않음: " .. npc)
	end
	local warp = bot:request(resp.warp, req.dialog { dialog_type = DIALOG_LIST, next = true, selected = selected }, function(p)
		return p.character.map == map_id
	end, 15000)
	if warp == false then
		return ctx:fail(bot:name() .. " 맵 이동 실패: " .. map_id)
	end
	return true
end

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
	name = "Party Quest: 해적 PQ",
	bot_count = 3,

	on_initialize = function(ctx)
		for i = 0, ctx:bot_count() - 1 do
			local bot = ctx:bot(i)
			if pq.command(bot, "/레벨바꾸기 55", "레벨 설정") == false then
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
			if select_warp(ctx, ctx:bot(0), GUIDE, 2, STAGE1) == false then
				return false
			end
			return pq.wait_all(ctx, STAGE1, 15000)
		end,
		function(ctx)
			return pass_all(ctx, 1523, 82, BOW)
		end,
		function(ctx)
			local leader = ctx:bot(0)
			for _, seal in ipairs(SEALS) do
				if pq.talk(ctx, leader, GUON, "나타나는 해적을 잡고") == false then
					return false
				end
				if pq.collect(ctx, leader, seal.mob, seal.item, 20) == false then
					return false
				end
				if pq.talk(ctx, leader, GUON, "모두 모아오셨군요") == false then
					return false
				end
			end
			return pass_all(ctx, 1324, 240, DECK1)
		end,
		function(ctx)
			return pass_all(ctx, 2133, 240, DECK2)
		end,
		function(ctx)
			return pass_all(ctx, 2133, 240, DOORS)
		end,
		function(ctx)
			local leader = ctx:bot(0)
			if pq.collect(ctx, leader, nil, KEY, 4) == false then
				return false
			end
			for i = 0, 3 do
				if pq.feed(ctx, leader, DOOR + i, KEY, 1) == false then
					return false
				end
			end
			for _ = 1, 5 do
				if clear_mobs(ctx, leader) == false then
					return false
				end
				if pq.move(leader, 1327, 238) == false then
					return ctx:fail("보스 포탈 앞으로 이동 실패")
				end
				leader:warp("next00")
				if leader:map() == BOSS then
					return pq.wait_all(ctx, BOSS, 15000)
				end
			end
			return ctx:fail("해적왕의 방으로 가는 포탈이 열리지 않음")
		end,
		function(ctx)
			local leader = ctx:bot(0)
			if leader:mobs(PIRATE_KING)[1] == nil and leader:request(resp.spawn_mob, nil, function(p)
				return p.mob.mob_id == PIRATE_KING
			end, 10000) == false then
				return ctx:fail("해적왕이 나타나지 않음")
			end
			if clear_mobs(ctx, leader) == false then
				return false
			end
			if click_warp(ctx, leader, GUON, REWARD) == false then
				return false
			end
			return pq.wait_all(ctx, REWARD, 15000)
		end,
		function(ctx)
			for i = 0, ctx:bot_count() - 1 do
				local bot = ctx:bot(i)
				if select_warp(ctx, bot, WYANG, 1, EXIT) == false then
					return false
				end
				if click_warp(ctx, bot, GUON, ENTRY) == false then
					return false
				end
			end
			return true
		end,
	},
}
