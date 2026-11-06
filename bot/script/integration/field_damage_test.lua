local pq = require("script/integration/lib/party_quest")

local EQUIP = 1
local HP = 0x400
local COLD_MAP = 211040000
local COLD_DAMAGE = 10
local COLD_INTERVAL_MS = 10000
local COLD_CAPE = 1102109
local CAPE_SLOT = -9
local START_HP = 1000

-- A map's WZ decHP takes HP every decInterval, and its protect item keeps it off.
local function cold_field(ctx)
    local bot = ctx:bot(0)
    if pq.command(bot, "/봇초기화 50 100 0 " .. COLD_CAPE .. ":1", "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    if pq.command(bot, "/체력바꾸기 " .. START_HP, "체력 설정") == false then
        return ctx:fail("체력 설정 실패")
    end
    if bot:force_instance_move(COLD_MAP, "추위로 HP가 주는 엘나스 산맥") == false then
        return ctx:fail("엘나스 산맥 이동 실패")
    end
    local hurt = bot:request(resp.update_stats, nil, function(p)
        return p.stats[HP] ~= nil and p.stats[HP] < START_HP
    end, COLD_INTERVAL_MS + 3000)
    if hurt == false then
        return ctx:fail("보호 장비 없이 있어도 HP가 줄지 않음")
    end
    if hurt.stats[HP] ~= START_HP - COLD_DAMAGE then
        return ctx:fail(string.format("환경 피해가 다름: %d (기대 %d)", START_HP - hurt.stats[HP], COLD_DAMAGE))
    end

    local wear = req.move_item { inventory_type = EQUIP, source = bot:slot(COLD_CAPE) or 0, dest = CAPE_SLOT, count = 1 }
    if bot:request(resp.inventory_operation, wear, nil, 3000) == false then
        return ctx:fail("보호 망토 착용 실패")
    end
    if bot:request(resp.update_stats, nil, function(p)
            return p.stats[HP] ~= nil and p.stats[HP] < hurt.stats[HP]
        end, COLD_INTERVAL_MS + 3000) ~= false then
        return ctx:fail("보호 망토를 입었는데 HP가 줄어듦")
    end
    return true
end

test_suite {
    name = "FieldDamage: 맵 환경 피해·보호 장비",
    bot_count = 1,

    scenarios = {
        cold_field,
    },
}
