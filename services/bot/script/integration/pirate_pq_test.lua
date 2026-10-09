local pq = require("script/integration/lib/party_quest")

local ENTRY = 251010404
local STAGE1 = 925100000
local BOW = 925100100
local DECK1 = 925100200
local TREASURE1 = 925100201
local HIDEOUT1 = 925100202
local DECK2 = 925100300
local TREASURE2 = 925100301
local HIDEOUT2 = 925100302
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
local TREASURE = 2512001
local TREASURE_KEY = 4031437
local FURIOUS_DAVY = 9300106
local HIDEOUT_SECONDS = 5
local GUARD_SPAWNS = { { x = 0, y = 238 }, { x = 1700, y = 238 } }
local PIRATE_SPAWNS = { { x = 430, y = 75 }, { x = 1600, y = 75 }, { x = 430, y = 238 }, { x = 1600, y = 238 } }
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

local function sweep(ctx, bot, spawns)
	local spots = bot:mob_spots()
	for _, spawn in ipairs(spawns) do
		spots[#spots + 1] = spawn
	end
	for _, spot in ipairs(spots) do
		if pq.move(bot, spot.x, spot.y) == false then
			return ctx:fail(bot:name() .. " 몹 위치로 이동 실패")
		end
		if clear_mobs(ctx, bot) == false then
			return false
		end
	end
	return true
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

local function visit_hideout(ctx, x, y, hideout, deck)
	local leader = ctx:bot(0)
	if pq.move(leader, x, y) == false then
		return ctx:fail("숨겨진 방 포탈 앞으로 이동 실패")
	end
	leader:warp("in01")
	if pq.wait_all(ctx, hideout, 15000) == false then
		return false
	end
	if pq.command(leader, "/타이머 " .. HIDEOUT_SECONDS, "타이머 제한: " .. HIDEOUT_SECONDS) == false then
		return ctx:fail("타이머 제한 설정 실패")
	end
	if pq.command(leader, "/타이머 0", "타이머 제한: 0") == false then
		return ctx:fail("타이머 제한 해제 실패")
	end
	return pq.wait_all(ctx, deck, HIDEOUT_SECONDS * 1000 + 15000)
end

local function steal_treasure(ctx, x, y, room, deck)
	local leader = ctx:bot(0)
	if pq.move(leader, x, y) == false then
		return ctx:fail("보물방 포탈 앞으로 이동 실패")
	end
	if pq.portal(ctx, leader, "in00", room) == false then
		return false
	end
	local treasure = pq.seek_reactor(ctx, leader, pq.reactor_by_id(TREASURE))
	if treasure == false then
		return false
	end
	for _ = 1, 3 do
		if sweep(ctx, leader, GUARD_SPAWNS) == false then
			return false
		end
		local current = pq.find_reactor(leader, treasure.oid)
		if current ~= nil and current.state ~= 0 then
			break
		end
		leader:request(resp.trigger_reactor, nil, function(p)
			return p.reactor.oid == treasure.oid
		end, 3000)
	end
	local opened = pq.find_reactor(leader, treasure.oid)
	if opened == nil or opened.state == 0 then
		return ctx:fail("경비를 모두 처치해도 보물상자가 열리지 않음: " .. room)
	end
	if pq.command(leader, "/아이템생성 " .. TREASURE_KEY .. " 1", "아이템 생성") == false then
		return ctx:fail("보물상자 열쇠 생성 실패")
	end
	if pq.feed(ctx, leader, TREASURE, TREASURE_KEY, 1) == false then
		return false
	end
	if pq.move(leader, 744, 234) == false then
		return ctx:fail("보물방 출구로 이동 실패")
	end
	return pq.portal(ctx, leader, "out00", deck)
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
			return visit_hideout(ctx, 1653, 235, HIDEOUT1, DECK1)
		end,
		function(ctx)
			if steal_treasure(ctx, 523, 229, TREASURE1, DECK1) == false then
				return false
			end
			if sweep(ctx, ctx:bot(0), PIRATE_SPAWNS) == false then
				return false
			end
			return pass_all(ctx, 2133, 240, DECK2)
		end,
		function(ctx)
			return visit_hideout(ctx, 436, 233, HIDEOUT2, DECK2)
		end,
		function(ctx)
			if steal_treasure(ctx, 1653, 229, TREASURE2, DECK2) == false then
				return false
			end
			if sweep(ctx, ctx:bot(0), PIRATE_SPAWNS) == false then
				return false
			end
			return pass_all(ctx, 2133, 240, DOORS)
		end,
		function(ctx)
			local leader = ctx:bot(0)
			if pq.collect(ctx, leader, nil, KEY, 4) == false then
				return false
			end
			for i = 0, 3 do
				local door = pq.seek_reactor(ctx, leader, pq.reactor_by_id(DOOR + i))
				if door == false then
					return false
				end
				if door.state == 0 and pq.feed(ctx, leader, DOOR + i, KEY, 1) == false then
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
			if leader:mobs(FURIOUS_DAVY)[1] == nil and leader:request(resp.spawn_mob, nil, function(p)
				return p.mob.mob_id == FURIOUS_DAVY
			end, 10000) == false then
				return ctx:fail("보물을 두 번 훔쳤는데 몹시 화난 데비존이 나타나지 않음")
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
