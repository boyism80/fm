local pq = require("script/integration/lib/party_quest")

local ENTRY = 211042300
local MINE = 280010000
local EXIT = 280090000
local BREATH_ROAD = 280020000
local BREATH_END = 280020001
local ADOBIS = 2030008
local AURA = 2032002
local ALI = 2030011
local LIRA = 2032003
local KEY = 4001016
local KEY_BOXES = { [2112004] = true, [2112011] = true }
local CHEST = 2112014
local FIRE_ORE = 4001018
local FIRE_ORE_PIECE = 4031061
local VOLCANO_BREATH = 4031062
local EYE_OF_FIRE = 4001017
local DIALOG_DEFAULT = 0
local DIALOG_LIST = 4

local CAVES = {
	{ { "in04", 280010040 }, { "east00", 280010041 } },
	{ { "in09", 280010090 }, { "west00", 280010091 } },
	{ { "in11", 280010110 } },
	{ { "in14", 280010140 } },
}
local CAVE_BACK = {
	{ { "west00", 280010040 }, { "west00", MINE } },
	{ { "east00", 280010090 }, { "east00", MINE } },
	{ { "east00", MINE } },
	{ { "west00", MINE } },
}
local DEEP = { { "in10", 280010100 }, { "east00", 280010101 }, { "east00", 280011000 } }
local DEEP_ROOMS = { { "in02", 280011002 }, { "in03", 280011003 }, { "in05", 280011005 } }
local DEEP_BACK = { { "out00", 280011000 }, { "st00", MINE } }

local function walk(ctx, bot, path)
	for _, step in ipairs(path) do
		if pq.portal(ctx, bot, step[1], step[2]) == false then
			return false
		end
	end
	return true
end

local function open_key_boxes(ctx, bot)
	return pq.visit_reactors(ctx, bot, function(r)
		return KEY_BOXES[r.id] ~= nil
	end, function(r)
		if pq.break_reactor(ctx, bot, r) == false then
			return false
		end
		return pq.loot_spawn(ctx, bot, KEY)
	end)
end

local function adobis(ctx, bot)
	local oid = pq.npc(ctx, bot, ADOBIS)
	if oid == false then
		return false
	end
	local dlg = bot:npc_click(oid)
	if dlg == false or dlg.selections == nil then
		return ctx:fail(bot:name() .. " 아도비스 선택지가 오지 않음")
	end
	return true
end

local function adobis_yes(ctx, bot, selected)
	if adobis(ctx, bot) == false then
		return false
	end
	local confirm = bot:dialog(true, selected)
	if confirm == nil then
		return ctx:fail(bot:name() .. " 아도비스 확인 대화가 오지 않음")
	end
	return true
end

