local pq = require("script/integration/lib/party_quest")

local HENESYS = 100000000
local PARK = 100000200

local OUTLAW = 521
local DARK_KNIGHT = 132
local OCTOPUS = 5211001
local GAVIOTA = 5211002
local BEHOLDER = 1321007
local BEHOLDER_HEALING = 1320008
local HEAL = 40
local SNAIL = 100100
local FP_ARCH_MAGE = 212
local IL_ARCH_MAGE = 222
local ELQUINES = 2121005
local IFRIT = 2221005
local MOB_STUN = 0x80
local MOB_FREEZE = 0x100

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
    if bot:instance_move(HENESYS) == false then
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

local function reappears(ctx, bot, skill_id)
    local moved = bot:request(resp.spawn_summon,
        req.normal_chat { message = string.format("/인스턴스이동 %d %d", PARK, ctx:seat()) }, function(p)
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

    if check(ctx, reappears(ctx, outlaw, OCTOPUS) == false, "문어가 맵 이동을 따라옴") == false then
        return false
    end
    if check(ctx, outlaw:map() == PARK, "문어 주인 맵 이동 실패") == false then
        return false
    end
    return check(ctx, reappears(ctx, knight, BEHOLDER), "비홀더가 맵 이동을 따라오지 않음")
end

local function heal(bot, oid)
    local p = bot:request(resp.update_stats, req.summon_skill { summon_oid = oid, sub_skill_id = BEHOLDER_HEALING },
        function(p)
            return p.stats[HP] ~= nil
        end, 1500)
    if p == false then
        return nil
    end
    return p.stats[HP]
end

local function hurt(ctx, bot)
    if pq.command(bot, "/체력바꾸기 1000", "체력 설정: 1000") == false then
        return ctx:fail(bot:name() .. " 체력 설정 실패")
    end
    local hit = bot:request(resp.update_stats, req.damaged { type = HIT_ENV, damage = 500 }, function(p)
        return p.stats[HP] ~= nil
    end, 3000)
    if hit == false then
        return ctx:fail(bot:name() .. " 피해를 받지 않음")
    end
    return hit.stats[HP]
end

-- Beholder's healing works once per its interval, and only through the beholder.
local function healing_interval(ctx)
    local knight, outlaw = ctx:bot(1), ctx:bot(0)
    if setup(ctx, knight, DARK_KNIGHT, { BEHOLDER, BEHOLDER_HEALING }) == false or setup(ctx, outlaw, OUTLAW, { OCTOPUS, BEHOLDER_HEALING }) == false then
        return false
    end
    local beholder = summon(ctx, knight, BEHOLDER)
    if beholder == false then
        return false
    end
    local octopus = summon(ctx, outlaw, OCTOPUS)
    if octopus == false then
        return false
    end
    if hurt(ctx, outlaw) == false then
        return false
    end
    if check(ctx, heal(outlaw, octopus.oid) == nil, "비홀더가 아닌 소환수로 힐링이 됨") == false then
        return false
    end

    local hp = hurt(ctx, knight)
    if hp == false then
        return false
    end
    local healed = heal(knight, beholder.oid)
    if check(ctx, healed == hp + HEAL, string.format("첫 힐링 뒤 체력 %s (기대 %d)", tostring(healed), hp + HEAL)) == false then
        return false
    end
    return check(ctx, heal(knight, beholder.oid) == nil, "주기 안에 힐링이 다시 됨")
end

-- Gaviota leaves after its one attack.
local function gaviota_once(ctx)
    local outlaw = ctx:bot(0)
    if setup(ctx, outlaw, OUTLAW, { GAVIOTA }) == false then
        return false
    end
    local gaviota = summon(ctx, outlaw, GAVIOTA)
    if gaviota == false then
        return false
    end
    local removed = outlaw:request(resp.remove_summon, req.summon_attack { summon_oid = gaviota.oid, damages = {} },
        function(p)
            return p.oid == gaviota.oid
        end, 3000)
    return check(ctx, removed ~= false, "가비오타가 공격한 뒤에도 남음")
end

-- Another player sees the summon's attack under the summon's object id.
local function attack_seen(ctx)
    local outlaw, knight = ctx:bot(0), ctx:bot(1)
    if setup(ctx, outlaw, OUTLAW, { OCTOPUS }) == false or setup(ctx, knight, DARK_KNIGHT, {}) == false then
        return false
    end
    local octopus = summon(ctx, outlaw, OCTOPUS)
    if octopus == false then
        return false
    end
    local spawned = outlaw:request(resp.spawn_mob, req.normal_chat { message = "/몬스터생성 " .. SNAIL }, function(p)
        return p.mob.mob_id == SNAIL
    end, 5000)
    if spawned == false then
        return ctx:fail("달팽이 생성 응답 없음")
    end

    local seen = outlaw:request_on(knight, resp.summon_attack, req.summon_attack {
        summon_oid = octopus.oid,
        damages = { { oid = spawned.mob.oid, damage_pairs = { { damage = 1 } } } },
    }, function(p)
        return p.character_id == outlaw:id()
    end, 3000)
    if seen == false then
        return ctx:fail("다른 플레이어에게 소환수 공격이 보이지 않음")
    end
    return check(ctx, seen.summon_oid == octopus.oid,
        string.format("소환수 공격의 소환수 번호 %d (기대 OID %d)", seen.summon_oid, octopus.oid))
end

-- A summon whose WZ has no prop puts its status on every mob it damages.
local function status_hit(ctx, bot, class, skill_id, flag, name)
    if setup(ctx, bot, class, { skill_id }) == false then
        return false
    end
    local summoned = summon(ctx, bot, skill_id)
    if summoned == false then
        return false
    end
    local spawned = bot:request(resp.spawn_mob, req.normal_chat { message = "/몬스터생성 " .. SNAIL }, function(p)
        return p.mob.mob_id == SNAIL
    end, 5000)
    if spawned == false then
        return ctx:fail("달팽이 생성 응답 없음")
    end
    local oid = spawned.mob.oid
    local applied = bot:request(resp.apply_mob_buff, req.summon_attack {
        summon_oid = summoned.oid,
        damages = { { oid = oid, damage_pairs = { { damage = 1 } } } },
    }, function(p)
        return p.oid == oid and math.floor(p.status / flag) % 2 == 1
    end, 3000)
    return check(ctx, applied ~= false, name .. " 공격에 맞은 달팽이에 상태이상이 걸리지 않음")
end

local function elquines_stun(ctx)
    return status_hit(ctx, ctx:bot(0), FP_ARCH_MAGE, ELQUINES, MOB_STUN, "엘퀴네스")
end

local function ifrit_freeze(ctx)
    return status_hit(ctx, ctx:bot(1), IL_ARCH_MAGE, IFRIT, MOB_FREEZE, "이프리트")
end

test_suite {
    name = "Summon: 고정형 소환수 맵 이동·비홀더스 힐링·가비오타 한 번 공격·공격 표시·엘퀴네스 기절·이프리트 빙결",
    bot_count = 2,

    scenarios = {
        map_change,
        healing_interval,
        gaviota_once,
        attack_seen,
        elquines_stun,
        ifrit_freeze,
    },
}
