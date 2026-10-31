local pq = require("script/integration/lib/party_quest")

local HUNTING_GROUND = 104040000
-- The GM holy symbol reaches every player on the map, so the GM measures it alone, away from the other bots.
local GM_GROUND = 104040001
local NO_DOOR_MAP = 100000005
local STONE_GOLEM = 5130101
local MAGIC_ROCK = 4006000
local MYSTIC_DOOR = 2311002
local GM_HOLY_SYMBOL = 9001002
local HOLY_SYMBOL = 2311003
local SHADOW_WEB = 4111003
local SNAIL = 100100
local CLAW = 1472000
local STAR = 2070000
local LUCKY_SEVEN = 4001344
local POISON_MIST = 2111003
local MOB_POISON = 0x200
local POISON_BOMB = 14111006
-- Does not move, so it stays in the small poison bomb mist.
local POISON_FLOWER = 9300174
-- Takes mob hits (damagedByMob), so player skills cannot stun, slow, or poison it.
local MOON_BUNNY = 9300061
local SOUL_CHARGE = 11111007
local DAWN_WARRIOR3 = 1111
local SWORD = 1302000
-- Immune to lightning, neutral to holy, and not a boss.
local LIGHTNING_IMMUNE = 8148001
local CRUSADER = 111
local COMBO_ATTACK = 1111002
local COMA = 1111005
local QUIET_MS = 1500
local DRAGON_KNIGHT = 131
local SPEAR = 1432000
local SACRIFICE = 1311005
local SACRIFICE_HP_PERCENT = 5
local MUSHMOM = 6130101
local PRIEST = 231
local ROGUE = 400
local HERMIT = 411
local NIGHT_WALKER = 1411
local FP_MAGE = 211
local GM = 900
local EQUIP = 1
local WEAPON_SLOT = -11
local RED_POTION = 2000000
local ALCHEMIST = 4110000
local DARK_SIGHT = 4001003
local DARK_SIGHT_SPEED = -30
local ALCHEMIST_CYGNUS = 14110003
local IMPROVING_MP_RECOVERY = 2000000
local MP_BURN_MOB = 8110300
local HIT_ENV = -3
local NIGHT_LORD = 412
local PAGE = 120
local SHOWDOWN = 4121003
local THREATEN = 1201006
local MAGICIAN = 200
local MAGIC_GUARD = 2001002
local MAGIC_GUARD_MP = 6
local MOB_WATK = 0x1
local MOB_MDEF = 0x8
local MOB_SHOWDOWN = 0x800
local HP = 0x400
local SAUNA = 105040401
local SAUNA_ROBES = { 1050018, 1051017 }
local OVERALL_SLOT = -5
local HP_HEAL_INTERVAL_MS = 10000
local MP = 0x1000
local DRAGON_ROAR = 1311006
local DRAGON_ROAR_HP_PERCENT = 59
local DRAGON_ROAR_STUN_MS = 4000
local GM_HYPER_BODY = 9001008
local MARAUDER = 511
local TRANSFORM = 5111005
local BUFF = {
    weapon_def = { mask = 0x2, position = 4 },
    magic_def = { mask = 0x8, position = 4 },
    speed = { mask = 0x80, position = 4 },
    jump = { mask = 0x100, position = 4 },
    max_hp = { mask = 0x2000, position = 4 },
    morph = { mask = 0x2, position = 3 },
    darksight = { mask = 0x400, position = 4 },
}

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

    if bot:instance_move(NO_DOOR_MAP) == false then
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

    if bot:instance_move(HUNTING_GROUND) == false then
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

local function dead_mystic_door(ctx)
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
    if bot:instance_move(HUNTING_GROUND) == false then
        return ctx:fail("사냥터 이동 실패")
    end
    if pq.command(bot, "/체력바꾸기 0", "체력 설정") == false then
        return ctx:fail("체력 설정 실패")
    end

    local p = bot:request(resp.spawn_door, req.active_skill { skill_id = MYSTIC_DOOR, skill_level = 1 }, function(p)
        return p.owner_id == bot:id()
    end, 2000)
    if p ~= false then
        return ctx:fail("죽은 캐릭터가 미스틱 도어를 엶")
    end
    local rocks = bot:items()[MAGIC_ROCK] or 0
    if check(ctx, rocks == 5, "죽은 캐릭터의 마법의 돌이 소모됨: " .. rocks) == false then
        return false
    end
    if pq.command(bot, "/체력바꾸기 3000", "체력 설정") == false then
        return ctx:fail("체력 복구 실패")
    end
    return true
end

local function gm_holy_symbol(ctx)
    local bot = ctx:bot(1)
    if setup(ctx, bot, 120, GM, "-") == false then
        return false
    end
    if learn(ctx, bot, GM_HOLY_SYMBOL) == false then
        return false
    end
    if bot:instance_move(GM_GROUND) == false then
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

