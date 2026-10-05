local pq = require("script/integration/lib/party_quest")

local ROOMS = 39
local DUNGEONS = {
	{ name = "pig", base = 100020000, room = 100020100, entry = { 720, -175 }, exit = { -47, -55 } },
	{ name = "golem", base = 105040304, room = 105040320, entry = { 717, 674 }, exit = { 268, 1156 } },
	{ name = "mushroom", base = 105050100, room = 105050101, entry = { 1287, 968 }, exit = { 394, -1001 } },
	{ name = "rabbit", base = 221023400, room = 221023401, entry = { -233, -1572 }, exit = { 196, 466 } },
	{ name = "roundTable", base = 240020500, room = 240020501, entry = { 765, 120 }, exit = { 839, -779 } },
	{ name = "remember", base = 240040511, room = 240040800, entry = { 1040, 1092 }, exit = { 1080, 1094 } },
	{ name = "protect", base = 240040520, room = 240040900, entry = { -1082, 1106 }, exit = { 305, 445 } },
	{ name = "sand", base = 260020600, room = 260020630, entry = { -180, -178 }, exit = { 742, 97 } },
	{ name = "error", base = 261020300, room = 261020301, entry = { 245, -93 }, exit = { -310, -85 } },
	{ name = "boat", base = 541000300, room = 541000301, rooms = 20, entry = { 273, 214 }, exit = { 294, -180 } },
	{ name = "longkiss", base = 541020610, room = 541020620, rooms = 20, entry = { 252, -560 }, exit = { 823, -592 } },
	{ name = "high", base = 551030000, room = 551030001, rooms = 20, entry = { -160, 638 }, exit = { -131, 158 } },
}

local function in_rooms(map_id, dungeon)
	return map_id >= dungeon.room and map_id < dungeon.room + (dungeon.rooms or ROOMS)
end

local function wait_room(ctx, bot, dungeon)
	if in_rooms(bot:map(), dungeon) then
		return bot:map()
	end
	local warp = bot:request(resp.warp, nil, function(p)
		return in_rooms(p.character.map, dungeon)
	end, 5000)
	if warp == false then
		return ctx:fail(string.format("%s %s 미니던전 방으로 이동하지 않음 (현재 %d)", bot:name(), dungeon.name, bot:map()))
	end
	return bot:map()
end

local function walk_in(ctx, bot, x, y)
	if pq.move(bot, x, y) == false then
		return ctx:fail(bot:name() .. " 포탈 앞으로 이동 실패")
	end
	return bot:warp("MD00")
end

local function explore(ctx, dungeon)
	local leader = ctx:bot(0)
	local member = ctx:bot(1)
	for i = 0, ctx:bot_count() - 1 do
		if ctx:bot(i):map_move(dungeon.base) == false then
			return ctx:fail(dungeon.name .. " 입구로 이동 실패")
		end
	end
	if walk_in(ctx, leader, dungeon.entry[1], dungeon.entry[2]) == false then
		return false
	end
	local room = wait_room(ctx, leader, dungeon)
	if room == false then
		return false
	end
	if wait_room(ctx, member, dungeon) ~= room then
		return ctx:fail(dungeon.name .. " 파티원이 다른 방으로 이동: " .. member:map())
	end
	for i = 0, ctx:bot_count() - 1 do
		local bot = ctx:bot(i)
		if pq.move(bot, dungeon.exit[1], dungeon.exit[2]) == false then
			return ctx:fail(bot:name() .. " 나가는 포탈 앞으로 이동 실패")
		end
		if pq.portal(ctx, bot, "out00", dungeon.base) == false then
			return false
		end
	end
	return true
end

local scenarios = {
	function(ctx)
		local member = ctx:bot(1)
		local dungeon = DUNGEONS[1]
		if member:map_move(dungeon.base) == false then
			return ctx:fail("돼지 미니던전 입구로 이동 실패")
		end
		if pq.move(member, dungeon.entry[1], dungeon.entry[2]) == false then
			return ctx:fail("포탈 앞으로 이동 실패")
		end
		local p = member:request(resp.notice, req.warp { target = 4294967295, portal_name = "MD00" }, function(p)
			return p.message:find("파티장이 아닙니다", 1, true) ~= nil
		end, 5000)
		if p == false then
			return ctx:fail("파티원 단독 입장이 거절되지 않음 (현재 " .. member:map() .. ")")
		end
		return true
	end,
}
for _, dungeon in ipairs(DUNGEONS) do
	scenarios[#scenarios + 1] = function(ctx)
		return explore(ctx, dungeon)
	end
end

test_suite {
	name = "Mini Dungeon: 미니던전 입장·퇴장",
	bot_count = 2,

	on_initialize = function(ctx)
		for i = 0, ctx:bot_count() - 1 do
			if pq.command(ctx:bot(i), "/봇초기화 120 100 0 - -", "봇초기화 완료") == false then
				return ctx:fail("봇 초기화 실패")
			end
			if pq.command(ctx:bot(i), "/무적", "무적 상태: enabled") == false then
				return ctx:fail("무적 설정 실패")
			end
		end
		return pq.form_party(ctx)
	end,

	scenarios = scenarios,
}
