local pq = require("script/integration/lib/party_quest")
local route = require("script/integration/lib/route")

local DIALOG_DEFAULT = 0
local DIALOG_YES_NO = 1

local function click_next(ctx, bot, npc, map_id)
    local oid = pq.npc(ctx, bot, npc)
    if oid == false then
        return false
    end
    if bot:npc_click(oid) == false then
        return ctx:fail(bot:name() .. " 대화가 오지 않음: " .. npc)
    end
    local warp = bot:request(resp.warp, req.dialog { dialog_type = DIALOG_DEFAULT, next = true }, function(p)
        return p.character.map == map_id
    end, 10000)
    if warp == false then
        return ctx:fail(bot:name() .. " 맵 이동 실패: " .. map_id)
    end
    return true
end

local function click_warp(ctx, bot, npc, map_id)
    local oid = pq.npc(ctx, bot, npc)
    if oid == false then
        return false
    end
    local warp = bot:request(resp.warp, req.npc_click { oid = oid }, function(p)
        return p.character.map == map_id
    end, 10000)
    if warp == false then
        return ctx:fail(bot:name() .. " 맵 이동 실패: " .. map_id)
    end
    return true
end

local function confirm_warp(ctx, bot, npc, map_id)
    local oid = pq.npc(ctx, bot, npc)
    if oid == false then
        return false
    end
    if bot:npc_click(oid) == false then
        return ctx:fail(bot:name() .. " 대화가 오지 않음: " .. npc)
    end
    local warp = bot:request(resp.warp, req.dialog { dialog_type = DIALOG_YES_NO, next = true }, function(p)
        return p.character.map == map_id
    end, 10000)
    if warp == false then
        return ctx:fail(bot:name() .. " 맵 이동 실패: " .. map_id)
    end
    return true
end

local function kill_boss(ctx, bot, mob_id)
    local mob = bot:mobs(mob_id)[1]
    if mob == nil then
        local spawn = bot:request(resp.spawn_mob, nil, function(p)
            return p.mob.mob_id == mob_id
        end, 10000)
        if spawn == false then
            return ctx:fail("보스가 보이지 않음: " .. mob_id)
        end
        mob = { oid = spawn.mob.oid }
    end
    if bot:kill(mob.oid) == false then
        return ctx:fail("보스 처치 실패: " .. mob_id)
    end
    return true
end

-- A field boss comes back mobTime after it died, however long it lived first.
local function respawn_after_death(ctx)
    local bot = ctx:bot(0)
    local mushmom, field = 6130101, 100000005
    if bot:instance_move(field) == false then
        return ctx:fail("머쉬맘 맵 이동 실패")
    end
    if bot:mobs(mushmom)[1] == nil and bot:request(resp.spawn_mob, nil, function(p)
            return p.mob.mob_id == mushmom
        end, 15000) == false then
        return ctx:fail("머쉬맘이 나오지 않음")
    end
    -- The boss has now lived past its mobTime.
    if pq.command(bot, "/시간가속 00:31:00", "시간가속 적용 요청") == false then
        return ctx:fail("시간가속 실패")
    end
    ctx:sleep(2000)
    local killed = kill_boss(ctx, bot, mushmom)
    local respawned = false
    if killed ~= false then
        respawned = bot:request(resp.spawn_mob, nil, function(p)
            return p.mob.mob_id == mushmom
        end, 12000) ~= false
    end
    pq.command(bot, "/현재시간초기화", "현재 시간 보정 초기화")
    ctx:sleep(2000)
    if killed == false then
        return false
    end
    if respawned then
        return ctx:fail("오래 살아 있던 머쉬맘이 죽은 뒤 mobTime 안에 다시 나옴")
    end
    return true
end

local function hit(bot, oid)
    local x, y = bot:position()
    return req.attack {
        close_attack_info = {
            attack_header = { targets = 1, hits = 1 },
            damages = { { oid = oid, damage_pairs = { { damage = 1 } } } },
            position = { x = x, y = y },
        },
    }
end

