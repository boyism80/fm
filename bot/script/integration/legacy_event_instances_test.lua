local pq = require("script/integration/lib/party_quest")

local TIMEOUT_SECONDS = 3

local function start(ctx, bot, group, map)
    local warp = bot:request(resp.warp, req.normal_chat { message = "/SM시작 " .. group }, function(p)
        return p.character.map == map
    end, 10000)
    if warp == false then
        return ctx:fail(group .. " 입장 실패")
    end
    return true
end

local function find_mob(ctx, bot, mob_id, except)
    for _ = 1, 25 do
        for _, mob in ipairs(bot:mobs(mob_id)) do
            if mob.oid ~= except then
                return mob
            end
        end
        ctx:sleep(200)
    end
    return nil
end

local function timeout(ctx, bot, exit_map)
    if pq.command(bot, "/타이머 " .. TIMEOUT_SECONDS, "타이머 제한: " .. TIMEOUT_SECONDS) == false then
        return ctx:fail("타이머 제한 설정 실패")
    end
    return pq.wait_map(ctx, bot, exit_map, TIMEOUT_SECONDS * 1000 + 15000)
end

local function boss(group, map, exit_map, bosses, cleared_npc)
    return function(ctx)
        local bot = ctx:bot(0)
        if pq.command(bot, "/봇초기화 200 112 0 - -", "봇초기화 완료") == false then
            return ctx:fail("봇 초기화 실패")
        end
        if start(ctx, bot, group, map) == false then
            return false
        end
        local previous = nil
        for phase, mob_id in ipairs(bosses) do
            local mob = find_mob(ctx, bot, mob_id, previous)
            if mob == nil then
                return ctx:fail(string.format("%s %d번째 보스가 나오지 않음", group, phase))
            end
            if bot:kill(mob.oid) == false then
                return ctx:fail(string.format("%s %d번째 보스 처치 실패", group, phase))
            end
            previous = mob.oid
        end
        if cleared_npc ~= nil then
            local found = false
            for _ = 1, 25 do
                if bot:npc(cleared_npc) ~= nil then
                    found = true
                    break
                end
                ctx:sleep(200)
            end
            if not found then
                return ctx:fail(group .. " 처치 뒤 NPC가 나오지 않음")
            end
        end
        return timeout(ctx, bot, exit_map)
    end
end

local function rusa(group, map, exit_map)
    return function(ctx)
        local bot = ctx:bot(0)
        if start(ctx, bot, group, map) == false then
            return false
        end
        return timeout(ctx, bot, exit_map)
    end
end

test_suite {
    name = "Legacy event instances: 쇼와·베르가모트·월드 보스·블러드 루사",
    bot_count = 2,

    on_initialize = function(ctx)
        for i = 0, ctx:bot_count() - 1 do
            if pq.command(ctx:bot(i), "/플레이어모드", "플레이어 모드: enabled") == false then
                return ctx:fail("플레이어 모드 설정 실패")
            end
        end
        return true
    end,

    on_finished = function(ctx)
        pq.command(ctx:bot(0), "/타이머 0", "타이머 제한: 0")
    end,

    scenarios = {
        function(ctx)
            local bot, other = ctx:bot(0), ctx:bot(1)
            if start(ctx, bot, "showa_boss", 801040100) == false then
                return false
            end
            if pq.command(other, "/SM시작 showa_boss", "SM 시작 실패") == false then
                return ctx:fail("쇼와 은신처 보스에 두 번째 그룹이 들어감")
            end
            return timeout(ctx, bot, 801040004)
        end,
        boss("showa_boss", 801040100, 801040004, { 9400112, 9400112 }),
        boss("vergamot", 802000211, 802000212, { 9400276, 9400263 }, 9120026),
        boss("world_boss", 123456785, 123456789, { 5000007, 5000007 }),
        rusa("blood_rusa_1", 100000208, 100000200),
        rusa("blood_rusa_2", 100000209, 100000204),
        rusa("blood_rusa_3", 100000210, 100000200),
    },
}