local function solo_holy_symbol(ctx)
    local bot = ctx:bot(0)
    if setup(ctx, bot, 70, PRIEST, "-") == false then
        return false
    end
    if learn(ctx, bot, HOLY_SYMBOL) == false then
        return false
    end
    if pq.command(bot, "/마력바꾸기 3000", "마력 설정") == false then
        return ctx:fail("마력 설정 실패")
    end
    if bot:instance_move(HUNTING_GROUND) == false then
        return ctx:fail("사냥터 이동 실패")
    end

    local base = kill_exp(ctx, bot)
    if base == nil then
        return false
    end
    if cast(ctx, bot, HOLY_SYMBOL) == false then
        return false
    end
    local boosted = kill_exp(ctx, bot)
    if boosted == nil then
        return false
    end
    return check(ctx, base > 0 and boosted > base, string.format("파티 없이 쓴 홀리 심볼이 시전자에게 걸리지 않음: %d -> %d", base, boosted))
end

local function spawn_golem(ctx, bot, x, y)
    if pq.move(bot, x, y) == false then
        ctx:fail("스톤골렘 생성 위치로 이동 실패")
        return nil
    end
    local p = bot:request(resp.spawn_mob, req.normal_chat { message = "/몬스터생성 " .. STONE_GOLEM }, function(p)
        return p.mob.mob_id == STONE_GOLEM
    end)
    if p == false then
        ctx:fail("스톤골렘 생성 응답 없음")
        return nil
    end
    return p.mob.oid
end

local function shadow_web_front(ctx)
    local bot = ctx:bot(0)
    if setup(ctx, bot, 120, HERMIT, "-") == false then
        return false
    end
    if pq.command(bot, "/스킬레벨 " .. SHADOW_WEB .. " 20", "스킬레벨 설정") == false then
        return ctx:fail("쉐도우 웹 배우기 실패")
    end
    if pq.command(bot, "/마력바꾸기 3000", "마력 설정") == false then
        return ctx:fail("마력 설정 실패")
    end
    if bot:instance_move(HUNTING_GROUND) == false then
        return ctx:fail("사냥터 이동 실패")
    end

    local x, y = bot:position()
    local behind = spawn_golem(ctx, bot, x - 300, y)
    if behind == nil then
        return false
    end
    local front = spawn_golem(ctx, bot, x - 100, y)
    if front == nil then
        return false
    end
    if pq.move(bot, x - 200, y) == false then
        return ctx:fail("시전 위치로 이동 실패")
    end

    local webbed_front = false
    local webbed_behind = false
    local function watch(p)
        if p.oid == behind then
            webbed_behind = true
        end
        if p.oid == front then
            webbed_front = true
        end
        return webbed_front
    end
    for _ = 1, 6 do
        if webbed_front then
            break
        end
        bot:request(resp.apply_mob_buff, req.active_skill { skill_id = SHADOW_WEB, skill_level = 20 }, watch, 1500)
    end
    bot:request(resp.apply_mob_buff, nil, watch, 500)
    if check(ctx, webbed_front, "오른쪽을 보고 쓴 쉐도우 웹이 앞의 몹에 걸리지 않음") == false then
        return false
    end
    return check(ctx, webbed_behind == false, "오른쪽을 보고 쓴 쉐도우 웹이 등 뒤 몹에 걸림")
end

local function spawn_snail(ctx, bot)
    local p = bot:request(resp.spawn_mob, req.normal_chat { message = "/몬스터생성 " .. SNAIL }, function(p)
        return p.mob.mob_id == SNAIL
    end)
    if p == false then
        ctx:fail("달팽이 생성 응답 없음")
        return nil
    end
    return p.mob.oid
end

local function stars(bot)
    return bot:items()[STAR]
end

local function ranged_ammo(ctx)
    local bot = ctx:bot(0)
    if setup(ctx, bot, 30, ROGUE, CLAW .. ":1," .. STAR .. ":3") == false then
        return false
    end
    if learn(ctx, bot, LUCKY_SEVEN) == false then
        return false
    end
    if pq.command(bot, "/마력바꾸기 3000", "마력 설정") == false then
        return ctx:fail("마력 설정 실패")
    end
    local wear = req.move_item { inventory_type = EQUIP, source = bot:slot(CLAW) or 0, dest = WEAPON_SLOT, count = 1 }
    if bot:request(resp.inventory_operation, wear, nil, 3000) == false then
        return ctx:fail("아대 장착 실패")
    end
    if bot:instance_move(HUNTING_GROUND) == false then
        return ctx:fail("사냥터 이동 실패")
    end

    local slot = bot:slot(STAR)
    local oid = spawn_snail(ctx, bot)
    if oid == nil then
        return false
    end
    if check(ctx, bot:shoot(oid, slot, LUCKY_SEVEN, 2), "표창 3개로 쓴 럭키 세븐이 맞지 않음") == false then
        return false
    end
    if check(ctx, stars(bot) == 1, "럭키 세븐 뒤 표창 개수가 다름: " .. tostring(stars(bot))) == false then
        return false
    end

    oid = spawn_snail(ctx, bot)
    if oid == nil then
        return false
    end
    if check(ctx, bot:shoot(oid, slot, LUCKY_SEVEN, 2) == false, "표창 1개로 쓴 럭키 세븐이 맞음") == false then
        return false
    end
    if check(ctx, stars(bot) == 1, "거부된 럭키 세븐이 표창을 소모함: " .. tostring(stars(bot))) == false then
        return false
    end

    if check(ctx, bot:shoot(oid, slot), "남은 표창 1개로 던진 공격이 맞지 않음") == false then
        return false
    end
    return check(ctx, bot:slot(STAR) == slot and stars(bot) == 0, "다 쓴 표창 칸이 0개로 남지 않음: " .. tostring(stars(bot)))
