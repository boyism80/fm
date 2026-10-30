local pq = require("script/integration/lib/party_quest")

local HUNTING_GROUND = 104040000
-- A boss whose WZ buff is 2022449 (Pink Bean expedition victory: weapon attack 35) and that has no die script.
local BUFF_BOSS = 5000007
local VICTORY = 2022449
local WEAPON_ATK = { mask = 0x1, position = 4 }

local received = {}

local function victory_attack(p)
    if p.buff_id ~= VICTORY and p.buff_id ~= -VICTORY then
        return nil
    end
    for _, entry in ipairs(p.buffs or {}) do
        if entry.buff.mask == WEAPON_ATK.mask and entry.buff.position == WEAPON_ATK.position then
            return entry.value
        end
    end
    return nil
end

-- A mob with a WZ buff puts that item's effect on every living character on the map when it dies.
local function victory_buff(ctx)
    local killer = ctx:bot(0)
    local other = ctx:bot(1)
    for _, bot in ipairs({ killer, other }) do
        if pq.command(bot, "/봇초기화 70 100 0 -", "봇초기화 완료") == false then
            return ctx:fail(bot:name() .. " 봇 초기화 실패")
        end
        if bot:instance_move(HUNTING_GROUND) == false then
            return ctx:fail(bot:name() .. " 사냥터 이동 실패")
        end
    end
    if killer:request(resp.spawn_mob, req.normal_chat { message = "/몬스터생성 " .. BUFF_BOSS }, function(p)
            return p.mob.mob_id == BUFF_BOSS
        end) == false then
        return ctx:fail("버프 보스 생성 응답 없음")
    end

    received = {}
    ctx:hook("mob_buff", function(_, bot, p, name)
        if name == "update_self_buff" and victory_attack(p) ~= nil then
            received[bot:id()] = victory_attack(p)
        end
    end)
    if pq.command(killer, "/몬스터죽이기", "마리 제거") == false then
        return ctx:fail("몬스터 처치 명령 실패")
    end
    for _ = 1, 30 do
        if received[killer:id()] ~= nil and received[other:id()] ~= nil then
            break
        end
        ctx:sleep(100)
    end
    ctx:unhook("mob_buff")
    for _, bot in ipairs({ killer, other }) do
        if received[bot:id()] ~= 35 then
            return ctx:fail(bot:name() .. " 처치 버프(물공 35)를 받지 않음: " .. tostring(received[bot:id()]))
        end
    end
    return true
end

test_suite {
    name = "MobBuff: 처치 시 맵 전원 버프",
    bot_count = 2,

    scenarios = {
        victory_buff,
    },
}
