local pq = require("script/integration/lib/party_quest")

local WAITING_ROOM = 123256780
local BOSSES = {
    { npc = 2038, name = "힐라", map = 123356780, mob = 8870000, x = 164, y = 196 },
    { npc = 2022, name = "매그너스", map = 123356785, mob = 8880000, x = 2365, y = -1347 },
    { npc = 2023, name = "스우", map = 123356782, mob = 9801028, x = -110, y = -16 },
    { npc = 2024, name = "카오스 파풀라투스", map = 123356781, mob = 8820119, x = -393, y = -386 },
    { npc = 2025, name = "데미안", map = 123356786, mob = 9300890, x = 1222, y = 16 },
    { npc = 2026, name = "검은 마법사", map = 123356788, mob = 8880500, x = -3, y = 85 },
    { npc = 2039, name = "루시드", map = 450004150, mob = 8880166, x = 1008, y = 48 },
}

local function waiting_room(ctx, bot)
    if bot:map() == WAITING_ROOM then
        return true
    end
    local warp = bot:request(resp.warp, req.normal_chat { message = "!보스" }, function(p)
        return p.character.map == WAITING_ROOM
    end, 10000)
    if warp == false then
        return ctx:fail("!보스로 보스 대기실에 가지 못함")
    end
    return true
end

local function enter_and_kill(ctx, bot, boss)
    if waiting_room(ctx, bot) == false then
        return false
    end
    local oid = pq.npc(ctx, bot, boss.npc)
    if oid == false then
        return false
    end
    local dlg = bot:npc_click(oid)
    if dlg == false then
        return ctx:fail(boss.name .. " 입장 대화가 오지 않음")
    end
    local warp = bot:request(resp.warp, req.dialog { dialog_type = dlg.enable_escape and 11 or 12, next = true },
        function(p)
            return p.character.map == boss.map
        end, 10000)
    if warp == false then
        return ctx:fail(boss.name .. " 보스 맵으로 이동하지 않음")
    end
    if pq.move(bot, boss.x, boss.y) == false then
        return ctx:fail(boss.name .. " 보스 위치로 이동 실패")
    end
    local mob = bot:mobs(boss.mob)[1]
    if mob == nil then
        local spawn = bot:request(resp.spawn_mob, nil, function(p)
            return p.mob.mob_id == boss.mob
        end, 5000)
        if spawn == false then
            return ctx:fail(boss.name .. " 보스가 소환되지 않음: " .. boss.mob)
        end
        mob = { oid = spawn.mob.oid }
    end
    if bot:kill(mob.oid) == false then
        return ctx:fail(boss.name .. " 처치 실패: " .. boss.mob)
    end
    if bot:mobs(boss.mob)[1] ~= nil then
        return ctx:fail(boss.name .. " 처치 후에도 보스가 남아 있음")
    end
    return true
end

-- Suu is fought through three maps: each phase's death sends the party on to the next 5 seconds later, with the
-- next phase's boss there and a new clock, and the world hears of each.
local function suu_phases(ctx)
    local bot = ctx:bot(0)
    local phases = {
        { map = 123356782, mob = 9801028 },
        { map = 123356783, mob = 9801029 },
        { map = 123356784, mob = 9801030 },
    }
    local heard = {}
    ctx:hook("suu_notices", function(_, b, p, name)
        if b:id() == bot:id() and name == resp.notice and p.message:find("[헤이븐]", 1, true) ~= nil then
            heard[#heard + 1] = p.message
        end
    end)
    local function finish(ok)
        ctx:unhook("suu_notices")
        return ok
    end
    if enter_and_kill(ctx, bot, BOSSES[3]) == false then
        return finish(false)
    end
    for i = 2, #phases do
        local phase = phases[i]
        if bot:map() ~= phase.map and bot:request(resp.warp, nil, function(p)
                return p.character.map == phase.map
            end, 15000) == false then
            return finish(ctx:fail("스우 " .. (i - 1) .. "페이즈를 잡아도 다음 페이즈로 넘어가지 않음"))
        end
        local mob = bot:mobs(phase.mob)[1]
        if mob == nil then
            local spawn = bot:request(resp.spawn_mob, nil, function(p)
                return p.mob.mob_id == phase.mob
            end, 5000)
            if spawn == false then
                return finish(ctx:fail("스우 " .. i .. "페이즈 보스가 없음: " .. phase.mob))
            end
            mob = { oid = spawn.mob.oid }
        end
        if bot:kill(mob.oid) == false then
            return finish(ctx:fail("스우 " .. i .. "페이즈 처치 실패"))
        end
    end
    ctx:sleep(1000)
    local notices = #heard
    if waiting_room(ctx, bot) == false then
        return finish(false)
    end
    if notices ~= 3 then
        return finish(ctx:fail("스우 페이즈 공지가 3번 오지 않음: " .. table.concat(heard, " / ")))
    end
    return finish(true)
end

test_suite {
    name = "Boss: 보스 처치 7종·스우 페이즈 이동",
    bot_count = 1,

    on_initialize = function(ctx)
        local bot = ctx:bot(0)
        if pq.command(bot, "/봇초기화 120 100 0", "봇초기화 완료") == false then
            return ctx:fail("봇 초기화 실패")
        end
        if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
            return ctx:fail("플레이어 모드 설정 실패")
        end
        if pq.command(bot, "/무적", "무적 상태: enabled") == false then
            return ctx:fail("무적 설정 실패")
        end
        if pq.command(bot, "/즉사 on", "즉사 상태: enabled") == false then
            return ctx:fail("즉사 설정 실패")
        end
        return waiting_room(ctx, bot)
    end,

    on_finished = function(ctx)
        ctx:bot(0):request(resp.party_update_disband, req.party_operation { operation = PARTY.Leave }, nil, 5000)
    end,

    scenarios = {
        pq.form_party,
        function(ctx)
            local passed = true
            for _, boss in ipairs(BOSSES) do
                if enter_and_kill(ctx, ctx:bot(0), boss) == false then
                    passed = false
                end
            end
            if waiting_room(ctx, ctx:bot(0)) == false then
                return false
            end
            return passed
        end,
        suu_phases,
    },
}