end

-- A poison mist poisons a mob standing in it. Each tick poisons at the level's chance (70% at level 30).
local function poison_mist(ctx)
    local bot = ctx:bot(1)
    if setup(ctx, bot, 120, FP_MAGE, "-") == false then
        return false
    end
    if pq.command(bot, "/스킬레벨 " .. POISON_MIST .. " 30", "스킬레벨 설정") == false then
        return ctx:fail("독 안개 배우기 실패")
    end
    if pq.command(bot, "/마력바꾸기 3000", "마력 설정") == false then
        return ctx:fail("마력 설정 실패")
    end
    if bot:instance_move(HUNTING_GROUND) == false then
        return ctx:fail("사냥터 이동 실패")
    end

    local x, y = bot:position()
    local golem = spawn_golem(ctx, bot, x, y)
    if golem == nil then
        return false
    end
    local bunny = bot:request(resp.spawn_mob, req.normal_chat { message = "/몬스터생성 " .. MOON_BUNNY }, function(p)
        return p.mob.mob_id == MOON_BUNNY
    end)
    if bunny == false then
        return ctx:fail("월묘 생성 응답 없음")
    end

    local bunny_poisoned = false
    local function poisoned(oid)
        return function(p)
            if p.oid == bunny.mob.oid and math.floor(p.status / MOB_POISON) % 2 == 1 then
                bunny_poisoned = true
            end
            return p.oid == oid and math.floor(p.status / MOB_POISON) % 2 == 1
        end
    end
    local p = bot:request(resp.apply_mob_buff, req.active_skill { skill_id = POISON_MIST, skill_level = 30 },
        poisoned(golem), 13000)
    if check(ctx, p ~= false, "독 안개 안의 스톤골렘이 중독되지 않음") == false then
        return false
    end
    -- A friendly mob stays in the mist for several more ticks without being poisoned.
    bot:request(resp.apply_mob_buff, nil, poisoned(bunny.mob.oid), 5000)
    -- The mist lasts 40 s and would poison the golem of a later scenario on this map.
    bot:request(resp.notice, req.normal_chat { message = "/맵리셋 인스턴스" }, nil, 300)
    return check(ctx, bunny_poisoned == false, "독 안개 안의 월묘가 중독됨")
end

-- A poison bomb leaves a poison mist where it is thrown, which poisons a mob that comes in later.
local function poison_bomb(ctx)
    local bot = ctx:bot(0)
    if setup(ctx, bot, 120, NIGHT_WALKER, CLAW .. ":1," .. STAR .. ":100") == false then
        return false
    end
    if pq.command(bot, "/스킬레벨 " .. POISON_BOMB .. " 30", "스킬레벨 설정") == false then
        return ctx:fail("포이즌 봄 배우기 실패")
    end
    if pq.command(bot, "/마력바꾸기 3000", "마력 설정") == false then
        return ctx:fail("마력 설정 실패")
    end
    local wear = req.move_item { inventory_type = EQUIP, source = bot:slot(CLAW) or 0, dest = WEAPON_SLOT, count = 1 }
    if bot:request(resp.inventory_operation, wear, nil, 3000) == false then
        return ctx:fail("아대 장착 실패")
    end
    if bot:instance_move(HUNTING_GROUND) == false then
        return ctx:fail("사냥터 이동 실패")
    end

    -- The flower lands on the foothold below, so the bomb is thrown from where it stands.
    local flower = bot:request(resp.spawn_mob, req.normal_chat { message = "/몬스터생성 " .. POISON_FLOWER }, function(p)
        return p.mob.mob_id == POISON_FLOWER
    end)
    if flower == false then
        return ctx:fail("포이즌 플라워 생성 응답 없음")
    end
    if pq.move(bot, flower.mob.position.x, flower.mob.position.y) == false then
        return ctx:fail("포이즌 플라워 위치로 이동 실패")
    end
    local snail = spawn_snail(ctx, bot)
    if snail == nil then
        return false
    end
    if check(ctx, bot:shoot(snail, bot:slot(STAR), POISON_BOMB), "포이즌 봄이 맞지 않음") == false then
        return false
    end
    if pq.command(bot, "/즉사 off", "즉사 상태: disabled") == false then
        return ctx:fail("즉사 끄기 실패")
    end
    local p = bot:request(resp.apply_mob_buff, nil, function(p)
        return p.oid == flower.mob.oid and math.floor(p.status / MOB_POISON) % 2 == 1
    end, 10000)
    -- The mist would poison the golem of a later scenario on this map.
    bot:request(resp.notice, req.normal_chat { message = "/맵리셋 인스턴스" }, nil, 300)
    return check(ctx, p ~= false, "포이즌 봄 독구름 안의 몹이 중독되지 않음")
end

local function hp_changed(p)
    return p.stats[HP] ~= nil
end

