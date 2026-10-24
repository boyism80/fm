local pq = require("script/integration/lib/party_quest")

local HUNTING_GROUND = 104040000
local DISEASE_MOB = 5000007
local MP_BURN_MOB = 8110300
local BANISH_MOB = 9300160
local BANISH_MAP = 925010000
local STUN = 123
local SEAL = 120
local ZOMBIFY = 133
local HIT_COLLIDE = -1
local HIT_MAP_DEBUFF = -2
local HIT_ENV = -3
local HP = 0x400
local MP = 0x1000
local BISHOP = 232
local HOLY_SHIELD = 2321005
local RED_POTION = 2000000

local function setup(ctx, bot)
    if pq.command(bot, "/봇초기화 70 0 0 -", "봇초기화 완료") == false then
        return ctx:fail(bot:name() .. " 봇 초기화 실패")
    end
    if pq.command(bot, "/체력바꾸기 3000", "체력 설정") == false then
        return ctx:fail(bot:name() .. " 체력 설정 실패")
    end
    if bot:instance_move(HUNTING_GROUND) == false then
        return ctx:fail(bot:name() .. " 사냥터 이동 실패")
    end
    return true
end

local function spawn(ctx, bot, mob_id)
    local p = bot:request(resp.spawn_mob, req.normal_chat { message = "/몬스터생성 " .. mob_id }, function(p)
        return p.mob.mob_id == mob_id
    end, 5000)
    if p == false then
        return ctx:fail(mob_id .. " 생성 응답 없음")
    end
    return p.mob.oid
end

local function hit(mob_id, oid, attack, damage)
    return req.damaged { type = attack, damage = damage, mob_id = mob_id, oid = oid }
end

-- A mob attack carries its WZ disease, and a map hazard carries the mob skill in the packet.
local function diseases(ctx)
    local actor, observer = ctx:bot(0), ctx:bot(1)
    if setup(ctx, actor) == false or setup(ctx, observer) == false then
        return false
    end
    local oid = spawn(ctx, actor, DISEASE_MOB)
    if oid == false then
        return false
    end

    local stunned = actor:request_on(observer, resp.give_debuff, hit(DISEASE_MOB, oid, 0, 1), function(p)
        return p.character_id == actor:id() and p.skill_id == STUN
    end, 5000)
    if stunned == false then
        return ctx:fail("attack1에 맞았는데 기절이 걸리지 않음")
    end

    local sealed = actor:request_on(observer, resp.give_debuff,
        req.damaged { type = HIT_MAP_DEBUFF, damage = 1, level = 1, skill_id = SEAL }, function(p)
            return p.character_id == actor:id() and p.skill_id == SEAL
        end, 5000)
    if sealed == false then
        return ctx:fail("맵 디버프 피격에 봉인이 걸리지 않음")
    end
    return true
end

-- A deadly attack leaves 1 HP and 1 MP, and an MP burn attack takes MP besides the damage.
local function hp_mp(ctx)
    local actor = ctx:bot(0)
    if setup(ctx, actor) == false then
        return false
    end

    local oid = spawn(ctx, actor, MP_BURN_MOB)
    if oid == false then
        return false
    end
    local burned = actor:request(resp.update_stats, hit(MP_BURN_MOB, oid, 0, 1), function(p)
        return p.stats[MP] ~= nil
    end, 5000)
    if burned == false then
        return ctx:fail("MP 흡수 공격에 맞았는데 MP가 줄지 않음")
    end

    oid = spawn(ctx, actor, DISEASE_MOB)
    if oid == false then
        return false
    end
    local deadly = actor:request(resp.update_stats, hit(DISEASE_MOB, oid, 3, 1), function(p)
        return p.stats[HP] ~= nil
    end, 5000)
    if deadly == false or deadly.stats[HP] ~= 1 then
        return ctx:fail("HP를 1로 만드는 공격 뒤 HP " .. tostring(deadly and deadly.stats[HP]))
    end
    return true
end

-- A mob that banishes on touch sends the player to its ban map.
local function banish(ctx)
    local actor = ctx:bot(0)
    if setup(ctx, actor) == false then
        return false
    end
    local oid = spawn(ctx, actor, BANISH_MOB)
    if oid == false then
        return false
    end
    local warp = actor:request(resp.warp, hit(BANISH_MOB, oid, HIT_COLLIDE, 1), function(p)
        return p.character.map == BANISH_MAP
    end, 5000)
    if warp == false then
        return ctx:fail("몸에 닿으면 추방하는 몹에 닿았는데 추방되지 않음")
    end
    return true
end

