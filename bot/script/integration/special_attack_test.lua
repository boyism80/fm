local pq = require("script/integration/lib/party_quest")

-- Hypnotize only checks the first mob in range, so the map must have no mobs of its own.
local HENESYS = 100000000
local SNAIL = 100100
-- Self-destructs, and its attack reaches 110px for at most 145 damage.
local BOMB = 5100002
local STONE_GOLEM = 5130101
local ZOMBIE_MUSHROOM = 2230101
local CLERIC = 230
local MARAUDER = 511
local CORSAIR = 522
local ENERGY_CHARGE = 5110001
local HYPNOTIZE = 5221009
local HEAL = 2301002
local QUIET_MS = 1500
local PRIEST = 231
local MYSTIC_DOOR = 2311002
local MAGIC_ROCK = 4006000
local DOOR_FIELD = 100010000

local function setup(ctx, bot, level, class)
    if pq.command(bot, "/봇초기화 " .. level .. " " .. class .. " 0 -", "봇초기화 완료") == false then
        return ctx:fail(bot:name() .. " 봇 초기화 실패")
    end
    if bot:instance_move(HENESYS) == false then
        return ctx:fail(bot:name() .. " 헤네시스 이동 실패")
    end
    return true
end

local function spawn(ctx, bot, mob_id, count)
    local before = {}
    for _, mob in ipairs(bot:mobs(mob_id)) do
        before[mob.oid] = true
    end
    if bot:request(resp.spawn_mob, req.normal_chat { message = "/몬스터생성 " .. mob_id .. " " .. count }) == false then
        ctx:fail("몬스터 생성 실패: " .. mob_id)
        return nil
    end
    local oids = {}
    for _ = 1, 20 do
        oids = {}
        for _, mob in ipairs(bot:mobs(mob_id)) do
            if before[mob.oid] == nil then
                table.insert(oids, mob.oid)
            end
        end
        if #oids >= count then
            return oids
        end
        ctx:sleep(100)
    end
    ctx:fail(string.format("생성한 몬스터를 찾지 못함: %d (%d/%d)", mob_id, #oids, count))
    return nil
end

local function died(bot, oid, pkt, timeout)
    return bot:request(resp.die_mob, pkt, function(p)
        return p.oid == oid
    end, timeout)
end

local function monster_bomb(ctx)
    local bot = ctx:bot(0)
    if setup(ctx, bot, 30, 0) == false then
        return false
    end

    local snail = spawn(ctx, bot, SNAIL, 1)
    if snail == nil then
        return false
    end
    if died(bot, snail[1], req.monster_bomb { oid = snail[1] }, QUIET_MS) ~= false then
        return ctx:fail("자폭하지 않는 몬스터가 자폭함")
    end

    local bomb = spawn(ctx, bot, BOMB, 1)
    if bomb == nil then
        return false
    end
    if died(bot, bomb[1], req.monster_bomb { oid = bomb[1] }, 5000) == false then
        return ctx:fail("자폭 몬스터가 자폭하지 않음")
    end
    return true
end

local function hypnotize_damage(ctx)
    local bot = ctx:bot(0)
    if setup(ctx, bot, 120, CORSAIR) == false then
        return false
    end
    if pq.command(bot, "/스킬레벨 " .. HYPNOTIZE .. " 20", "스킬레벨 설정") == false then
        return ctx:fail("하이퍼노타이즈 배우기 실패")
    end

    local mobs = spawn(ctx, bot, BOMB, 2)
    if mobs == nil then
        return false
    end
    local function hit(source, target, damage, timeout)
        return bot:request(resp.show_mob_hp, req.hypnotize_damage {
            source_oid = source,
            character_id = bot:id(),
            target_oid = target,
            damage = damage,
        }, function(p)
            return p.oid == target
        end, timeout)
    end

    if hit(mobs[1], mobs[2], 100, QUIET_MS) ~= false then
        return ctx:fail("최면에 걸리지 않은 몬스터가 다른 몬스터를 때림")
    end

    local source, target
    for _ = 1, 10 do
        if pq.command(bot, "/마력바꾸기 30000", "마력 설정") == false then
            return ctx:fail("마력 설정 실패")
        end
        if bot:request(resp.update_stats, req.active_skill { skill_id = HYPNOTIZE, skill_level = 20 }, function(p)
                return p.unlock_action
            end, 3000) == false then
            return ctx:fail("하이퍼노타이즈 사용 후 행동이 풀리지 않음")
        end
        if hit(mobs[1], mobs[2], 100, QUIET_MS) ~= false then
            source, target = mobs[1], mobs[2]
            break
        end
        if hit(mobs[2], mobs[1], 100, QUIET_MS) ~= false then
            source, target = mobs[2], mobs[1]
            break
        end
        if pq.command(bot, "/쿨타임초기화", "쿨타임 초기화 완료") == false then
            return ctx:fail("쿨타임 초기화 실패")
        end
    end
    if source == nil then
        return ctx:fail("하이퍼노타이즈 후 최면 피해가 들어가지 않음")
    end

    if hit(source, target, 146, QUIET_MS) ~= false then
        return ctx:fail("몬스터 공격력을 넘는 최면 피해가 들어감")
    end
    if bot:request(resp.show_mob_hp, req.hypnotize_damage {
            source_oid = source,
            character_id = bot:id() + 1,
            target_oid = target,
            damage = 100,
        }, nil, QUIET_MS) ~= false then
        return ctx:fail("다른 캐릭터 이름으로 보낸 최면 피해가 들어감")
    end
    return true
end

local function heal_attack(ctx)
    local bot = ctx:bot(0)
    if setup(ctx, bot, 50, CLERIC) == false then
        return false
    end
    if pq.command(bot, "/스킬레벨 " .. HEAL .. " 30", "스킬레벨 설정") == false then
        return ctx:fail("힐 배우기 실패")
    end
    if pq.command(bot, "/전능 999", "전능 설정") == false then
        return ctx:fail("스탯 설정 실패")
    end
    if pq.command(bot, "/마력바꾸기 30000", "마력 설정") == false then
        return ctx:fail("마력 설정 실패")
    end

    local snail = spawn(ctx, bot, SNAIL, 1)
    local zombie = spawn(ctx, bot, ZOMBIE_MUSHROOM, 1)
    if snail == nil or zombie == nil then
        return false
    end
    local function killed(oid)
        bot:magic_attack(oid, 0, 1, HEAL)
        return died(bot, oid, nil, QUIET_MS) ~= false
    end

    if killed(snail[1]) then
        return ctx:fail("힐이 언데드가 아닌 몬스터를 때림")
    end
    if killed(zombie[1]) == false then
        return ctx:fail("힐이 언데드 몬스터를 때리지 못함")
    end
    return true
end

local function energy_attack(ctx)
    local bot = ctx:bot(0)
    local observer = ctx:bot(1)
    if setup(ctx, bot, 70, MARAUDER) == false then
        return false
    end
    if observer:instance_move(HENESYS) == false then
        return ctx:fail("관찰 봇 헤네시스 이동 실패")
    end
    if pq.command(bot, "/스킬레벨 " .. ENERGY_CHARGE .. " 40", "스킬레벨 설정") == false then
        return ctx:fail("에너지 차지 배우기 실패")
    end
    -- Starting stats leave the golem's defense above any hit, and a rejected hit charges nothing.
    if pq.command(bot, "/전능 999", "전능 설정") == false then
        return ctx:fail("스탯 설정 실패")
    end

    local golem = spawn(ctx, bot, STONE_GOLEM, 1)
    if golem == nil then
        return false
    end
    local function released()
        bot:energy_attack(golem[1], 1, 1)
        return bot:request_on(observer, resp.energy_attack, nil, function(p)
            return p.attack.character_id == bot:id()
        end, QUIET_MS) ~= false
    end

    if released() then
        return ctx:fail("에너지를 모으기 전에 에너지 어택이 나감")
    end
    for _ = 1, 12 do
        for _ = 1, 10 do
            bot:attack(golem[1], 1, 1)
            ctx:sleep(50)
        end
        if released() then
            return true
        end
    end
    return ctx:fail("에너지를 다 모았는데 에너지 어택이 나가지 않음")
end

local function door_attack(ctx)
    local bot = ctx:bot(0)
    if pq.command(bot, "/봇초기화 70 " .. PRIEST .. " 0 " .. MAGIC_ROCK .. ":2 -", "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    if bot:instance_move(DOOR_FIELD) == false then
        return ctx:fail("미스틱 도어 시험 맵 이동 실패")
    end
    if pq.command(bot, "/스킬레벨 " .. MYSTIC_DOOR .. " 20", "스킬레벨 설정") == false then
        return ctx:fail("미스틱 도어 배우기 실패")
    end
    if pq.command(bot, "/마력바꾸기 30000", "마력 설정") == false then
        return ctx:fail("마력 설정 실패")
    end
    local snail = spawn(ctx, bot, SNAIL, 1)
    if snail == nil then
        return false
    end

    bot:magic_attack(snail[1], 0, 1, MYSTIC_DOOR)
    if bot:request(resp.spawn_door, nil, nil, QUIET_MS) ~= false then
        return ctx:fail("공격 패킷으로 미스틱 도어가 열림")
    end
    if (bot:items()[MAGIC_ROCK] or 0) ~= 2 then
        return ctx:fail("공격 패킷 미스틱 도어가 마법의 돌을 씀")
    end

    if bot:request(resp.spawn_door, req.active_skill { skill_id = MYSTIC_DOOR, skill_level = 20 }, nil, 3000) == false then
        return ctx:fail("정상 사용한 미스틱 도어가 열리지 않음")
    end
    if (bot:items()[MAGIC_ROCK] or 0) ~= 1 then
        return ctx:fail("정상 사용한 미스틱 도어가 마법의 돌을 쓰지 않음")
    end
    return true
end

test_suite {
    name = "SpecialAttack: 몬스터 자폭·최면 피해·힐 공격·에너지 어택·공격 패킷 비공격 스킬",
    bot_count = 2,

    scenarios = {
        monster_bomb,
        hypnotize_damage,
        heal_attack,
        energy_attack,
        door_attack,
    },
}
