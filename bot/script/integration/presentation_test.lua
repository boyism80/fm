local pq = require("script/integration/lib/party_quest")

local HUNTING_GROUND = 104040000
local BRAWLER = 510
local CORKSCREW_BLOW = 5101004
local KNUCKLE_BOOSTER = 5101006
local VOMIT = 5160000
local PANIC = 5160001
local SUN_EFFECT = 5010000
local MOON_EFFECT = 5010001
local QUIET_MS = 1500

local function setup(ctx, bot, level, class, items)
    if pq.command(bot, "/봇초기화 " .. level .. " " .. class .. " 0 " .. items, "봇초기화 완료") == false then
        return ctx:fail(bot:name() .. " 봇 초기화 실패")
    end
    return true
end

local function meet(ctx, actor, observer)
    if actor:instance_move(HUNTING_GROUND) == false or observer:instance_move(HUNTING_GROUND) == false then
        return ctx:fail("사냥터 이동 실패")
    end
    return true
end

local function seen(actor, observer, name, pkt, match, timeout)
    return actor:request_on(observer, name, pkt, function(p)
        return p.character_id == actor:id() and (match == nil or match(p))
    end, timeout or 5000)
end

local function face_expression(ctx)
    local actor = ctx:bot(0)
    local observer = ctx:bot(1)
    if setup(ctx, actor, 10, 0, PANIC .. ":1") == false then
        return false
    end
    if meet(ctx, actor, observer) == false then
        return false
    end

    if seen(actor, observer, resp.face_expression, req.face_expression { expression = 1 }, function(p)
            return p.expression == 1
        end) == false then
        return ctx:fail("기본 표정이 다른 플레이어에게 보이지 않음")
    end
    if seen(actor, observer, resp.face_expression, req.face_expression { expression = 8 }, nil, QUIET_MS) ~= false then
        return ctx:fail("표정 아이템 없이 캐시 표정이 보임")
    end
    if seen(actor, observer, resp.face_expression, req.face_expression { expression = 23 }, nil, QUIET_MS) ~= false then
        return ctx:fail("범위를 벗어난 표정이 보임")
    end
    if seen(actor, observer, resp.face_expression, req.face_expression { expression = PANIC - 5159992 }, function(p)
            return p.expression == PANIC - 5159992
        end) == false then
        return ctx:fail("가진 표정 아이템의 표정이 보이지 않음")
    end
    return true
end

local function item_effect(ctx)
    local actor = ctx:bot(0)
    local observer = ctx:bot(1)
    if setup(ctx, actor, 10, 0, SUN_EFFECT .. ":1") == false then
        return false
    end
    if meet(ctx, actor, observer) == false then
        return false
    end

    if seen(actor, observer, resp.item_effect, req.item_effect { item_id = MOON_EFFECT }, nil, QUIET_MS) ~= false then
        return ctx:fail("갖지 않은 이펙트 아이템이 켜짐")
    end
    if seen(actor, observer, resp.item_effect, req.item_effect { item_id = VOMIT }, nil, QUIET_MS) ~= false then
        return ctx:fail("이펙트 아이템이 아닌 아이템으로 이펙트가 켜짐")
    end
    if seen(actor, observer, resp.item_effect, req.item_effect { item_id = SUN_EFFECT }, function(p)
            return p.item_id == SUN_EFFECT
        end) == false then
        return ctx:fail("이펙트 아이템이 켜지지 않음")
    end
    if seen(actor, observer, resp.item_effect, req.item_effect { item_id = SUN_EFFECT }, nil, QUIET_MS) ~= false then
        return ctx:fail("같은 이펙트를 다시 켰는데 다시 알림")
    end
    if seen(actor, observer, resp.item_effect, req.item_effect { item_id = 0 }, function(p)
            return p.item_id == 0
        end) == false then
        return ctx:fail("이펙트가 꺼지지 않음")
    end
    return true
end

local function close_chalkboard(ctx)
    local actor = ctx:bot(0)
    local observer = ctx:bot(1)
    if meet(ctx, actor, observer) == false then
        return false
    end
    if seen(actor, observer, resp.chalkboard, req.close_chalkboard {}, nil, QUIET_MS) ~= false then
        return ctx:fail("칠판이 없는데 칠판 닫기를 알림")
    end
    return true
end

local function charge_skill_effect(ctx)
    local actor = ctx:bot(0)
    local observer = ctx:bot(1)
    if setup(ctx, actor, 30, BRAWLER, "-") == false then
        return false
    end
    if meet(ctx, actor, observer) == false then
        return false
    end

    if seen(actor, observer, resp.charge_skill_effect, req.skill_effect { skill_id = CORKSCREW_BLOW, level = 1, speed = 4 }, nil, QUIET_MS) ~= false then
        return ctx:fail("배우지 않은 스킬의 충전 이펙트가 보임")
    end
    if pq.command(actor, "/스킬레벨 " .. CORKSCREW_BLOW .. " 1", "스킬레벨 설정") == false then
        return ctx:fail("코크스크류 블로우 배우기 실패")
    end
    if seen(actor, observer, resp.charge_skill_effect, req.skill_effect { skill_id = CORKSCREW_BLOW, level = 2, speed = 4 }, nil, QUIET_MS) ~= false then
        return ctx:fail("배운 레벨과 다른 충전 이펙트가 보임")
    end
    if seen(actor, observer, resp.charge_skill_effect, req.skill_effect { skill_id = CORKSCREW_BLOW, level = 1, flags = 2, speed = 4 }, function(
            p)
            return p.skill_id == CORKSCREW_BLOW and p.level == 1 and p.flags == 2 and p.speed == 4
        end) == false then
        return ctx:fail("충전 이펙트가 다른 플레이어에게 보이지 않음")
    end
    if seen(actor, observer, resp.cancel_charge_skill_effect, req.cancel_buff { source_id = CORKSCREW_BLOW }, function(p)
            return p.skill_id == CORKSCREW_BLOW
        end) == false then
        return ctx:fail("충전 취소가 다른 플레이어에게 보이지 않음")
    end

    if pq.command(actor, "/스킬레벨 " .. KNUCKLE_BOOSTER .. " 1", "스킬레벨 설정") == false then
        return ctx:fail("너클 부스터 배우기 실패")
    end
    if seen(actor, observer, resp.cancel_charge_skill_effect, req.cancel_buff { source_id = KNUCKLE_BOOSTER }, nil, QUIET_MS) ~= false then
        return ctx:fail("충전 스킬이 아닌 버프 취소에 충전 취소가 보임")
    end
    return true
end

test_suite {
    name = "Presentation: 표정·아이템 이펙트·칠판 닫기·충전 이펙트·충전 취소",
    bot_count = 2,

    scenarios = {
        face_expression,
        item_effect,
        close_chalkboard,
        charge_skill_effect,
    },
}