-- The GM HP command sets the max too, so the HP a potion can fill is made by a hit.
local function potion_heal(ctx, bot)
    if pq.command(bot, "/체력바꾸기 1000", "체력 설정") == false then
        ctx:fail(bot:name() .. " 체력 설정 실패")
        return nil
    end
    local hit = bot:request(resp.update_stats, req.damaged { type = HIT_ENV, damage = 900 }, hp_changed, 3000)
    if hit == false then
        ctx:fail(bot:name() .. " 피해를 받지 않음")
        return nil
    end
    local p = bot:request(resp.update_stats, req.use_item { slot = bot:slot(RED_POTION) or 0, item_id = RED_POTION },
        hp_changed, 3000)
    if p == false then
        ctx:fail(bot:name() .. " 빨간 포션을 먹었는데 체력이 오르지 않음")
        return nil
    end
    return p.stats[HP] - hit.stats[HP]
end

-- Alchemist at level 20 makes a potion heal 150%, for the Night Walker as for the Hermit.
local function alchemist(ctx, bot, class, skill_id)
    if setup(ctx, bot, 70, class, RED_POTION .. ":2") == false then
        return false
    end
    if bot:instance_move(HUNTING_GROUND) == false then
        return ctx:fail("사냥터 이동 실패")
    end
    local base = potion_heal(ctx, bot)
    if base == nil then
        return false
    end
    if pq.command(bot, "/스킬레벨 " .. skill_id .. " 20", "스킬레벨 설정") == false then
        return ctx:fail("알케미스트 배우기 실패: " .. skill_id)
    end
    local boosted = potion_heal(ctx, bot)
    if boosted == nil then
        return false
    end
    return check(ctx, boosted * 2 == base * 3,
        string.format("알케미스트 %d 20레벨 포션 회복이 150%%가 아님: %d -> %d", skill_id, base, boosted))
end

local function hermit_alchemist(ctx)
    return alchemist(ctx, ctx:bot(0), HERMIT, ALCHEMIST)
end

local function night_walker_alchemist(ctx)
    return alchemist(ctx, ctx:bot(1), NIGHT_WALKER, ALCHEMIST_CYGNUS)
end

-- A magician's natural MP recovery grows with Improving MP Recovery and the character level (16 x 70 / 10 = 112).
local function magician_mp_recovery(ctx)
    local bot = ctx:bot(0)
    if setup(ctx, bot, 70, FP_MAGE, "-") == false then
        return false
    end
    if pq.command(bot, "/스킬레벨 " .. IMPROVING_MP_RECOVERY .. " 16", "스킬레벨 설정") == false then
        return ctx:fail("MP 회복력 향상 배우기 실패")
    end
    if bot:instance_move(HUNTING_GROUND) == false then
        return ctx:fail("사냥터 이동 실패")
    end
    if pq.command(bot, "/체력바꾸기 3000", "체력 설정") == false or pq.command(bot, "/마력바꾸기 1000", "마력 설정") == false then
        return ctx:fail("체력·마력 설정 실패")
    end
    local spawned = bot:request(resp.spawn_mob, req.normal_chat { message = "/몬스터생성 " .. MP_BURN_MOB }, function(p)
        return p.mob.mob_id == MP_BURN_MOB
    end, 5000)
    if spawned == false then
        return ctx:fail("MP 흡수 몹 생성 응답 없음")
    end
    local function mp_changed(p)
        return p.stats[MP] ~= nil
    end
    local burned = bot:request(resp.update_stats,
        req.damaged { type = 0, damage = 1, mob_id = MP_BURN_MOB, oid = spawned.mob.oid }, mp_changed, 3000)
    if burned == false then
        return ctx:fail("MP 흡수 공격에 MP가 줄지 않음")
    end
    local p = bot:request(resp.update_stats, req.heal_over_time { heal_mp = 200 }, mp_changed, 3000)
    if p == false then
        return ctx:fail("자연 회복으로 MP가 오르지 않음")
    end
    local gain = p.stats[MP] - burned.stats[MP]
    return check(ctx, gain > 100 and gain < 200, "MP 회복력 향상 16레벨 70레벨 마법사의 자연 회복 상한이 다름: " .. gain)
end

-- A sauna doubles natural HP recovery (10 x 2.0), and a sauna robe makes it half again (10 x 3.0).
local function sauna_robe_recovery(ctx)
    local bot = ctx:bot(0)
    local items = string.format("%d:1,%d:1", SAUNA_ROBES[1], SAUNA_ROBES[2])
    if setup(ctx, bot, 30, 0, items) == false then
        return false
    end
    if bot:force_instance_move(SAUNA, "사우나 자연 회복") == false then
        return ctx:fail("사우나 이동 실패")
    end
    if pq.command(bot, "/체력바꾸기 3000", "체력 설정") == false then
        return ctx:fail("체력 설정 실패")
    end
    local function hp_changed(p)
        return p.stats[HP] ~= nil
    end
    local function heal()
        local hurt = bot:request(resp.update_stats, req.damaged { type = HIT_ENV, damage = 1000 }, hp_changed, 3000)
        if hurt == false then
            return nil
        end
        local p = bot:request(resp.update_stats, req.heal_over_time { heal_hp = 200 }, hp_changed, 3000)
        if p == false then
            return nil
        end
        return p.stats[HP] - hurt.stats[HP]
    end

    local plain = heal()
    if check(ctx, plain == 20, "사우나의 자연 회복이 2배가 아님: " .. tostring(plain)) == false then
        return false
    end
    local worn = false
    for _, robe in ipairs(SAUNA_ROBES) do
        bot:request(resp.inventory_operation,
            req.move_item { inventory_type = EQUIP, source = bot:slot(robe) or 0, dest = OVERALL_SLOT, count = 1 }, nil,
            2000)
        ctx:sleep(500)
        if bot:slot(robe) == nil then
            worn = true
            break
        end
    end
    if check(ctx, worn, "사우나 가운을 입지 못함") == false then
        return false
    end
    ctx:sleep(HP_HEAL_INTERVAL_MS)
    local robed = heal()
    return check(ctx, robed == 30, "사우나 가운을 입은 자연 회복이 1.5배가 아님: " .. tostring(robed))
