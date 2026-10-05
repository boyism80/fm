local pq = require("script/integration/lib/party_quest")

local ROMEO = {
	name = "로미오",
	town = 261000011,
	base = 926100000,
	guide = 2112004,
	hub_guide = 2112006,
	investigate = 2112007,
	letter = 4001131,
	marble = 4001159,
	urete = 2112002,
}
local JULIET = {
	name = "줄리엣",
	town = 261000021,
	base = 926110000,
	guide = 2112003,
	hub_guide = 2112005,
	investigate = 2112013,
	letter = 4001130,
	marble = 4001160,
	urete = 2112012,
}
local HALLWAY = 1
local BEAKERS = 100
local HUB = 200
local LAB1 = 201
local LAB2 = 202
local OFFICE = 203
local TOWER = 300
local SUMMIT = 400
local BOSS_ROOM = 401
local URETE_ROOM = 500
local REWARD = 600
local EXIT = 700
local BEAKER = 2618000
local BEAKER_FILL = 7
local LIQUID = 4001132
local LIQUID_MOB = 9300147
local KEY = 4001133
local KEY_MOB = 9300148
local LAB1_DOOR = 2618002
local LAB2_DOOR = 2618001
local LAB1_BOX = 2612002
local LAB2_BOX = 2612001
local LAB1_NOTE = 4001134
local LAB2_NOTE = 4001135
local URETE = 2112000
local BOSS_URETE = 2112010
local PERSUADED_URETE = 9300140
local INVESTIGATE_RANGE = 5000
local DIALOG_LIST = 4
local ANY_PORTAL = 4294967295

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

local function converse(ctx, bot, oid, map_id)
	local function arrived(p, name)
		return name == resp.dialog or p.character.map == map_id
	end
	local p, name = bot:request({ resp.dialog, resp.warp }, req.npc_click { oid = oid }, arrived, 15000)
	for _ = 1, 5 do
		if p == false then
			return ctx:fail(bot:name() .. " 대화 후 맵 이동 없음: " .. map_id)
		end
		if name == resp.warp then
			return true
		end
		p, name = bot:request({ resp.dialog, resp.warp }, req.dialog { dialog_type = 0, next = true }, arrived, 15000)
	end
	return ctx:fail(bot:name() .. " 대화가 끝나지 않음: " .. map_id)
end

local function start(ctx, path)
	local leader = ctx:bot(0)
	local oid = pq.npc(ctx, leader, path.guide)
	if oid == false then
		return false
	end
	local dlg = leader:npc_click(oid)
	if dlg == false or dlg.selections == nil then
		return ctx:fail(path.name .. " 선택지가 오지 않음")
	end
	local warp = leader:request(resp.warp, req.dialog { dialog_type = DIALOG_LIST, next = true, selected = 2 }, function(p)
		return p.character.map == path.base
	end, 15000)
	if warp == false then
		return ctx:fail(path.name .. " 파티 퀘스트 시작 실패")
	end
	return pq.wait_all(ctx, path.base, 15000)
end

local function investigate(ctx, path)
	local leader = ctx:bot(0)
	local opened = false
	for _, spot in ipairs(leader:npc_spots()) do
		if opened and (leader:items()[path.letter] or 0) > 0 then
			break
		end
		if spot.id == path.investigate then
			if pq.move(leader, spot.x, spot.y) == false then
				return ctx:fail("조사 위치로 이동 실패")
			end
			for _, npc in ipairs(leader:npcs(path.investigate)) do
				local dx = npc.x - spot.x
				local dy = npc.y - spot.y
				if dx * dx + dy * dy <= INVESTIGATE_RANGE then
					local p, name = leader:request({ resp.dialog, resp.notice }, req.npc_click { oid = npc.oid }, function(p, name)
						return name == resp.dialog or p.message:find("포탈", 1, true) ~= nil
					end, 5000)
					if p == false then
						return ctx:fail("조사 응답 없음: " .. npc.oid)
					end
					if name == resp.dialog then
						leader:dialog(false)
					else
						opened = true
					end
				end
			end
		end
	end
	if opened == false then
		return ctx:fail(path.name .. " 모든 곳을 조사해도 포탈이 열리지 않음")
	end
	if (leader:items()[path.letter] or 0) == 0 then
		return ctx:fail(path.name .. " 편지를 찾지 못함")
	end
	return pass_all(ctx, "pt00", path.base + HALLWAY)
end

local function clear_hallway(ctx, path)
	if sweep(ctx, ctx:bot(0)) == false then
		return false
	end
	return pass_all(ctx, "pt00", path.base + BEAKERS)
end

