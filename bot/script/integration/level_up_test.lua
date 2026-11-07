local pq = require("script/integration/lib/party_quest")

local LEVEL = 0x10
local EXP = 0x10000
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
local HP = 0x400
local HIT_ENV = -3
local FIELD = 100000000
local CLASS_CHANGE_EFFECT = 8
local CRUSADER = 111
local HERO = 112
local ACHILLES = 1120004

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

-- A Cygnus Knight starts level 11 with a tenth of its 2360 EXP already filled; a beginner starts it empty.
local function cygnus_level_exp(ctx)
    for _, case in ipairs({ { class = NOBLESSE, exp = 236 }, { class = 0, exp = 0 } }) do
        local p = level_up(ctx, ctx:bot(0), 10, case.class, 1716)
        if p == nil then
            return false
        end
        if p.stats[LEVEL] ~= 11 or p.stats[EXP] ~= case.exp then
            return ctx:fail(string.format("직업 %d가 레벨 %d, 경험치 %s로 시작함 (기대 11, %d)", case.class,
                p.stats[LEVEL], tostring(p.stats[EXP]), case.exp))
        end
    end
    return true
end

-- An explorer's first advancement adds 200~250 max HP and four slots to each tab, fills HP, and shows the others the
-- effect; a Cygnus advancement only fills HP.
local function class_change(ctx)
    local bot, observer = ctx:bot(0), ctx:bot(1)
    for _, case in ipairs({ { from = 0, to = WARRIOR, hp = { 200, 250 }, slots = 36 }, { from = NOBLESSE, to = 1100, hp = { 0, 0 } } }) do
        if pq.command(bot, "/봇초기화 10 " .. case.from .. " 0 - -", "봇초기화 완료") == false then
            return ctx:fail("봇 초기화 실패")
        end
        if pq.command(bot, "/체력바꾸기 " .. START_HP, "체력 설정") == false then
            return ctx:fail("체력 설정 실패")
        end
        if bot:instance_move(FIELD) == false or observer:instance_move(FIELD) == false then
            return ctx:fail("맵 이동 실패")
        end
        if bot:request(resp.update_stats, req.damaged { type = HIT_ENV, damage = 300 }, nil, 3000) == false then
            return ctx:fail("체력을 깎지 못함")
        end
        local slots, effect = {}, false
        ctx:hook("class_change", function(_, b, p, name)
            if b:id() == bot:id() and name == resp.inventory_slot_limit then
                slots[p.inventory_type] = p.limit
            elseif b:id() == observer:id() and name == resp.show_effect and p.character_id == bot:id() and p.type == CLASS_CHANGE_EFFECT then
                effect = true
            end
        end)
        local p = bot:request(resp.update_stats, req.normal_chat { message = "/직업바꾸기 " .. case.to }, function(p)
            return p.stats[MAX_HP] ~= nil
        end, 3000)
        if p == false then
            return ctx:fail("직업 " .. case.to .. " 전직에 최대 HP 갱신이 오지 않음")
        end
        local gained = p.stats[MAX_HP] - START_HP
        if gained < case.hp[1] or gained > case.hp[2] or p.stats[HP] ~= p.stats[MAX_HP] then
            return ctx:fail(string.format("직업 %d 전직: 최대 HP +%d, HP %d/%d", case.to, gained, p.stats[HP], p.stats[MAX_HP]))
        end
        ctx:sleep(500)
        for tab = 1, 4 do
            if slots[tab] ~= case.slots then
                return ctx:fail(string.format("직업 %d 전직: %d번 탭 칸 %s (기대 %s)", case.to, tab, tostring(slots[tab]),
                    tostring(case.slots)))
            end
        end
        if not effect then
            return ctx:fail("직업 " .. case.to .. " 전직 이펙트가 주변에 오지 않음")
        end
    end
    return true
end

-- The fourth advancement opens the class's fourth class skills at level 0 and their WZ master level.
local function fourth_class_skills(ctx)
    local bot = ctx:bot(0)
    if pq.command(bot, "/봇초기화 120 " .. CRUSADER .. " 0 - -", "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    local p = bot:request(resp.update_skills, req.normal_chat { message = "/직업바꾸기 " .. HERO }, function(p)
        return p.skill_id == ACHILLES
    end, 3000)
    if p == false then
        return ctx:fail("4차 전직 뒤 아킬레스가 스킬 목록에 오지 않음")
    end
    if p.level ~= 0 or p.master_level ~= 10 then
        return ctx:fail(string.format("아킬레스가 레벨 %d, 마스터 레벨 %d로 옴 (기대 0, 10)", p.level, p.master_level))
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
    name = "LevelUp: 초보 직업 SP·시그너스 경험치 보너스·직업별 HP 증가량·장비 INT의 MP 증가량·전직 보상",
    bot_count = 2,

    scenarios = {
        noblesse_sp,
        warrior_hp,
        equipped_int_mp,
        cygnus_level_exp,
        class_change,
        fourth_class_skills,
    },
}