end

-- A skill takes its MP once, and a use refused for want of MP takes nothing.
local function skill_cost(ctx)
    local bot = ctx:bot(1)
    if setup(ctx, bot, 30, MAGICIAN, "-") == false or learn(ctx, bot, MAGIC_GUARD) == false then
        return false
    end
    if bot:instance_move(HUNTING_GROUND) == false then
        return ctx:fail("사냥터 이동 실패")
    end
    local function mp_changed(p)
        return p.stats[MP] ~= nil
    end
    if pq.command(bot, "/마력바꾸기 100", "마력 설정") == false then
        return ctx:fail("마력 설정 실패")
    end
    local paid = bot:request(resp.update_stats, req.active_skill { skill_id = MAGIC_GUARD, skill_level = 1 }, mp_changed,
        3000)
    if check(ctx, paid ~= false and paid.stats[MP] == 100 - MAGIC_GUARD_MP,
            "매직 가드 MP 비용이 한 번만 빠지지 않음: " .. tostring(paid and paid.stats[MP])) == false then
        return false
    end
    if pq.command(bot, "/마력바꾸기 " .. (MAGIC_GUARD_MP - 1), "마력 설정") == false then
        return ctx:fail("마력 설정 실패")
    end
    local refused = bot:request(resp.update_stats, req.active_skill { skill_id = MAGIC_GUARD, skill_level = 1 },
        mp_changed, 1500)
    return check(ctx, refused == false, "MP가 모자라 거절된 매직 가드가 MP를 바꿈")
end

local function has_flag(status, flag)
    return math.floor(status / flag) % 2 == 1
end

-- Threaten shares only the weapon defense flag with Showdown, so it takes that one flag over and leaves Showdown's
-- other statuses (the EXP and drop bonus, the magic defense) on the mob.
local function showdown_threaten(ctx)
    local thief, page = ctx:bot(0), ctx:bot(1)
    if setup(ctx, thief, 120, NIGHT_LORD, CLAW .. ":1," .. STAR .. ":100") == false or setup(ctx, page, 70, PAGE, "-") == false then
        return false
    end
    local wear = req.move_item { inventory_type = EQUIP, source = thief:slot(CLAW) or 0, dest = WEAPON_SLOT, count = 1 }
    if thief:request(resp.inventory_operation, wear, nil, 3000) == false then
        return ctx:fail("아대 장착 실패")
    end
    if pq.command(thief, "/스킬레벨 " .. SHOWDOWN .. " 30", "스킬레벨 설정") == false
        or pq.command(page, "/스킬레벨 " .. THREATEN .. " 20", "스킬레벨 설정") == false then
        return ctx:fail("쇼다운·위협 배우기 실패")
    end
    if pq.command(thief, "/마력바꾸기 3000", "마력 설정") == false or pq.command(page, "/마력바꾸기 3000", "마력 설정") == false then
        return ctx:fail("마력 설정 실패")
    end
    -- A damage check against the golem needs the thief to hit it, so its stats are raised and instant kill is off.
    if pq.command(thief, "/전능 999", "전능 설정") == false or pq.command(thief, "/즉사 off", "즉사 상태: disabled") == false then
        return ctx:fail("능력치·즉사 설정 실패")
    end
    if thief:instance_move(HUNTING_GROUND) == false or page:instance_move(HUNTING_GROUND) == false then
        return ctx:fail("사냥터 이동 실패")
    end
    -- Threaten takes the first six mobs in its range, so the map's own mobs are cleared first.
    if pq.command(thief, "/몬스터죽이기", "마리 제거") == false then
        return ctx:fail("사냥터 몬스터 제거 실패")
    end
    local x, y = page:position()
    local golem = spawn_golem(ctx, thief, x, y)
    if golem == nil then
        return false
    end

    local applied, cancelled = {}, {}
    ctx:hook("showdown_threaten", function(_, bot, p, name)
        if bot:id() ~= thief:id() or p.oid ~= golem then
            return
        end
        local seen = nil
        if name == "apply_mob_buff" then
            seen = applied
        elseif name == "cancel_mob_buff" then
            seen = cancelled
        else
            return
        end
        for _, flag in ipairs({ MOB_WATK, MOB_MDEF, MOB_SHOWDOWN }) do
            if has_flag(p.status, flag) then
                seen[flag] = true
            end
        end
    end)
    thief:attack(golem, 1, 1, SHOWDOWN)
    ctx:sleep(500)
    local showdown = applied[MOB_SHOWDOWN] == true
    page:send(req.active_skill { skill_id = THREATEN, skill_level = 20 })
    ctx:sleep(1500)
    ctx:unhook("showdown_threaten")
    if check(ctx, showdown, "쇼다운이 스톤골렘에 걸리지 않음") == false then
        return false
    end
    if check(ctx, applied[MOB_WATK] == true, "위협이 스톤골렘에 걸리지 않음") == false then
        return false
    end
    if check(ctx, cancelled[MOB_SHOWDOWN] == nil, "위협이 쇼다운을 풀어 버림") == false then
        return false
    end
    return check(ctx, cancelled[MOB_MDEF] == nil, "위협이 쇼다운의 마법 방어 감소를 풀어 버림")
