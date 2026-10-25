local pq = require("script/integration/lib/party_quest")

local EQUIP = 1
local CASH_WEAPON = 1702000
local CASH_WEAPON_SLOT = -111

local function check(ctx, ok, message)
    if ok == false then
        return ctx:fail(message)
    end
    return true
end

-- The character list shows the stats the character was saved with and its cash weapon in the weapon sticker field.
local function saved_overview(ctx)
    local bot = ctx:bot(0)
    if pq.command(bot, "/봇초기화 30 100 0 " .. CASH_WEAPON .. ":1 -", "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    local commands = {
        { "/힘바꾸기 41", "힘 설정" },
        { "/덱스바꾸기 32", "덱스 설정" },
        { "/인트바꾸기 23", "인트 설정" },
        { "/럭바꾸기 14", "럭 설정" },
        { "/체력바꾸기 1234", "체력 설정" },
        { "/마력바꾸기 567", "마력 설정" },
        { "/AP바꾸기 7", "AP 설정" },
        { "/스킬포인트 5", "스킬포인트 설정" },
        { "/인기도 9", "인기도 설정" },
        { "/경험치얻기 100", "경험치 100 획득" },
    }
    for _, c in ipairs(commands) do
        if pq.command(bot, c[1], c[2]) == false then
            return ctx:fail(c[1] .. " 실패")
        end
    end
    local wear = bot:request(resp.inventory_operation,
        req.move_item { inventory_type = EQUIP, source = bot:slot(CASH_WEAPON) or 0, dest = CASH_WEAPON_SLOT, count = 1 },
        nil, 3000)
    if check(ctx, wear ~= false, "캐시 무기 착용 응답 없음") == false then
        return false
    end

    local r = bot:relog({ passwords = { "" } })
    if r == false then
        return ctx:fail("재접속 실패")
    end
    local c = r.listed
    if check(ctx, c ~= nil, "캐릭터 목록에 캐릭터가 없음") == false then
        return false
    end
    local got = string.format("lv=%d str=%d dex=%d int=%d luk=%d hp=%d/%d mp=%d/%d ap=%d sp=%d exp=%d pop=%d weapon=%d",
        c.level, c.str, c.dex, c.int, c.luk, c.hp, c.max_hp, c.mp, c.max_mp, c.ability_point, c.skill_point, c.exp,
        c.population, c.weapon)
    return check(ctx,
        c.level == 30 and c.str == 41 and c.dex == 32 and c.int == 23 and c.luk == 14 and c.hp == 1234 and
        c.max_hp == 1234 and c.mp == 567 and c.max_mp == 567 and c.ability_point == 7 and c.skill_point == 5 and
        c.exp == 100 and c.population == 9 and c.weapon == CASH_WEAPON,
        "캐릭터 목록 값이 저장한 값과 다름: " .. got)
end

test_suite {
    name = "CharacterList: 능력치·캐시 무기",
    bot_count = 1,

    scenarios = {
        saved_overview,
    },
}
