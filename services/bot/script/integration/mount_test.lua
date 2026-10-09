local pq = require("script/integration/lib/party_quest")

local HENESYS_PARK = 100000200
local NO_MOUNT_MAP = 101000100
local TAMING_MOB = 1902000
local SADDLE = 1912000
local FOOD = 2260000
local MONSTER_RIDER = 1004
local BATTLESHIP = 5221006
local BATTLESHIP_VEHICLE = 1932000
local BATTLESHIP_GAUGE = 5221999
local EQUIP = 1
local TAMING_MOB_SLOT = -18
local SADDLE_SLOT = -19
local HIT_ENV = -3

local function check(ctx, ok, message)
	if ok == false then
		return ctx:fail(message)
	end
	return true
end

local function setup(ctx, bot, level, class, items)
	if pq.command(bot, "/봇초기화 " .. level .. " " .. class .. " 0 " .. items, "봇초기화 완료") == false then
		return ctx:fail(bot:name() .. " 봇 초기화 실패")
	end
	if bot:map_move(HENESYS_PARK) == false then
		return ctx:fail(bot:name() .. " 맵 이동 실패")
	end
	return true
end

local function mount(ctx, bot)
	local p = pq.command(bot, "/탈것정보", "탈것 레벨")
	if p == false then
		ctx:fail(bot:name() .. " 탈것 정보 없음")
		return nil
	end
	local level, exp, fatigue, hp, vehicle = p.message:match("탈것 레벨 (%d+) 경험치 (%d+) 피로도 (%d+) HP (%d+) 탑승 (%d+)")
	return {
		level = tonumber(level),
		exp = tonumber(exp),
		fatigue = tonumber(fatigue),
		hp = tonumber(hp),
		vehicle = tonumber(vehicle),
	}
end

local function riding(ctx, bot, vehicle, what)
	local m = mount(ctx, bot)
	if m == nil then
		return false
	end
	return check(ctx, m.vehicle == vehicle, bot:name() .. " " .. what .. ": 탑승 " .. tostring(m.vehicle))
end

local function unlocked(ctx, bot, pkt, what)
	local p = bot:request(resp.update_stats, pkt, function(p)
		return p.unlock_action
	end, 3000)
	if p == false then
		return ctx:fail(bot:name() .. " 잠금이 풀리지 않음: " .. what)
	end
	return true
end

local function learn(ctx, bot, skill_id)
	if pq.command(bot, "/스킬레벨 " .. skill_id .. " 1", "스킬레벨 설정") == false then
		return ctx:fail(bot:name() .. " 스킬 배우기 실패: " .. skill_id)
	end
	return true
end

local function ride(ctx, bot, skill_id)
	return unlocked(ctx, bot, req.active_skill { skill_id = skill_id, skill_level = 1 }, "탑승 " .. skill_id)
end

local function wear(ctx, bot, item_id, slot)
	if bot:request(resp.inventory_operation, req.move_item { inventory_type = EQUIP, source = bot:slot(item_id) or 0, dest = slot, count = 1 }, nil, 3000) == false then
		return ctx:fail(bot:name() .. " 장착 실패: " .. item_id)
	end
	return true
end

local function rider(ctx, bot)
	if setup(ctx, bot, 70, 0, TAMING_MOB .. ":1," .. SADDLE .. ":1," .. FOOD .. ":5") == false then
		return false
	end
	if learn(ctx, bot, MONSTER_RIDER) == false then
		return false
	end
	if wear(ctx, bot, TAMING_MOB, TAMING_MOB_SLOT) == false then
		return false
	end
	return wear(ctx, bot, SADDLE, SADDLE_SLOT)
end

local function ride_without_equipment(ctx)
	local bot = ctx:bot(0)
	if setup(ctx, bot, 70, 0, TAMING_MOB .. ":1") == false then
		return false
	end
	if learn(ctx, bot, MONSTER_RIDER) == false then
		return false
	end
	if ride(ctx, bot, MONSTER_RIDER) == false then
		return false
	end
	if riding(ctx, bot, 0, "장비 없이 탑승") == false then
		return false
	end

	if wear(ctx, bot, TAMING_MOB, TAMING_MOB_SLOT) == false then
		return false
	end
	if ride(ctx, bot, MONSTER_RIDER) == false then
		return false
	end
	return riding(ctx, bot, 0, "안장 없이 탑승")
end