end

-- The attack is sent by the request itself, so its reply cannot arrive before the wait starts.
local function melee(bot, oid, damage, skill_id)
    local x, y = bot:position()
    return req.attack {
        close_attack_info = {
            attack_header = { targets = 1, hits = 1, skill = skill_id },
            damages = { { oid = oid, damage_pairs = { { damage = damage } } } },
            position = { x = x, y = y },
        },
    }
end

-- Soul charge gives the sword the holy element, so a lightning-immune mob still takes the hit.
local function soul_charge(ctx)
    local bot = ctx:bot(0)
    if setup(ctx, bot, 70, DAWN_WARRIOR3, SWORD .. ":1") == false then
        return false
    end
    local wear = req.move_item { inventory_type = EQUIP, source = bot:slot(SWORD) or 0, dest = WEAPON_SLOT, count = 1 }
    if bot:request(resp.inventory_operation, wear, nil, 3000) == false then
        return ctx:fail("검 착용 실패")
    end
    if pq.command(bot, "/스킬레벨 " .. SOUL_CHARGE .. " 20", "스킬레벨 설정") == false then
        return ctx:fail("소울 차지 배우기 실패")
    end
    -- Starting stats leave the mob's defense above any hit.
    if pq.command(bot, "/전능 999", "전능 설정") == false or pq.command(bot, "/마력바꾸기 3000", "마력 설정") == false then
        return ctx:fail("능력치 설정 실패")
    end
    if pq.command(bot, "/즉사 off", "즉사 상태: disabled") == false then
        return ctx:fail("즉사 끄기 실패")
    end
    if bot:instance_move(HUNTING_GROUND) == false then
        return ctx:fail("사냥터 이동 실패")
    end
    if bot:request(resp.update_stats, req.active_skill { skill_id = SOUL_CHARGE, skill_level = 20 }, function(p)
            return p.unlock_action
        end, 3000) == false then
        return ctx:fail("소울 차지 사용 실패")
    end

    local spawned = bot:request(resp.spawn_mob, req.normal_chat { message = "/몬스터생성 " .. LIGHTNING_IMMUNE }, function(p)
        return p.mob.mob_id == LIGHTNING_IMMUNE
    end)
    if spawned == false then
        return ctx:fail("번개 면역 몹 생성 응답 없음")
    end
    local hit = bot:request(resp.show_mob_hp, melee(bot, spawned.mob.oid, 10, 0), function(p)
        return p.oid == spawned.mob.oid
    end, 3000)
    return check(ctx, hit ~= false, "소울 차지 공격이 번개 면역 몹에 들어가지 않음")
end

-- A finisher needs a combo orb; without one the whole attack is dropped.
local function finisher_orbs(ctx)
    local bot = ctx:bot(0)
    if setup(ctx, bot, 70, CRUSADER, SWORD .. ":1") == false then
        return false
    end
    local wear = req.move_item { inventory_type = EQUIP, source = bot:slot(SWORD) or 0, dest = WEAPON_SLOT, count = 1 }
    if bot:request(resp.inventory_operation, wear, nil, 3000) == false then
        return ctx:fail("검 착용 실패")
    end
    if pq.command(bot, "/스킬레벨 " .. COMBO_ATTACK .. " 30", "스킬레벨 설정") == false
        or pq.command(bot, "/스킬레벨 " .. COMA .. " 30", "스킬레벨 설정") == false then
        return ctx:fail("콤보 어택·코마 배우기 실패")
    end
    -- Starting stats leave the mob's defense above any hit.
    if pq.command(bot, "/전능 999", "전능 설정") == false or pq.command(bot, "/마력바꾸기 3000", "마력 설정") == false then
        return ctx:fail("능력치 설정 실패")
    end
    if pq.command(bot, "/즉사 off", "즉사 상태: disabled") == false then
        return ctx:fail("즉사 끄기 실패")
    end
    if bot:instance_move(HUNTING_GROUND) == false then
        return ctx:fail("사냥터 이동 실패")
    end
    if bot:request(resp.update_stats, req.active_skill { skill_id = COMBO_ATTACK, skill_level = 30 }, function(p)
            return p.unlock_action
        end, 3000) == false then
        return ctx:fail("콤보 어택 사용 실패")
    end
    local x, y = bot:position()
    local golem = spawn_golem(ctx, bot, x, y)
    if golem == nil then
        return false
    end
    local function hit(skill_id, timeout)
        return bot:request(resp.show_mob_hp, melee(bot, golem, 10, skill_id), function(p)
            return p.oid == golem
        end, timeout) ~= false
    end

    if hit(COMA, QUIET_MS) then
        return ctx:fail("구슬 없는 코마가 들어감")
    end
    if hit(0, 3000) == false then
        return ctx:fail("일반 공격이 들어가지 않음")
    end
    return check(ctx, hit(COMA, 3000), "구슬을 모은 뒤 코마가 들어가지 않음")
