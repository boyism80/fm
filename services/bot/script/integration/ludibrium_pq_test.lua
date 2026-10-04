local pq = require("script/integration/lib/party_quest")

local LOBBY = 221024500
local ENTRY_NPC = 2040034
local START = 922010100
local STAGE2_HIDDEN = 922010201
local CLIMB_MAP = 922010600
local AREA_MAP = 922010800
local BONUS = 922011000
local REWARD = 922011100
local PASS = 4001022
local KEY = 4001023
local WARP_BOX = 2200002
local ROMBARD_REACTOR = 2201002
local ALISHAR_REACTOR = 2201003
local ALISHAR = 9300012
local DIALOG_LIST = 4
local OPENED = "다음 스테이지로 통하는 포탈이 열렸습니다"

local CLIMB = {
	"h002", "h005", "h006", "h010", "h013", "h015", "h018", "h021",
	"h025", "h027", "h032", "h034", "h036", "h039", "h044",
}

local function climb(ctx, bot)
	for _, portal in ipairs(CLIMB) do
		if pq.portal(ctx, bot, portal, CLIMB_MAP) == false then
			return false
		end
	end
	return pq.portal(ctx, bot, "next00", CLIMB_MAP + 100)
end

local function catch_up(ctx, bot, map_id)
	while bot:map() < map_id do
		if bot:map() == CLIMB_MAP then
			if climb(ctx, bot) == false then
				return false
			end
		elseif pq.portal(ctx, bot, "next00", bot:map() + 100) == false then
			return false
		end
	end
	return true
end

local function sweep_rooms(ctx, bot, rooms, count)
	local hall = bot:map()
	for i = 1, rooms do
		if (bot:items()[PASS] or 0) >= count then
			return true
		end
		if pq.portal(ctx, bot, string.format("in%02d", i), hall + i) == false then
			return false
		end
		if pq.sweep(ctx, bot, PASS, count) == false then
			return false
		end
		if pq.portal(ctx, bot, "out00", hall) == false then
			return false
		end
	end
	return true
end

local function find_reactor(ctx, bot, template)
	for _, spot in ipairs(bot:reactor_spots()) do
		if spot.id == template then
			if pq.move(bot, spot.x, spot.y) == false then
				return ctx:fail("리액터 앞으로 이동 실패: " .. template)
			end
			local reactor = bot:reactors(template)[1]
			if reactor == nil then
				return ctx:fail("리액터가 보이지 않음: " .. template)
			end
			return reactor
		end
	end
	return ctx:fail("맵에 리액터 배치가 없음: " .. template)
end

local function pass_stage(ctx, npc, guide, count, collect)
	local leader = ctx:bot(0)
	if pq.talk(ctx, leader, npc, guide) == false then
		return false
	end
	if collect(leader) == false then
		return false
	end
	local have = leader:items()[PASS] or 0
	if have < count then
		return ctx:fail(string.format("통행증 부족: %d/%d (맵 %d)", have, count, leader:map()))
	end
	if pq.talk(ctx, leader, npc, OPENED) == false then
		return false
	end
	return pq.portal(ctx, leader, "next00", leader:map() + 100)
end

local function summon_alishar(ctx, leader)
	for _ = 1, 3 do
		local alishar = leader:mobs(ALISHAR)[1]
		if alishar ~= nil then
			return alishar.oid
		end
		local oid = leader:drop(PASS, 1)
		if oid == nil then
			return ctx:fail("통행증 버리기 실패")
		end
		local spawn = leader:request(resp.spawn_mob, nil, function(p)
			return p.mob.mob_id == ALISHAR
		end, 8000)
		if spawn ~= false then
			return spawn.mob.oid
		end
		if leader:loot(oid) == false then
			return ctx:fail("통행증 다시 줍기 실패")
		end
	end
	return ctx:fail("알리샤르가 나타나지 않음")
end