local function fill_beakers(ctx, path)
	local leader = ctx:bot(0)
	if pq.collect(ctx, leader, LIQUID_MOB, LIQUID, BEAKER_FILL * 3) == false then
		return false
	end
	local filled = 0
	local visited = pq.visit_reactors(ctx, leader, pq.reactor_by_id(BEAKER), function(beaker)
		for _ = 1, BEAKER_FILL do
			if leader:drop(LIQUID, 1) == nil then
				return ctx:fail("수상한 액체 버리기 실패")
			end
			local p = leader:request(resp.trigger_reactor, nil, function(p)
				return p.reactor.oid == beaker.oid
			end, 12000)
			if p == false then
				return ctx:fail("비커가 액체에 반응하지 않음: " .. beaker.oid)
			end
		end
		filled = filled + 1
		return true
	end)
	if visited == false then
		return false
	end
	if filled ~= 3 then
		return ctx:fail("비커 3개를 모두 찾지 못함: " .. filled)
	end
	return pass_all(ctx, "pt00", path.base + HUB)
end

local function search_lab(ctx, path, door, portal, lab, box, note)
	local leader = ctx:bot(0)
	if pq.feed(ctx, leader, door, KEY, 1) == false then
		return false
	end
	if pq.portal(ctx, leader, portal, path.base + lab) == false then
		return false
	end
	if pq.open_box(ctx, leader, box, note) == false then
		return false
	end
	return pq.portal(ctx, leader, "out00", path.base + HUB)
end

local function search_labs(ctx, path)
	local leader = ctx:bot(0)
	if pq.collect(ctx, leader, KEY_MOB, KEY, 2) == false then
		return false
	end
	if search_lab(ctx, path, LAB1_DOOR, "in00", LAB1, LAB1_BOX, LAB1_NOTE) == false then
		return false
	end
	if search_lab(ctx, path, LAB2_DOOR, "in01", LAB2, LAB2_BOX, LAB2_NOTE) == false then
		return false
	end
	if pq.talk(ctx, leader, path.hub_guide, "편지") == false then
		return false
	end
	if pq.talk(ctx, leader, path.hub_guide, "자료를 찾아오셨군요") == false then
		return false
	end
	if pq.talk(ctx, leader, path.hub_guide, "자료도 찾아오셨군요") == false then
		return false
	end
	return pass_all(ctx, "pt00", path.base + OFFICE)
end

local function break_into_office(ctx, path)
	local leader = ctx:bot(0)
	if pq.move(leader, 200, 188) == false then
		return ctx:fail("유레테 앞으로 이동 실패")
	end
	local oid = leader:npc(URETE)
	if oid == nil then
		return ctx:fail("연구실에 유레테가 없음")
	end
	if leader:request(resp.notice, req.npc_click { oid = oid }, nil, 5000) == false then
		return ctx:fail("유레테가 반응하지 않음")
	end
	if pq.move(leader, 68, 214) == false then
		return ctx:fail("연구실 안쪽으로 이동 실패")
	end
	if leader:request(resp.spawn_mob, req.warp { target = ANY_PORTAL, portal_name = "pt00" }, nil, 10000) == false then
		return ctx:fail("연구실에 몹이 나타나지 않음")
	end
	if clear_mobs(ctx, leader) == false then
		return false
	end
	return pass_all(ctx, "out00", path.base + TOWER)
end

local function climb(ctx, bot)
	for row = 0, 9 do
		local passed = false
		for column = 0, 3 do
			local portal = string.format("pt%d%d", row, column)
			if bot:request(resp.field_relocate, req.warp { target = ANY_PORTAL, portal_name = portal }, nil, 5000) == false then
				return ctx:fail(bot:name() .. " 탑 포탈 반응 없음: " .. portal)
			end
			if bot:spawn_portal() == "np0" .. row then
				passed = true
				break
			end
		end
		if passed == false then
			return ctx:fail(bot:name() .. " 탑 " .. row .. "층에 맞는 포탈이 없음")
		end
	end
	return true
end

local function climb_towers(ctx, path)
	for i = 0, ctx:bot_count() - 1 do
		local bot = ctx:bot(i)
		if pq.portal(ctx, bot, "pt0" .. i, path.base + TOWER + 1 + i) == false then
			return false
		end
		if climb(ctx, bot) == false then
			return false
		end
		if pq.portal(ctx, bot, "out00", path.base + SUMMIT) == false then
			return false
		end
	end
	if pq.portal(ctx, ctx:bot(0), "pt00", path.base + BOSS_ROOM) == false then
		return false
	end
	return pq.wait_all(ctx, path.base + BOSS_ROOM, 15000)
end