local function ride_and_dismount(ctx)
	local bot = ctx:bot(1)
	local watcher = ctx:bot(2)
	if setup(ctx, watcher, 30, 0, "-") == false then
		return false
	end
	if rider(ctx, bot) == false then
		return false
	end

	local p = bot:request_on(watcher, resp.update_ridding, req.active_skill { skill_id = MONSTER_RIDER, skill_level = 1 }, function(p)
		return p.character_id == bot:id()
	end, 3000)
	if p == false then
		return ctx:fail("주변 유저가 탑승을 보지 못함")
	end
	if check(ctx, p.mount_id == TAMING_MOB and p.skill_id == MONSTER_RIDER, "주변 탑승 패킷 값이 다름: " .. p.mount_id .. " " .. p.skill_id) == false then
		return false
	end
	if riding(ctx, bot, TAMING_MOB, "탑승") == false then
		return false
	end

	bot:send(req.cancel_buff { source_id = MONSTER_RIDER })
	if riding(ctx, bot, 0, "버프 취소 후") == false then
		return false
	end

	if ride(ctx, bot, MONSTER_RIDER) == false then
		return false
	end
	if riding(ctx, bot, TAMING_MOB, "다시 탑승") == false then
		return false
	end
	if bot:request(resp.inventory_operation, req.move_item { inventory_type = EQUIP, source = TAMING_MOB_SLOT, dest = 1, count = 1 }, nil, 3000) == false then
		return ctx:fail("길들인 몬스터 해제 실패")
	end
	return riding(ctx, bot, 0, "길들인 몬스터 해제 후")
end

local function feed(ctx)
	local bot = ctx:bot(3)
	if rider(ctx, bot) == false then
		return false
	end

	if unlocked(ctx, bot, req.use_mount_food { slot = bot:slot(FOOD) or 0, item_id = FOOD }, "탑승 전 먹이") == false then
		return false
	end
	ctx:sleep(300)
	if check(ctx, (bot:items()[FOOD] or 0) == 5, "탑승하지 않았는데 먹이가 사라짐") == false then
		return false
	end

	if ride(ctx, bot, MONSTER_RIDER) == false then
		return false
	end
	local p = bot:request(resp.update_mount, req.normal_chat { message = "/탈것피로도 50" }, nil, 3000)
	if p == false then
		return ctx:fail("피로도 설정 후 탈것 정보 없음")
	end
	if check(ctx, p.fatigue == 50, "피로도 설정이 다름: " .. p.fatigue) == false then
		return false
	end

	p = bot:request(resp.update_mount, req.use_mount_food { slot = bot:slot(FOOD) or 0, item_id = FOOD }, nil, 3000)
	if p == false then
		return ctx:fail("먹이 후 탈것 정보 없음")
	end
	if check(ctx, p.fatigue == 20, "먹이 후 피로도가 다름: " .. p.fatigue) == false then
		return false
	end
	if check(ctx, p.level == 2 and p.level_up and p.exp >= 15 and p.exp <= 24, string.format("먹이 경험치/레벨업이 다름: 레벨 %d 경험치 %d", p.level, p.exp)) == false then
		return false
	end

	p = bot:request(resp.update_mount, req.use_mount_food { slot = bot:slot(FOOD) or 0, item_id = FOOD }, nil, 3000)
	if p == false then
		return ctx:fail("두 번째 먹이 후 탈것 정보 없음")
	end
	if check(ctx, p.fatigue == 0, "피로도가 0 아래로 내려감: " .. p.fatigue) == false then
		return false
	end
	local exp = p.exp

	p = bot:request(resp.update_mount, req.use_mount_food { slot = bot:slot(FOOD) or 0, item_id = FOOD }, nil, 3000)
	if p == false then
		return ctx:fail("피로도 0 먹이 후 탈것 정보 없음")
	end
	if check(ctx, p.exp == exp and p.level_up == false, "피로도 0에서 경험치를 얻음") == false then
		return false
	end
	ctx:sleep(300)
	return check(ctx, (bot:items()[FOOD] or 0) == 2, "먹이 개수가 다름: " .. tostring(bot:items()[FOOD]))
end

local function tired(ctx)
	local bot = ctx:bot(4)
	if rider(ctx, bot) == false then
		return false
	end
	if ride(ctx, bot, MONSTER_RIDER) == false then
		return false
	end
	if pq.command(bot, "/탈것피로도 99", "탈것 피로도 99") == false then
		return ctx:fail("피로도 설정 실패")
	end
	if pq.command(bot, "/탈것피로", "탈것 피로도 100") == false then
		return ctx:fail("피로도가 한계에 이르지 않음")
	end
	if riding(ctx, bot, 0, "피로도 한계") == false then
		return false
	end

	if ride(ctx, bot, MONSTER_RIDER) == false then
		return false
	end
	return riding(ctx, bot, 0, "피로도 100에서 탑승")
