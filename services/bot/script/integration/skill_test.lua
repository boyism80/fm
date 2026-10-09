local pq = require("script/integration/lib/party_quest")

local HUNTING_GROUND = 104040000
local NO_DOOR_MAP = 100000005
local STONE_GOLEM = 5130101
local MAGIC_ROCK = 4006000
local MYSTIC_DOOR = 2311002
local GM_HOLY_SYMBOL = 9001002
local PRIEST = 231
local GM = 900

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
	return true
end

local function learn(ctx, bot, skill_id)
	if pq.command(bot, "/스킬레벨 " .. skill_id .. " 1", "스킬레벨 설정") == false then
		return ctx:fail(bot:name() .. " 스킬 배우기 실패: " .. skill_id)
	end
	return true
end

local function cast(ctx, bot, skill_id)
	local p = bot:request(resp.update_stats, req.active_skill { skill_id = skill_id, skill_level = 1 }, function(p)
		return p.unlock_action
	end, 3000)
	if p == false then
		return ctx:fail(bot:name() .. " 스킬 사용 후 잠금이 풀리지 않음: " .. skill_id)
	end
	return true
end

local function kill_exp(ctx, bot)
	local before = {}
	for _, mob in ipairs(bot:mobs(STONE_GOLEM)) do
		before[mob.oid] = true
	end
	if bot:request(resp.spawn_mob, req.normal_chat { message = "/몬스터생성 " .. STONE_GOLEM }) == false then
		ctx:fail("스톤골렘 생성 응답 없음")
		return nil
	end
	local oid = nil
	for _, mob in ipairs(bot:mobs(STONE_GOLEM)) do
		if before[mob.oid] == nil then
			oid = mob.oid
		end
	end
	if oid == nil then
		ctx:fail("생성한 스톤골렘을 찾지 못함")
		return nil
	end

	local exp = bot:exp()
	if bot:kill(oid) == false then
		ctx:fail("스톤골렘 처치 실패")
		return nil
	end
	ctx:sleep(500)
	return bot:exp() - exp
end

local function mystic_door(ctx)
	local bot = ctx:bot(0)
	if setup(ctx, bot, 70, PRIEST, MAGIC_ROCK .. ":5") == false then
		return false
	end
	if learn(ctx, bot, MYSTIC_DOOR) == false then
		return false
	end
	if pq.command(bot, "/마력바꾸기 3000", "마력 설정") == false then
		return ctx:fail("마력 설정 실패")
	end

	if bot:map_move(NO_DOOR_MAP) == false then
		return ctx:fail("미스틱 도어 금지 맵 이동 실패")
	end
	if cast(ctx, bot, MYSTIC_DOOR) == false then
		return false
	end
	ctx:sleep(500)
	local rocks = bot:items()[MAGIC_ROCK] or 0
	if check(ctx, rocks == 5, "미스틱 도어 금지 맵에서 마법의 돌이 소모됨: " .. rocks) == false then
		return false
	end

	if bot:map_move(HUNTING_GROUND) == false then
		return ctx:fail("사냥터 이동 실패")
	end
	local p = bot:request(resp.spawn_door, req.active_skill { skill_id = MYSTIC_DOOR, skill_level = 1 }, function(p)
		return p.owner_id == bot:id()
	end, 5000)
	if p == false then
		return ctx:fail("사냥터에서 미스틱 도어가 생기지 않음")
	end
	rocks = bot:items()[MAGIC_ROCK] or 0
	return check(ctx, rocks == 4, "미스틱 도어 사용 후 마법의 돌 개수가 다름: " .. rocks)
end

local function gm_holy_symbol(ctx)
	local bot = ctx:bot(1)
	if setup(ctx, bot, 120, GM, "-") == false then
		return false
	end
	if learn(ctx, bot, GM_HOLY_SYMBOL) == false then
		return false
	end
	if bot:instance_move(HUNTING_GROUND) == false then
		return ctx:fail("사냥터 이동 실패")
	end

	local base = kill_exp(ctx, bot)
	if base == nil then
		return false
	end
	if cast(ctx, bot, GM_HOLY_SYMBOL) == false then
		return false
	end
	local boosted = kill_exp(ctx, bot)
	if boosted == nil then
		return false
	end
	return check(ctx, base > 0 and boosted * 100 >= base * 143 and boosted * 100 <= base * 147,
		string.format("혼자 받은 GM 홀리 심볼 경험치가 145%%가 아님: %d -> %d", base, boosted))
end

test_suite {
	name = "Skill: 미스틱 도어 맵 제한·GM 홀리 심볼",
	bot_count = 2,

	scenarios = {
		mystic_door,
		gm_holy_symbol,
	},
}
