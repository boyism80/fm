local pq = require("script/integration/lib/party_quest")

local LEVEL = 0x10
local MAX_HP = 0x800
local MAX_MP = 0x2000
local AVAILABLE_SP = 0x8000
local WARRIOR = 100
local NOBLESSE = 1000
local START_HP = 1000
local START_MP = 1000
local EQUIP = 1
local CAP_SLOT = -1
local INT_CAP = 1002140

local function level_up(ctx, bot, level, class, exp, items)
    if pq.command(bot, "/봇초기화 " .. level .. " " .. class .. " 0 " .. (items or "-") .. " -", "봇초기화 완료") == false then
        ctx:fail("봇 초기화 실패")
        return nil
    end
    if pq.command(bot, "/체력바꾸기 " .. START_HP, "체력 설정") == false then
        ctx:fail("체력 설정 실패")
        return nil
    end
    if pq.command(bot, "/마력바꾸기 " .. START_MP, "마력 설정") == false then
        ctx:fail("마력 설정 실패")
        return nil
    end
    if items ~= nil then
        local worn = bot:request(resp.inventory_operation,
            req.move_item { inventory_type = EQUIP, source = bot:slot(INT_CAP) or 0, dest = CAP_SLOT, count = 1 }, nil,
            3000)
        if worn == false then
            ctx:fail("모자 착용 응답 없음")
            return nil
        end
    end
    local p = bot:request(resp.update_stats, req.normal_chat { message = "/경험치얻기 " .. exp }, function(p)
        return p.stats[LEVEL] ~= nil and p.stats[MAX_HP] ~= nil and p.stats[MAX_MP] ~= nil and
        p.stats[AVAILABLE_SP] ~= nil
    end, 3000)
    if p == false then
        ctx:fail("경험치를 얻었는데 레벨이 오르지 않음")
        return nil
    end
    return p
end

-- Noblesse is a beginner class, so its level ups give no SP.
local function noblesse_sp(ctx)
    local p = level_up(ctx, ctx:bot(0), 2, NOBLESSE, 1000)
    if p == nil then
        return false
    end
    if p.stats[AVAILABLE_SP] ~= nil and p.stats[AVAILABLE_SP] ~= 0 then
        return ctx:fail(string.format("노블레스가 레벨 %d에서 SP %d를 받음", p.stats[LEVEL], p.stats[AVAILABLE_SP]))
    end
    return true
end

-- A warrior's level up adds 24~28 HP a level, not a level-scaled amount.
local function warrior_hp(ctx)
    local from = 50
    local p = level_up(ctx, ctx:bot(0), from, WARRIOR, 3000000)
    if p == nil then
        return false
    end
    local levels = p.stats[LEVEL] - from
    local gained = p.stats[MAX_HP] - START_HP
    if gained < 24 * levels or gained > 28 * levels then
        return ctx:fail(string.format("전사 %d레벨 상승에 HP %d (기대 %d~%d)", levels, gained, 24 * levels, 28 * levels))
    end
    return true
end

-- The INT behind a level up's MP counts the equipped items, so an INT 999 cap adds about 100 MP a level.
local function equipped_int_mp(ctx)
    local from = 50
    local p = level_up(ctx, ctx:bot(0), from, WARRIOR, 3000000, INT_CAP .. ":1")
    if p == nil then
        return false
    end
    local levels = p.stats[LEVEL] - from
    local gained = p.stats[MAX_MP] - START_MP
    if gained < (4 + 99) * levels then
        return ctx:fail(string.format("INT 999 모자를 끼고 %d레벨 상승에 MP %d (기대 %d 이상)", levels, gained, (4 + 99) * levels))
    end
    return true
end

test_suite {
    name = "LevelUp: 초보 직업 SP·직업별 HP 증가량·장비 INT의 MP 증가량",
    bot_count = 1,

    scenarios = {
        noblesse_sp,
        warrior_hp,
        equipped_int_mp,
    },
}