-- A boss with an HP tag shows the boss HP bar when hit and clears it when it dies; a boss without one shows no HP.
local function boss_hp_bar(ctx, bot, oid)
    local mano, zakum_arm = 2220000, 8800003
    if pq.command(bot, "/즉사 off", "즉사 상태: disabled") == false then
        return ctx:fail("즉사 끄기 실패")
    end
    -- Starting stats miss Mano, and a missed hit is rejected before any HP shows.
    if pq.command(bot, "/전능 999", "전능 설정") == false then
        return ctx:fail("스탯 설정 실패")
    end
    local bar = bot:request(resp.show_boss_hp, hit(bot, oid), function(p)
        return p.mob_id == mano
    end, 3000)
    if bar == false or bar.current_hp ~= bar.max_hp - 1 then
        return ctx:fail("마노를 때렸는데 보스 HP 바가 오지 않음")
    end
    local cleared = false
    ctx:hook("mano_hp_bar", function(_, b, p, name)
        if b:id() == bot:id() and name == resp.show_boss_hp and p.mob_id == mano and p.current_hp == -1 then
            cleared = true
        end
    end)
    if bot:kill(oid) == false then
        return ctx:fail("마노 처치 실패")
    end
    if cleared == false then
        return ctx:fail("마노가 죽었는데 보스 HP 바가 지워지지 않음")
    end

    if pq.command(bot, "/즉사 off", "즉사 상태: disabled") == false then
        return ctx:fail("즉사 끄기 실패")
    end
    local arm = bot:request(resp.spawn_mob, req.normal_chat { message = "/몬스터생성 " .. zakum_arm }, function(p)
        return p.mob.mob_id == zakum_arm
    end, 5000)
    if arm == false then
        return ctx:fail("자쿰 팔 생성 응답 없음")
    end
    if bot:request({ resp.show_mob_hp, resp.show_boss_hp }, hit(bot, arm.mob.oid), nil, 1500) ~= false then
        return ctx:fail("태그 없는 자쿰 팔의 HP가 보임")
    end
    return true
end

-- Mano is not in its map's WZ life; it comes from the field boss table with a notice to the map.
local function mano_from_table(ctx)
    local bot = ctx:bot(0)
    local mano, field = 2220000, 104000400
    if bot:instance_move(field) == false then
        return ctx:fail("마노 맵 이동 실패")
    end
    if kill_boss(ctx, bot, mano) == false then
        return false
    end
    -- A field boss belongs to no carnival team.
    local team = nil
    ctx:hook("mano_team", function(_, b, p, name)
        if b:id() == bot:id() and name == resp.spawn_mob and p.mob.mob_id == mano then
            team = p.mob.carnival_team
        end
    end)
    if pq.command(bot, "/시간가속 00:09:00", "시간가속 적용 요청") == false then
        return ctx:fail("시간가속 실패")
    end
    local notice = bot:request(resp.notice, nil, function(p)
        return p.message == "서늘한 기운이 감돌며 마노가 나타났습니다."
    end, 15000)
    pq.command(bot, "/현재시간초기화", "현재 시간 보정 초기화")
    ctx:sleep(2000)
    if notice == false then
        return ctx:fail("마노가 다시 나오며 공지가 오지 않음")
    end
    if team ~= -1 then
        return ctx:fail("다시 나온 마노의 카니발 팀이 다름: " .. tostring(team))
    end
    local respawned = bot:mobs(mano)[1]
    if respawned == nil then
        return ctx:fail("공지 뒤에 마노가 보이지 않음")
    end
    return boss_hp_bar(ctx, bot, respawned.oid)
end

test_suite {
    name = "Boss: 대왕지네·포장마차·필드 보스 재생성",
    bot_count = 1,

    on_initialize = function(ctx)
        local bot = ctx:bot(0)
        if pq.command(bot, "/봇초기화 120 100 0 - 4103:1,4014:2,4013:1", "봇초기화 완료") == false then
            return ctx:fail("봇 초기화 실패")
        end
        if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
            return ctx:fail("플레이어 모드 설정 실패")
        end
        return true
    end,

    scenarios = {
        function(ctx)
            local bot = ctx:bot(0)
            if route.travel_instance(ctx, bot, 701010321) == false then
                return ctx:fail("검은 숲으로 이동 실패")
            end
            if click_warp(ctx, bot, 9310005, 701010322) == false then
                return false
            end
            if click_next(ctx, bot, 9310006, 701010323) == false then
                return false
            end
            if kill_boss(ctx, bot, 9600009) == false then
                return false
            end
            return confirm_warp(ctx, bot, 9310007, 701010320)
        end,
        function(ctx)
            local bot = ctx:bot(0)
            if route.travel_instance(ctx, bot, 741020100) == false then
                return ctx:fail("야시장 뒷골목으로 이동 실패")
            end
            if click_next(ctx, bot, 9330028, 741020101) == false then
                return false
            end
            if kill_boss(ctx, bot, 9410014) == false then
                return false
            end
            return confirm_warp(ctx, bot, 9330032, 741020100)
        end,
        respawn_after_death,
        mano_from_table,
    },
}
