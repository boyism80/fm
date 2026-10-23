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
local PRIEST = 231
local ROGUE = 400
local HERMIT = 411
local FP_MAGE = 211
local GM = 900
local EQUIP = 1
local WEAPON_SLOT = -11

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
    local p = bot:request(resp.apply_mob_buff, req.active_skill { skill_id = POISON_MIST, skill_level = 30 }, function(p)
        return p.oid == golem and math.floor(p.status / MOB_POISON) % 2 == 1
    end, 13000)
    return check(ctx, p ~= false, "독 안개 안의 스톤골렘이 중독되지 않음")
end

test_suite {
    name = "Skill: 미스틱 도어 맵 제한·사망 시 사용 불가·GM 홀리 심볼·솔로 파티 버프·스킬 범위 방향·원거리 탄 소모·독 안개",
    bot_count = 2,

    scenarios = {
        mystic_door,
        dead_mystic_door,
        gm_holy_symbol,
        solo_holy_symbol,
        shadow_web_front,
        ranged_ammo,
        poison_mist,
    },
}