end

-- Sacrifice costs the attacker a part of the damage dealt, a boss included.
local function sacrifice_boss(ctx)
    local bot = ctx:bot(0)
    if setup(ctx, bot, 70, DRAGON_KNIGHT, SPEAR .. ":1") == false then
        return false
    end
    local wear = req.move_item { inventory_type = EQUIP, source = bot:slot(SPEAR) or 0, dest = WEAPON_SLOT, count = 1 }
    if bot:request(resp.inventory_operation, wear, nil, 3000) == false then
        return ctx:fail("창 착용 실패")
    end
    if pq.command(bot, "/스킬레벨 " .. SACRIFICE .. " 30", "스킬레벨 설정") == false then
        return ctx:fail("새크리파이스 배우기 실패")
    end
    if pq.command(bot, "/전능 999", "전능 설정") == false or pq.command(bot, "/마력바꾸기 3000", "마력 설정") == false
        or pq.command(bot, "/체력바꾸기 3000", "체력 설정") == false then
        return ctx:fail("능력치 설정 실패")
    end
    if pq.command(bot, "/즉사 off", "즉사 상태: disabled") == false then
        return ctx:fail("즉사 끄기 실패")
    end
    if bot:instance_move(HUNTING_GROUND) == false then
        return ctx:fail("사냥터 이동 실패")
    end
    local spawned = bot:request(resp.spawn_mob, req.normal_chat { message = "/몬스터생성 " .. MUSHMOM }, function(p)
        return p.mob.mob_id == MUSHMOM
    end)
    if spawned == false then
        return ctx:fail("머쉬맘 생성 응답 없음")
    end

    local p = bot:request(resp.update_stats, melee(bot, spawned.mob.oid, 1000, SACRIFICE), function(p)
        return p.stats[HP] ~= nil and p.stats[HP] < 3000
    end, 3000)
    if p == false then
        return ctx:fail("보스에게 쓴 새크리파이스가 HP를 깎지 않음")
    end
    return check(ctx, p.stats[HP] == 3000 - 1000 * SACRIFICE_HP_PERCENT / 100,
        "새크리파이스 HP 감소량이 피해의 5%가 아님: " .. tostring(p.stats[HP]))
end

local function buff_value(p, flag)
    for _, entry in ipairs(p.buffs or {}) do
        if entry.buff.mask == flag.mask and entry.buff.position == flag.position then
            return entry.value
        end
    end
    return nil
end

-- The GM hyper body reaches the living players in its area, not only the GM.
local function gm_hyper_body(ctx)
    local gm = ctx:bot(1)
    local other = ctx:bot(0)
    if setup(ctx, gm, 120, GM, "-") == false or setup(ctx, other, 30, 0, "-") == false then
        return false
    end
    if learn(ctx, gm, GM_HYPER_BODY) == false then
        return false
    end
    if gm:instance_move(GM_GROUND) == false or other:instance_move(GM_GROUND) == false then
        return ctx:fail("사냥터 이동 실패")
    end
    local x, y = gm:position()
    if pq.move(other, x, y) == false then
        return ctx:fail("GM 옆으로 이동 실패")
    end
    local p = gm:request_on(other, resp.update_self_buff,
        req.active_skill { skill_id = GM_HYPER_BODY, skill_level = 1 }, function(p)
            return p.buff_id == GM_HYPER_BODY
        end, 3000)
    if p == false then
        return ctx:fail("범위 안의 다른 캐릭터가 GM 하이퍼 바디를 받지 않음")
    end
    return check(ctx, buff_value(p, BUFF.max_hp) == 60,
        "GM 하이퍼 바디 최대 HP 증가가 60%가 아님: " .. tostring(buff_value(p, BUFF.max_hp)))
end

-- A transformation carries the skill's defense, speed, and jump along with the look.
local function transform_stats(ctx)
    local bot = ctx:bot(0)
    if setup(ctx, bot, 70, MARAUDER, "-") == false or learn(ctx, bot, TRANSFORM) == false then
        return false
    end
    if pq.command(bot, "/마력바꾸기 3000", "마력 설정") == false then
        return ctx:fail("마력 설정 실패")
    end
    if bot:instance_move(HUNTING_GROUND) == false then
        return ctx:fail("사냥터 이동 실패")
    end
    local p = bot:request(resp.update_self_buff, req.active_skill { skill_id = TRANSFORM, skill_level = 1 }, function(p)
        return p.buff_id == TRANSFORM
    end, 3000)
    if p == false then
        return ctx:fail("트랜스폼 버프 응답 없음")
    end
    local got = string.format("morph=%s speed=%s jump=%s wdef=%s mdef=%s", tostring(buff_value(p, BUFF.morph)),
        tostring(buff_value(p, BUFF.speed)), tostring(buff_value(p, BUFF.jump)), tostring(buff_value(p, BUFF.weapon_def)),
        tostring(buff_value(p, BUFF.magic_def)))
    return check(ctx, (buff_value(p, BUFF.morph) or 0) > 0 and buff_value(p, BUFF.speed) == 40 and
        buff_value(p, BUFF.jump) == 20 and buff_value(p, BUFF.weapon_def) == 2 and buff_value(p, BUFF.magic_def) == 2,
        "트랜스폼 버프 능력치가 WZ와 다름: " .. got)