end

local function persistence(ctx)
	local bot = ctx:bot(5)
	if setup(ctx, bot, 30, 0, "-") == false then
		return false
	end
	if pq.command(bot, "/탈것레벨 5", "탈것 레벨 5 경험치 105") == false then
		return ctx:fail("탈것 레벨 설정 실패")
	end
	if pq.command(bot, "/탈것피로도 30", "탈것 피로도 30") == false then
		return ctx:fail("탈것 피로도 설정 실패")
	end
	if bot:transfer(1) == false then
		return ctx:fail("채널 이동 실패")
	end
	local m = mount(ctx, bot)
	if m == nil then
		return false
	end
	return check(ctx, m.level == 5 and m.exp == 105 and m.fatigue == 30, string.format("채널 이동 후 탈것 상태가 다름: %d %d %d", m.level, m.exp, m.fatigue))
end

local function no_mount_map(ctx)
	local bot = ctx:bot(6)
	if rider(ctx, bot) == false then
		return false
	end
	if ride(ctx, bot, MONSTER_RIDER) == false then
		return false
	end
	if riding(ctx, bot, TAMING_MOB, "탑승") == false then
		return false
	end
	if bot:map_move(NO_MOUNT_MAP) == false then
		return ctx:fail("탈것 금지 맵 이동 실패")
	end
	return riding(ctx, bot, 0, "탈것 금지 맵 진입")
end

local function battleship(ctx)
	local bot = ctx:bot(7)
	if setup(ctx, bot, 130, 522, "-") == false then
		return false
	end
	if learn(ctx, bot, BATTLESHIP) == false then
		return false
	end
	if pq.command(bot, "/체력바꾸기 30000", "체력 설정") == false then
		return ctx:fail("체력 설정 실패")
	end
	if pq.command(bot, "/마력바꾸기 30000", "마력 설정") == false then
		return ctx:fail("마력 설정 실패")
	end

	if ride(ctx, bot, BATTLESHIP) == false then
		return false
	end
	local m = mount(ctx, bot)
	if m == nil then
		return false
	end
	if check(ctx, m.vehicle == BATTLESHIP_VEHICLE and m.hp == 2400, string.format("배틀쉽 탑승 상태가 다름: 탑승 %d HP %d", m.vehicle, m.hp)) == false then
		return false
	end

	bot:send(req.cancel_buff { source_id = BATTLESHIP })
	if riding(ctx, bot, 0, "배틀쉽 버프 취소 후") == false then
		return false
	end
	if ride(ctx, bot, BATTLESHIP) == false then
		return false
	end
	if riding(ctx, bot, BATTLESHIP_VEHICLE, "배틀쉽 다시 탑승") == false then
		return false
	end

	local p = bot:request(resp.skill_cooldown, req.damaged { type = HIT_ENV, damage = 100 }, function(p)
		return p.skill_id == BATTLESHIP_GAUGE
	end, 3000)
	if p == false then
		return ctx:fail("배틀쉽 게이지 갱신 없음")
	end
	if check(ctx, p.remaining_sec == 2300, "배틀쉽 HP가 다름: " .. p.remaining_sec) == false then
		return false
	end

	p = bot:request(resp.skill_cooldown, req.damaged { type = HIT_ENV, damage = 3000 }, function(p)
		return p.skill_id == BATTLESHIP
	end, 3000)
	if p == false then
		return ctx:fail("배틀쉽 파괴 후 쿨다운 없음")
	end
	if check(ctx, p.remaining_sec == 90, "배틀쉽 쿨다운이 다름: " .. p.remaining_sec) == false then
		return false
	end
	if riding(ctx, bot, 0, "배틀쉽 파괴") == false then
		return false
	end

	if ride(ctx, bot, BATTLESHIP) == false then
		return false
	end
	return riding(ctx, bot, 0, "쿨다운 중 배틀쉽 탑승")
end

test_suite {
	name = "Mount: 탑승·하차·피로도·먹이·배틀쉽",
	bot_count = 8,

	scenarios = {
		ride_without_equipment,
		ride_and_dismount,
		feed,
		tired,
		persistence,
		no_mount_map,
		battleship,
	},
}