test_suite {
	name = "Party Quest: 루디브리엄 PQ",
	bot_count = 5,

	on_initialize = function(ctx)
		for i = 0, ctx:bot_count() - 1 do
			local bot = ctx:bot(i)
			if pq.command(bot, "/레벨바꾸기 35", "레벨 설정") == false then
				return ctx:fail(bot:name() .. " 레벨 설정 실패")
			end
			if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
				return ctx:fail(bot:name() .. " 플레이어 모드 설정 실패")
			end
			if bot:map_move(LOBBY) == false then
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
			local oid = pq.npc(ctx, leader, ENTRY_NPC)
			if oid == false then
				return false
			end
			local dlg = leader:npc_click(oid)
			if dlg == false or dlg.selections == nil then
				return ctx:fail("표지판 선택지가 오지 않음")
			end
			local warp = leader:request(resp.warp, req.dialog { dialog_type = DIALOG_LIST, next = true, selected = #dlg.selections - 1 }, function(p)
				return p.character.map == START
			end, 15000)
			if warp == false then
				return ctx:fail("1스테이지 입장 실패")
			end
			for i = 1, ctx:bot_count() - 1 do
				if pq.wait_map(ctx, ctx:bot(i), START, 15000) == false then
					return false
				end
			end
			return true
		end,
		function(ctx)
			return pass_stage(ctx, 2040036, "첫번째 스테이지", 25, function(leader)
				return pq.sweep(ctx, leader, PASS, 25)
			end)
		end,
		function(ctx)
			return pass_stage(ctx, 2040037, "두번째 스테이지", 15, function(leader)
				local hall = leader:map()
				if pq.sweep(ctx, leader, PASS, 15, { [WARP_BOX] = true }) == false then
					return false
				end
				local box = find_reactor(ctx, leader, WARP_BOX)
				if box == false then
					return false
				end
				for _ = 1, 10 do
					local p, name = leader:request({ resp.trigger_reactor, resp.warp }, req.damage_reactor { oid = box.oid, hit_side = 3 }, function(p, name)
						return name == resp.warp or p.reactor.oid == box.oid
					end)
					if p == false then
						return ctx:fail("추방 상자 타격 응답 없음")
					end
					if name == resp.warp then
						break
					end
				end
				if leader:map() ~= STAGE2_HIDDEN then
					return ctx:fail("추방 상자로 숨겨진 맵에 가지 못함: " .. leader:map())
				end
				if pq.sweep(ctx, leader, PASS, 15) == false then
					return false
				end
				return pq.portal(ctx, leader, "out00", hall)
			end)
		end,
		function(ctx)
			return pass_stage(ctx, 2040038, "세번째 스테이지", 32, function(leader)
				return pq.sweep(ctx, leader, PASS, 32)
			end)
		end,
		function(ctx)
			return pass_stage(ctx, 2040039, "네번째 스테이지", 6, function(leader)
				return sweep_rooms(ctx, leader, 5, 6)
			end)
		end,
		function(ctx)
			return pass_stage(ctx, 2040040, "다섯번째 스테이지", 24, function(leader)
				if pq.sweep(ctx, leader, PASS, 24) == false then
					return false
				end
				return sweep_rooms(ctx, leader, 6, 24)
			end)
		end,
		function(ctx)
			return climb(ctx, ctx:bot(0))
		end,
		function(ctx)
			return pass_stage(ctx, 2040042, "일곱번째 스테이지", 3, function(leader)
				return pq.sweep(ctx, leader, PASS, 3, { [ROMBARD_REACTOR] = true })
			end)
		end,
		function(ctx)
			for i = 1, ctx:bot_count() - 1 do
				if catch_up(ctx, ctx:bot(i), AREA_MAP) == false then
					return false
				end
			end

			local leader = ctx:bot(0)
			if pq.talk(ctx, leader, 2040043, "여덟번째 스테이지") == false then
				return false
			end
			local oid = leader:npc(2040043)
			if oid == nil then
				return ctx:fail("블루 벌룬이 보이지 않음")
			end
			local areas = leader:areas()
			for _, combo in ipairs(pq.combinations(#areas, 5)) do
				for j, index in ipairs(combo) do
					local bot = ctx:bot(j - 1)
					if pq.move(bot, areas[index].x, areas[index].y) == false then
						return ctx:fail(bot:name() .. " 발판 이동 실패: " .. index)
					end
				end
				local p, name = leader:request({ resp.dialog, resp.environment_change }, req.npc_click { oid = oid }, function(p, name)
					return name == resp.dialog or p.env:find("wrong", 1, true) ~= nil
				end)
				if p == false then
					return ctx:fail("발판 판정 응답 없음")
				end
				if name == resp.dialog then
					leader:dialog(false)
					if p.text:find(OPENED, 1, true) == nil then
						return ctx:fail("예상과 다른 판정 대화: " .. p.text)
					end
					for i = 0, ctx:bot_count() - 1 do
						if pq.portal(ctx, ctx:bot(i), "next00", AREA_MAP + 100) == false then
							return false
						end
					end
					return true
				end
			end
			return ctx:fail("정답 발판 조합을 찾지 못함")
		end,
		function(ctx)
			local leader = ctx:bot(0)
			if pq.talk(ctx, leader, 2040044, "드디어 여기까지") == false then
				return false
			end
			if pq.sweep(ctx, leader, PASS, 1, { [ALISHAR_REACTOR] = true }) == false then
				return false
			end
			if find_reactor(ctx, leader, ALISHAR_REACTOR) == false then
				return false
			end
			if (leader:items()[KEY] or 0) == 0 and #leader:drops(KEY) == 0 then
				local alishar = summon_alishar(ctx, leader)
				if alishar == false then
					return false
				end
				if leader:kill(alishar) == false then
					return ctx:fail("알리샤르 처치 실패")
				end
				leader:request(resp.spawn_item, nil, nil, 2000)
			end
			if pq.clear_visible(ctx, leader, KEY, { [ALISHAR_REACTOR] = true }) == false then
				return false
			end
			if (leader:items()[KEY] or 0) < 1 then
				return ctx:fail("차원의 열쇠를 얻지 못함")
			end
			if pq.talk(ctx, leader, 2040044, "차원의 열쇠를 가져 오셨군요") == false then
				return false
			end

			local oid = leader:npc(2040044)
			local warp = leader:request(resp.warp, req.npc_click { oid = oid }, function(p)
				return p.character.map == BONUS
			end)
			if warp == false then
				return ctx:fail("보너스 맵 이동 실패")
			end
			for i = 1, ctx:bot_count() - 1 do
				if pq.wait_map(ctx, ctx:bot(i), BONUS, 10000) == false then
					return false
				end
			end
			return true
		end,
		function(ctx)
			if pq.command(ctx:bot(0), "/타이머 1", "남은 시간") == false then
				return ctx:fail("보너스 시간 단축 실패")
			end
			for i = 0, ctx:bot_count() - 1 do
				if pq.wait_map(ctx, ctx:bot(i), REWARD, 15000) == false then
					return false
				end
			end
			return true
		end,
	},
}