end

local function dark_sight_speed(ctx)
    local bot = ctx:bot(0)
    if setup(ctx, bot, 30, ROGUE, "-") == false or learn(ctx, bot, DARK_SIGHT) == false then
        return false
    end
    if pq.command(bot, "/마력바꾸기 1000", "마력 설정") == false then
        return ctx:fail("마력 설정 실패")
    end
    local p = bot:request(resp.update_self_buff, req.active_skill { skill_id = DARK_SIGHT, skill_level = 1 }, function(p)
        return p.buff_id == DARK_SIGHT
    end, 3000)
    if p == false then
        return ctx:fail("다크 사이트 버프 응답 없음")
    end
    return check(ctx, buff_value(p, BUFF.darksight) == 1 and buff_value(p, BUFF.speed) == DARK_SIGHT_SPEED,
        string.format("다크 사이트 버프가 WZ와 다름: darksight=%s speed=%s (기대 1, %d)",
            tostring(buff_value(p, BUFF.darksight)), tostring(buff_value(p, BUFF.speed)), DARK_SIGHT_SPEED))
end

-- Dragon Roar under half HP is refused before it hits; above half it costs the WZ share of max HP.
local function dragon_roar_half_hp(ctx)
    local bot = ctx:bot(0)
    if setup(ctx, bot, 70, DRAGON_KNIGHT, SPEAR .. ":1") == false then
        return false
    end
    local wear = req.move_item { inventory_type = EQUIP, source = bot:slot(SPEAR) or 0, dest = WEAPON_SLOT, count = 1 }
    if bot:request(resp.inventory_operation, wear, nil, 3000) == false then
        return ctx:fail("창 착용 실패")
    end
    if learn(ctx, bot, DRAGON_ROAR) == false then
        return false
    end
    if pq.command(bot, "/전능 999", "전능 설정") == false or pq.command(bot, "/마력바꾸기 3000", "마력 설정") == false
        or pq.command(bot, "/체력바꾸기 3000", "체력 설정") == false then
        return ctx:fail("능력치 설정 실패")
    end
    if pq.command(bot, "/즉사 off", "즉사 상태: disabled") == false then
        return ctx:fail("즉사 끄기 실패")
    end
    if bot:instance_move(HUNTING_GROUND) == false then
        return ctx:fail("사냥터 이동 실패")
    end
    local x, y = bot:position()
    local golem = spawn_golem(ctx, bot, x, y)
    if golem == nil then
        return false
    end

    for _ = 1, 2 do
        local hurt = bot:request(resp.update_stats, req.damaged { type = HIT_ENV, damage = 1000 }, function(p)
            return p.stats[HP] ~= nil
        end, 3000)
        if hurt == false then
            return ctx:fail("피해를 받지 않음")
        end
    end
    if bot:request(resp.show_mob_hp, melee(bot, golem, 10, DRAGON_ROAR), function(p)
            return p.oid == golem
        end, QUIET_MS) ~= false then
        return ctx:fail("HP 절반 아래에서 드래곤 로어가 들어감")
    end

    if pq.command(bot, "/체력바꾸기 3000", "체력 설정") == false then
        return ctx:fail("체력 회복 실패")
    end
    local p = bot:request(resp.update_stats, melee(bot, golem, 10, DRAGON_ROAR), function(p)
        return p.stats[HP] ~= nil and p.stats[HP] < 3000
    end, 3000)
    if p == false then
        return ctx:fail("HP 절반 이상에서 드래곤 로어가 HP를 쓰지 않음")
    end
    local expected = 3000 - math.floor(3000 * DRAGON_ROAR_HP_PERCENT / 100)
    if check(ctx, p.stats[HP] == expected, string.format("드래곤 로어 HP 소모가 다름: %d (기대 %d)", p.stats[HP], expected)) == false then
        return false
    end
    ctx:sleep(DRAGON_ROAR_STUN_MS + 500)
    return true
end

test_suite {
    name = "Skill: 미스틱 도어 맵 제한·사망 시 사용 불가·GM 홀리 심볼·솔로 파티 버프·스킬 범위 방향·원거리 탄 소모·독 안개·알케미스트·MP 회복력 향상·쇼다운 뒤 위협·스킬 MP 비용·GM 하이퍼 바디·변신 능력치·드래곤 로어 HP 조건·다크 사이트 감속",
    bot_count = 2,

    scenarios = {
        mystic_door,
        dead_mystic_door,
        gm_holy_symbol,
        solo_holy_symbol,
        shadow_web_front,
        ranged_ammo,
        poison_mist,
        poison_bomb,
        hermit_alchemist,
        night_walker_alchemist,
        magician_mp_recovery,
        showdown_threaten,
        soul_charge,
        finisher_orbs,
        sacrifice_boss,
        skill_cost,
        sauna_robe_recovery,
        dragon_roar_half_hp,
        transform_stats,
        dark_sight_speed,
        -- Last: the hyper body stays on the other bot and would raise its max HP.
        gm_hyper_body,
    },
}
