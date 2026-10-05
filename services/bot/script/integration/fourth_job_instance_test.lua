local pq = require("script/integration/lib/party_quest")

local function setup(ctx, bot, args, map_id)
	if pq.command(bot, "/봇초기화 120 " .. args, "봇초기화 완료") == false then
		return ctx:fail("봇 초기화 실패: " .. args)
	end
	if bot:map_move(map_id) == false then
		return ctx:fail("맵 이동 실패: " .. map_id)
	end
	return true
end

local function pass(ctx, bot, x, y, portal, map_id)
	if pq.move(bot, x, y) == false then
		return ctx:fail("포탈 앞으로 이동 실패: " .. portal)
	end
	local message = ""
	for _ = 1, 10 do
		local p, name = bot:request({ resp.warp, resp.notice }, req.warp { target = 4294967295, portal_name = portal })
		if name ~= resp.notice then
			return pq.wait_map(ctx, bot, map_id, 5000)
		end
		message = p.message
		ctx:sleep(500)
	end
	return ctx:fail("포탈 입장 거부: " .. portal .. " " .. message)
end

local function offer(ctx, bot, reactor, item_id)
	if bot:drop(item_id, 1) == nil then
		return ctx:fail("아이템 버리기 실패: " .. item_id)
	end
	local p = bot:request({ resp.trigger_reactor, resp.destroy_reactor }, nil, function(p)
		return p.reactor.oid == reactor.oid
	end, 12000)
	if p == false then
		return ctx:fail(string.format("리액터가 아이템에 반응하지 않음: %d <- %d", reactor.id, item_id))
	end
	return true
end

local function hatch(ctx, bot, args, reactor_id, egg_id, mob_id)
	if setup(ctx, bot, args, 240010700) == false then
		return false
	end
	if pass(ctx, bot, -397, -462, "pt00", 924000100) == false then
		return false
	end
	if pq.feed(ctx, bot, reactor_id, egg_id, 1) == false then
		return false
	end
	local reactor = pq.seek_reactor(ctx, bot, pq.reactor_by_id(reactor_id))
	if reactor == false then
		return false
	end
	if pq.break_reactor(ctx, bot, reactor) == false then
		return false
	end
	if pq.kill_at(ctx, bot, mob_id, reactor.x, reactor.y) == false then
		return false
	end
	return pass(ctx, bot, 316, 422, "out00", 240010700)
end

test_suite {
	name = "Quest: 4차 스킬 인스턴스",
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
			local bot = ctx:bot(0)
			if setup(ctx, bot, "232 0 4001108:1,4001107:1 6131:2,6132:1", 230040300) == false then
				return false
			end
			if pass(ctx, bot, 123, 1375, "in01", 923000100) == false then
				return false
			end
			local altar = pq.seek_reactor(ctx, bot, pq.reactor_by_id(2302003))
			if altar == false then
				return false
			end
			if pq.move(bot, altar.x, altar.y - 89) == false then
				return ctx:fail("제단 위로 이동 실패")
			end
			if offer(ctx, bot, altar, 4001108) == false then
				return false
			end
			if offer(ctx, bot, altar, 4001107) == false then
				return false
			end
			if pq.pick_up(ctx, bot, 4161017) == false then
				return false
			end
			return pass(ctx, bot, -334, 167, "out00", 230040300)
		end,
		function(ctx)
			local bot = ctx:bot(0)
			if setup(ctx, bot, "232 0 - 6133:2,6134:1", 220070400) == false then
				return false
			end
			if pass(ctx, bot, -791, -130, "pt00", 922020000) == false then
				return false
			end
			if pass(ctx, bot, 700, -405, "out00", 220070400) == false then
				return false
			end
			if (bot:items()[4031448] or 0) == 0 then
				return ctx:fail("죽음의 체험 증표를 받지 못함")
			end
			return true
		end,
		function(ctx)
			local bot = ctx:bot(0)
			if setup(ctx, bot, "312 0 - 6240:1", 211042200) == false then
				return false
			end
			if pass(ctx, bot, -374, -277, "in00", 921100200) == false then
				return false
			end
			if pq.open_box(ctx, bot, 2112016, 4001113) == false then
				return false
			end
			return pass(ctx, bot, -291, -244, "out00", 211042200)
		end,
		function(ctx)
			return hatch(ctx, ctx:bot(0), "312 0 4001113:1 6240:2,6241:1", 2401001, 4001113, 9300089)
		end,
		function(ctx)
			return hatch(ctx, ctx:bot(0), "322 0 4001114:1 6242:2,6243:1", 2401002, 4001114, 9300090)
		end,
	},
}