test_suite {
	name = "Party Quest: 자쿰 PQ 1~3단계",
	bot_count = 2,

	on_initialize = function(ctx)
		for i = 0, ctx:bot_count() - 1 do
			local bot = ctx:bot(i)
			if pq.command(bot, "/봇초기화 50 100 0 4000051:30 100000:1", "봇초기화 완료") == false then
				return ctx:fail(bot:name() .. " 봇 초기화 실패")
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
			if adobis(ctx, leader) == false then
				return false
			end
			local warp = leader:request(resp.warp, req.dialog { dialog_type = DIALOG_LIST, next = true, selected = 0 }, function(p)
				return p.character.map == MINE
			end, 15000)
			if warp == false then
				return ctx:fail("폐광 동굴 입장 실패")
			end
			return pq.wait_all(ctx, MINE, 15000)
		end,
		function(ctx)
			local leader = ctx:bot(0)
			for i, cave in ipairs(CAVES) do
				if walk(ctx, leader, cave) == false then
					return false
				end
				if open_key_boxes(ctx, leader) == false then
					return false
				end
				if walk(ctx, leader, CAVE_BACK[i]) == false then
					return false
				end
			end
			if walk(ctx, leader, DEEP) == false then
				return false
			end
			for _, room in ipairs(DEEP_ROOMS) do
				if walk(ctx, leader, { room }) == false then
					return false
				end
				if open_key_boxes(ctx, leader) == false then
					return false
				end
				if room[2] ~= 280011005 and walk(ctx, leader, { DEEP_BACK[1] }) == false then
					return false
				end
			end
			if (leader:items()[KEY] or 0) < 7 then
				return ctx:fail("열쇠 부족: " .. (leader:items()[KEY] or 0))
			end
			if pq.feed(ctx, leader, CHEST, KEY, 7) == false then
				return false
			end
			if leader:drops(FIRE_ORE)[1] == nil and leader:request(resp.spawn_item, nil, function(p)
				return p.item_model ~= nil and p.item_model.id == FIRE_ORE
			end, 30000) == false then
				return ctx:fail("보물상자에서 불의 원석이 나오지 않음")
			end
			if pq.pick_up(ctx, leader, FIRE_ORE) == false then
				return false
			end
			return walk(ctx, leader, DEEP_BACK)
		end,
		function(ctx)
			local leader = ctx:bot(0)
			local oid = pq.npc(ctx, leader, AURA)
			if oid == false then
				return false
			end
			if leader:npc_click(oid) == false then
				return ctx:fail("아우라 대화가 오지 않음")
			end
			for _ = 1, 2 do
				if leader:dialog(true) == nil then
					return ctx:fail("아우라 설명이 이어지지 않음")
				end
			end
			if leader:dialog(true) == nil then
				return ctx:fail("아우라 선택지가 오지 않음")
			end
			local confirm = leader:dialog(true, 0)
			if confirm == nil or confirm.text:find("불의 원석", 1, true) == nil then
				return ctx:fail("아우라가 불의 원석을 확인하지 않음")
			end
			local done = leader:dialog(true)
			if done == nil or done.text:find("1단계를 클리어", 1, true) == nil then
				return ctx:fail("아우라가 1단계를 클리어하지 않음: " .. (done and done.text or "응답 없음"))
			end
			leader:dialog(false)
			for i = 0, ctx:bot_count() - 1 do
				local bot = ctx:bot(i)
				if pq.portal(ctx, bot, "ps01", EXIT) == false then
					return false
				end
				if (bot:items()[FIRE_ORE_PIECE] or 0) < 1 then
					return ctx:fail(bot:name() .. " 불의 원석 조각을 받지 못함")
				end
			end
			return true
		end,
		function(ctx)
			for i = 0, ctx:bot_count() - 1 do
				local bot = ctx:bot(i)
				local oid = pq.npc(ctx, bot, ALI)
				if oid == false then
					return false
				end
				if bot:npc_click(oid) == false then
					return ctx:fail(bot:name() .. " 알리 대화가 오지 않음")
				end
				local warp = bot:request(resp.warp, req.dialog { dialog_type = DIALOG_DEFAULT, next = true }, function(p)
					return p.character.map == ENTRY
				end, 10000)
				if warp == false then
					return ctx:fail(bot:name() .. " 알리가 내보내지 않음")
				end
			end
			return true
		end,
		function(ctx)
			for i = 0, ctx:bot_count() - 1 do
				local bot = ctx:bot(i)
				if adobis_yes(ctx, bot, 1) == false then
					return false
				end
				if bot:dialog(true) == nil then
					return ctx:fail(bot:name() .. " 2단계 안내가 오지 않음")
				end
				local warp = bot:request(resp.warp, req.dialog { dialog_type = DIALOG_DEFAULT, next = true }, function(p)
					return p.character.map == BREATH_ROAD
				end, 10000)
				if warp == false then
					return ctx:fail(bot:name() .. " 화산의 숨결 1단계 입장 실패")
				end
				if pq.portal(ctx, bot, "east00", BREATH_END) == false then
					return false
				end
				local oid = pq.npc(ctx, bot, LIRA)
				if oid == false then
					return false
				end
				if bot:npc_click(oid) == false then
					return ctx:fail(bot:name() .. " 리라 대화가 오지 않음")
				end
				local home = bot:request(resp.warp, req.dialog { dialog_type = DIALOG_DEFAULT, next = true }, function(p)
					return p.character.map == ENTRY
				end, 10000)
				if home == false then
					return ctx:fail(bot:name() .. " 리라가 돌려보내지 않음")
				end
				if (bot:items()[VOLCANO_BREATH] or 0) < 1 then
					return ctx:fail(bot:name() .. " 화산의 숨결을 받지 못함")
				end
			end
			return true
		end,
		function(ctx)
			for i = 0, ctx:bot_count() - 1 do
				local bot = ctx:bot(i)
				if adobis_yes(ctx, bot, 2) == false then
					return false
				end
				local ask = bot:dialog(true)
				if ask == nil or ask.text:find("헥터의 꼬리 30개", 1, true) == nil then
					return ctx:fail(bot:name() .. " 아도비스가 헥터의 꼬리를 요구하지 않음")
				end
				bot:dialog(false)
				if adobis_yes(ctx, bot, 2) == false then
					return false
				end
				if bot:dialog(true) == nil then
					return ctx:fail(bot:name() .. " 제련 대화가 오지 않음")
				end
				local done = bot:dialog(true)
				if done == nil or done.text:find("불의 눈", 1, true) == nil then
					return ctx:fail(bot:name() .. " 제련 결과 대화가 오지 않음")
				end
				bot:dialog(false)
				if (bot:items()[EYE_OF_FIRE] or 0) < 5 then
					return ctx:fail(bot:name() .. " 불의 눈을 받지 못함: " .. (bot:items()[EYE_OF_FIRE] or 0))
				end
			end
			return true
		end,
	},
}
