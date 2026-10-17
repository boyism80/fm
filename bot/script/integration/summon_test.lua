local pq = require("script/integration/lib/party_quest")

local HENESYS = 100000000
local PARK = 100000200

local OUTLAW = 521
local DARK_KNIGHT = 132
local OCTOPUS = 5211001
local BEHOLDER = 1321007
local BEHOLDER_HEALING = 1320008
local HEAL = 40

local HP = 0x400
local HIT_ENV = -3

local function check(ctx, ok, message)
	if ok == false then
		return ctx:fail(message)
	end
	return true
end

local function setup(ctx, bot, class, skills)
	if pq.command(bot, "/봇초기화 120 " .. class .. " 0 -", "봇초기화 완료") == false then
		return ctx:fail(bot:name() .. " 봇 초기화 실패")
	end
	for _, skill_id in ipairs(skills) do
		if pq.command(bot, "/스킬레벨 " .. skill_id .. " 1", "스킬레벨 설정") == false then
			return ctx:fail(bot:name() .. " 스킬 배우기 실패: " .. skill_id)
		end
	end
	if pq.command(bot, "/마력바꾸기 5000", "마력 설정: 5000") == false then
		return ctx:fail(bot:name() .. " 마력 설정 실패")
	end
	if bot:map_move(HENESYS) == false then
		return ctx:fail(bot:name() .. " 맵 이동 실패")
	end
	return true
end

local function summon(ctx, bot, skill_id)
	if pq.command(bot, "/쿨타임초기화", "쿨타임 초기화 완료") == false then
		return ctx:fail(bot:name() .. " 쿨타임 초기화 실패")
	end
	local p = bot:request(resp.spawn_summon, req.active_skill { skill_id = skill_id, skill_level = 1 }, function(p)
		return p.owner_id == bot:id() and p.skill_id == skill_id
	end, 5000)
	if p == false then
		return ctx:fail(bot:name() .. " 소환수가 나타나지 않음: " .. skill_id)
	end
	return p
end

local function reappears(bot, skill_id)
	local moved = bot:request(resp.spawn_summon, req.normal_chat { message = "/맵이동 " .. PARK }, function(p)
		return p.owner_id == bot:id() and p.skill_id == skill_id
	end, 3000)
	return moved ~= false
end

-- A stationary summon stays behind and ends when its owner changes maps; a following summon comes along.
local function map_change(ctx)
	local outlaw, knight = ctx:bot(0), ctx:bot(1)
	if setup(ctx, outlaw, OUTLAW, { OCTOPUS }) == false or setup(ctx, knight, DARK_KNIGHT, { BEHOLDER }) == false then
		return false
	end
	if summon(ctx, outlaw, OCTOPUS) == false or summon(ctx, knight, BEHOLDER) == false then
		return false
	end

	if check(ctx, reappears(outlaw, OCTOPUS) == false, "문어가 맵 이동을 따라옴") == false then
		return false
	end
	if check(ctx, outlaw:map() == PARK, "문어 주인 맵 이동 실패") == false then
		return false
	end
	return check(ctx, reappears(knight, BEHOLDER), "비홀더가 맵 이동을 따라오지 않음")
end

local function heal(bot, oid)
	local p = bot:request(resp.update_stats, req.summon_skill { summon_oid = oid, sub_skill_id = BEHOLDER_HEALING }, function(p)
		return p.stats[HP] ~= nil
	end, 1500)
	if p == false then
		return nil
	end
	return p.stats[HP]
end

-- Beholder's healing works once per its interval.
local function healing_interval(ctx)
	local knight = ctx:bot(1)
	if setup(ctx, knight, DARK_KNIGHT, { BEHOLDER, BEHOLDER_HEALING }) == false then
		return false
	end
	local beholder = summon(ctx, knight, BEHOLDER)
	if beholder == false then
		return false
	end
	if pq.command(knight, "/체력바꾸기 1000", "체력 설정: 1000") == false then
		return ctx:fail("체력 설정 실패")
	end
	local hit = knight:request(resp.update_stats, req.damaged { type = HIT_ENV, damage = 500 }, function(p)
		return p.stats[HP] ~= nil
	end, 3000)
	if hit == false then
		return ctx:fail("피해를 받지 않음")
	end
	local hp = hit.stats[HP]

	local healed = heal(knight, beholder.oid)
	if check(ctx, healed == hp + HEAL, string.format("첫 힐링 뒤 체력 %s (기대 %d)", tostring(healed), hp + HEAL)) == false then
		return false
	end
	return check(ctx, heal(knight, beholder.oid) == nil, "주기 안에 힐링이 다시 됨")
end

test_suite {
	name = "Summon: 고정형 소환수 맵 이동·비홀더스 힐링 주기",
	bot_count = 2,

	scenarios = {
		map_change,
		healing_interval,
	},
}