-- A touch that kills leaves the body where it fell instead of banishing it.
local function lethal_touch(ctx)
    local actor = ctx:bot(0)
    if setup(ctx, actor) == false then
        return false
    end
    local oid = spawn(ctx, actor, BANISH_MOB)
    if oid == false then
        return false
    end
    local dead = actor:request(resp.update_stats, hit(BANISH_MOB, oid, HIT_COLLIDE, 30000), function(p)
        return p.stats[HP] == 0
    end, 5000)
    if dead == false then
        return ctx:fail("죽을 만큼 맞았는데 HP가 0이 되지 않음")
    end
    if actor:request(resp.warp, nil, nil, 1500) ~= false then
        return ctx:fail("죽은 캐릭터가 추방됨")
    end
    return true
end

local function map_debuff(skill_id)
    return req.damaged { type = HIT_MAP_DEBUFF, damage = 1, level = 1, skill_id = skill_id }
end

-- Holy Shield keeps off every disease but stun and seduce. The bishop is the bot the earlier scenarios left without
-- diseases, since a disease already there is not given again.
local function holy_shield(ctx)
    local actor, observer = ctx:bot(1), ctx:bot(0)
    if setup(ctx, observer) == false then
        return false
    end
    if pq.command(actor, "/봇초기화 120 " .. BISHOP .. " 0 -", "봇초기화 완료") == false then
        return ctx:fail("비숍 봇 초기화 실패")
    end
    if pq.command(actor, "/스킬레벨 " .. HOLY_SHIELD .. " 1", "스킬레벨 설정") == false then
        return ctx:fail("홀리 실드 배우기 실패")
    end
    if pq.command(actor, "/체력바꾸기 3000", "체력 설정") == false or pq.command(actor, "/마력바꾸기 3000", "마력 설정") == false then
        return ctx:fail("비숍 체력·마력 설정 실패")
    end
    if actor:instance_move(HUNTING_GROUND) == false then
        return ctx:fail("비숍 사냥터 이동 실패")
    end
    local cast = actor:request(resp.update_stats, req.active_skill { skill_id = HOLY_SHIELD, skill_level = 1 },
        function(p)
            return p.unlock_action
        end, 3000)
    if cast == false then
        return ctx:fail("홀리 실드 사용 응답 없음")
    end

    local function given(skill_id)
        return function(p)
            return p.character_id == actor:id() and p.skill_id == skill_id
        end
    end
    if actor:request_on(observer, resp.give_debuff, map_debuff(SEAL), given(SEAL), 1500) ~= false then
        return ctx:fail("홀리 실드 중에 봉인이 걸림")
    end
    if actor:request_on(observer, resp.give_debuff, map_debuff(STUN), given(STUN), 5000) == false then
        return ctx:fail("홀리 실드 중에 기절이 걸리지 않음")
    end
    return true
end

-- A zombified character recovers half the HP a potion gives.
local function zombify_potion(ctx)
    local actor = ctx:bot(0)
    if pq.command(actor, "/봇초기화 70 0 0 " .. RED_POTION .. ":2", "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    if actor:instance_move(HUNTING_GROUND) == false then
        return ctx:fail("사냥터 이동 실패")
    end
    local function hp_changed(p)
        return p.stats[HP] ~= nil
    end
    local function heal()
        if pq.command(actor, "/체력바꾸기 1000", "체력 설정") == false then
            return nil
        end
        local hit = actor:request(resp.update_stats, req.damaged { type = HIT_ENV, damage = 900 }, hp_changed, 3000)
        if hit == false then
            return nil
        end
        local p = actor:request(resp.update_stats,
            req.use_item { slot = actor:slot(RED_POTION) or 0, item_id = RED_POTION }, hp_changed, 3000)
        if p == false then
            return nil
        end
        return p.stats[HP] - hit.stats[HP]
    end

    local base = heal()
    if base == nil then
        return ctx:fail("언데드 전 포션 회복 확인 실패")
    end
    if actor:request(resp.update_stats, map_debuff(ZOMBIFY), hp_changed, 3000) == false then
        return ctx:fail("언데드 맵 디버프 피격 응답 없음")
    end
    local zombified = heal()
    if zombified == nil then
        return ctx:fail("언데드 중 포션 회복 확인 실패")
    end
    if zombified ~= math.floor(base / 2) then
        return ctx:fail(string.format("언데드 중 포션 회복이 절반이 아님: %d -> %d", base, zombified))
    end
    return true
end

test_suite {
    name = "Mob attack: 상태이상·맵 디버프·HP/MP 1·MP 흡수·몸박치기 추방·홀리 실드·언데드 포션",
    bot_count = 2,

    scenarios = {
        diseases,
        hp_mp,
        banish,
        lethal_touch,
        holy_shield,
        zombify_potion,
    },
}
