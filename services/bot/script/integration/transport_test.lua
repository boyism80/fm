local pq = require("script/integration/lib/party_quest")

local DIALOGS = { resp.dialog, resp.dialog_yes_no, resp.dialog_list, resp.warp }
local BALROG = 8150000

local function answer(p, name)
	if name == resp.dialog_yes_no then
		return req.dialog { dialog_type = 1, next = true }
	elseif name == resp.dialog_list then
		return req.dialog { dialog_type = 4, next = true, selected = 0 }
	end
	return req.dialog { dialog_type = 0, next = true }
end

local function talk_until_warp(ctx, bot, template, map_id)
	local oid = pq.npc(ctx, bot, template)
	if oid == false then
		return false
	end
	local p, name = bot:request(DIALOGS, req.npc_click { oid = oid })
	for _ = 1, 10 do
		if p == false then
			return ctx:fail(string.format("NPC %d 대화 중 응답 없음 (맵 %d)", template, bot:map()))
		end
		if name == resp.warp then
			if p.character.map ~= map_id then
				return ctx:fail(string.format("NPC %d 대화 뒤 다른 맵으로 이동: %d", template, p.character.map))
			end
			return true
		end
		p, name = bot:request(DIALOGS, answer(p, name))
	end
	return ctx:fail(string.format("NPC %d 대화가 끝나지 않음", template))
end

local function phase(ctx, bot, group, name, map_id)
	if pq.command(bot, "/수송 " .. group .. " " .. name, "강제 호출") == false then
		return ctx:fail(string.format("수송 단계 실패: %s %s", group, name))
	end
	if map_id == nil then
		return true
	end
	return pq.wait_map(ctx, bot, map_id, 5000)
end

local function repel(ctx, bot)
	if phase(ctx, bot, "Boats", "invasion") == false then
		return false
	end
	while #bot:mobs(BALROG) < 2 do
		if bot:request(resp.spawn_mob, nil, function(p)
				return p.mob.mob_id == BALROG
			end, 5000) == false then
			return ctx:fail("배에 발록이 나타나지 않음: " .. #bot:mobs(BALROG))
		end
	end
	for _, mob in ipairs(bot:mobs(BALROG)) do
		if bot:kill(mob.oid) == false then
			return ctx:fail("발록 처치 실패")
		end
	end
	return true
end

local function ride(ctx, route)
	local bot = ctx:bot(0)
	if pq.command(bot, "/봇초기화 40 100 0 " .. route.ticket .. ":1 -", "봇초기화 완료") == false then
		return ctx:fail("봇 초기화 실패")
	end
	if bot:map_move(route.station) == false then
		return ctx:fail("정거장으로 이동 실패: " .. route.station)
	end
	if phase(ctx, bot, route.group, "dock") == false then
		return false
	end
	if talk_until_warp(ctx, bot, route.npc, route.waiting) == false then
		return false
	end
	if (bot:items()[route.ticket] or 0) ~= 0 then
		return ctx:fail("탑승권이 회수되지 않음: " .. route.ticket)
	end
	if phase(ctx, bot, route.group, "takeoff", route.ride) == false then
		return false
	end
	if route.group == "Boats" and repel(ctx, bot) == false then
		return false
	end
	return phase(ctx, bot, route.group, "arrived", route.dest)
end

local function pass(ctx, bot, x, y, map_id)
	if pq.move(bot, x, y) == false then
		return ctx:fail("엘리베이터 앞으로 이동 실패")
	end
	return pq.portal(ctx, bot, "in00", map_id)
end

test_suite {
	name = "Transport: 배·비행선·지니·열차·엘리베이터",
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
			return ride(ctx, { group = "Boats", station = 200000111, npc = 2012001, ticket = 4031047, waiting = 200000112, ride = 200090000, dest = 101000300 })
		end,
		function(ctx)
			return ride(ctx, { group = "Boats", station = 101000300, npc = 1032008, ticket = 4031045, waiting = 101000301, ride = 200090010, dest = 200000100 })
		end,
		function(ctx)
			return ride(ctx, { group = "Flight", station = 200000131, npc = 2012021, ticket = 4031331, waiting = 200000132, ride = 200090200, dest = 240000100 })
		end,
		function(ctx)
			return ride(ctx, { group = "Flight", station = 240000110, npc = 2082001, ticket = 4031045, waiting = 240000111, ride = 200090210, dest = 200000100 })
		end,
		function(ctx)
			return ride(ctx, { group = "Geenie", station = 200000151, npc = 2012025, ticket = 4031576, waiting = 200000152, ride = 200090400, dest = 260000100 })
		end,
		function(ctx)
			return ride(ctx, { group = "Geenie", station = 260000100, npc = 2102000, ticket = 4031045, waiting = 260000110, ride = 200090410, dest = 200000100 })
		end,
		function(ctx)
			return ride(ctx, { group = "Trains", station = 200000121, npc = 2012013, ticket = 4031074, waiting = 200000122, ride = 200090100, dest = 220000110 })
		end,
		function(ctx)
			return ride(ctx, { group = "Trains", station = 220000110, npc = 2041000, ticket = 4031074, waiting = 220000111, ride = 200090110, dest = 200000121 })
		end,
		function(ctx)
			local bot = ctx:bot(0)
			if pq.command(bot, "/봇초기화 40 100 0 - -", "봇초기화 완료") == false then
				return ctx:fail("봇 초기화 실패")
			end
			if bot:map_move(222020100) == false then
				return ctx:fail("헬리오스 탑 2층으로 이동 실패")
			end
			if phase(ctx, bot, "elevator", "waiting_to_up") == false then
				return false
			end
			if pass(ctx, bot, -139, 286, 222020110) == false then
				return false
			end
			if phase(ctx, bot, "elevator", "run_to_up", 222020111) == false then
				return false
			end
			if phase(ctx, bot, "elevator", "waiting_to_down", 222020200) == false then
				return false
			end
			if pass(ctx, bot, -133, 1963, 222020210) == false then
				return false
			end
			if phase(ctx, bot, "elevator", "run_to_down", 222020211) == false then
				return false
			end
			return phase(ctx, bot, "elevator", "waiting_to_up", 222020100)
		end,
	},
}