local function persuade_urete(ctx, path)
	local leader = ctx:bot(0)
	if pq.move(leader, 242, 150) == false then
		return ctx:fail("유레테 앞으로 이동 실패")
	end
	local oid = leader:npc(BOSS_URETE)
	if oid == nil then
		return ctx:fail("보스방에 유레테가 없음")
	end
	local dlg = leader:npc_click(oid)
	if dlg == false or dlg.selections == nil then
		return ctx:fail("편지를 전했는데 유레테 설득 선택지가 없음")
	end
	local spawn, name = leader:request({ resp.dialog, resp.spawn_mob }, req.dialog { dialog_type = DIALOG_LIST, next = true, selected = 1 }, function(p, name)
		return name == resp.dialog or p.mob.mob_id == PERSUADED_URETE
	end, 15000)
	for _ = 1, 5 do
		if spawn == false then
			return ctx:fail("설득한 유레테가 나타나지 않음")
		end
		if name == resp.spawn_mob then
			break
		end
		spawn, name = leader:request({ resp.dialog, resp.spawn_mob }, req.dialog { dialog_type = 0, next = true }, function(p, name)
			return name == resp.dialog or p.mob.mob_id == PERSUADED_URETE
		end, 15000)
	end
	if name ~= resp.spawn_mob then
		return ctx:fail("유레테 대화가 끝나지 않음")
	end
	if leader:kill(spawn.mob.oid) == false then
		return ctx:fail("유레테 처치 실패")
	end

	local guide = leader:npc(path.guide, 10000)
	if guide == nil then
		return ctx:fail(path.name .. " 보호에 실패해 " .. path.name .. "가 나타나지 않음")
	end
	if converse(ctx, leader, guide, path.base + URETE_ROOM) == false then
		return false
	end
	if pq.wait_all(ctx, path.base + URETE_ROOM, 15000) == false then
		return false
	end
	if pq.move(leader, 232, 150) == false then
		return ctx:fail("유레테 앞으로 이동 실패")
	end
	local urete = leader:npc(path.urete)
	if urete == nil then
		return ctx:fail(path.name .. " 설득에 성공했는데 유레테가 없음")
	end
	if converse(ctx, leader, urete, path.base + REWARD) == false then
		return false
	end
	return pq.wait_all(ctx, path.base + REWARD, 15000)
end

local function reward(ctx, path)
	for i = 0, ctx:bot_count() - 1 do
		local bot = ctx:bot(i)
		local before = bot:items()[path.marble] or 0
		if pq.move(bot, 107, 128) == false then
			return ctx:fail(bot:name() .. " " .. path.name .. " 앞으로 이동 실패")
		end
		local oid = bot:npc(path.guide)
		if oid == nil then
			return ctx:fail(bot:name() .. " 보상 맵에 " .. path.name .. "가 없음")
		end
		if converse(ctx, bot, oid, path.base + EXIT) == false then
			return false
		end
		if (bot:items()[path.marble] or 0) ~= before + 1 then
			return ctx:fail(bot:name() .. " 구슬을 받지 못함: " .. path.marble)
		end
		if pq.portal(ctx, bot, "out00", path.town) == false then
			return false
		end
	end
	return true
end

local function enter_town(ctx, path)
	for i = 0, ctx:bot_count() - 1 do
		if ctx:bot(i):map_move(path.town) == false then
			return false
		end
	end
	return true
end

local function stages(path)
	return {
		function(ctx)
			if enter_town(ctx, path) == false then
				return false
			end
			return start(ctx, path)
		end,
		function(ctx)
			return investigate(ctx, path)
		end,
		function(ctx)
			return clear_hallway(ctx, path)
		end,
		function(ctx)
			return fill_beakers(ctx, path)
		end,
		function(ctx)
			return search_labs(ctx, path)
		end,
		function(ctx)
			return break_into_office(ctx, path)
		end,
		function(ctx)
			return climb_towers(ctx, path)
		end,
		function(ctx)
			return persuade_urete(ctx, path)
		end,
		function(ctx)
			return reward(ctx, path)
		end,
	}
end

local scenarios = { pq.form_party }
for _, path in ipairs({ ROMEO, JULIET }) do
	for _, stage in ipairs(stages(path)) do
		scenarios[#scenarios + 1] = stage
	end
end

test_suite {
	name = "Party Quest: 마가티아 PQ",
	bot_count = 4,

	on_initialize = function(ctx)
		for i = 0, ctx:bot_count() - 1 do
			local bot = ctx:bot(i)
			if pq.command(bot, "/레벨바꾸기 71", "레벨 설정") == false then
				return ctx:fail(bot:name() .. " 레벨 설정 실패")
			end
			if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
				return ctx:fail(bot:name() .. " 플레이어 모드 설정 실패")
			end
			if pq.command(bot, "/즉사", "즉사 상태: enabled") == false then
				return ctx:fail(bot:name() .. " 즉사 설정 실패")
			end
		end
		return true
	end,

	on_finished = function(ctx)
		ctx:bot(0):request(resp.party_update_disband, req.party_operation { operation = PARTY.Leave }, nil, 5000)
	end,

	scenarios = scenarios,
}
