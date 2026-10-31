local pq = require("script/integration/lib/party_quest")

local GUILD_HQ = 200000301
local HENESYS = 100000000
local GUILD_SHOW_INFO = 0x1A
local GUILD_MESO = 2000000
local LEVEL = 0x10
local EXP = 2000000
local WATCHED = {
    guild_member_level_class_update = true,
    family_info = true,
    family_chart = true,
}

local seen = {}

local function check(ctx, ok, message)
    if ok == false then
        return ctx:fail(message)
    end
    return true
end

local function watch(ctx)
    seen = {}
    ctx:hook("member_level", function(_, bot, p, name)
        if WATCHED[name] == nil then
            return
        end
        local key = bot:id() .. name
        seen[key] = seen[key] or {}
        table.insert(seen[key], p)
    end)
end

local function wait_for(ctx, bot, name, match, what)
    for _ = 1, 20 do
        for _, p in ipairs(seen[bot:id() .. name] or {}) do
            if match(p) then
                return p
            end
        end
        ctx:sleep(500)
    end
    return ctx:fail(string.format("%s %s 없음", bot:name(), what))
end

local function setup(ctx, bot, level, map, meso)
    if pq.command(bot, string.format("/봇초기화 %d 0 %d - -", level, meso or 0), "봇초기화 완료") == false then
        return ctx:fail(bot:name() .. " 봇 초기화 실패")
    end
    if bot:instance_move(map) == false then
        return ctx:fail(bot:name() .. " 맵 이동 실패")
    end
    return true
end

local function level_up(ctx, bot)
    local p = bot:request(resp.update_stats, req.normal_chat { message = "/경험치얻기 " .. EXP }, function(p)
        return p.stats[LEVEL] ~= nil
    end, 3000)
    if p == false then
        ctx:fail(bot:name() .. " 레벨업 응답 없음")
        return nil
    end
    return p.stats[LEVEL]
end

local function show_info(p)
    return p.code == GUILD_SHOW_INFO
end

local function guild_level(ctx)
    local master, member = ctx:bot(0), ctx:bot(1)
    watch(ctx)
    if setup(ctx, master, 50, GUILD_HQ, GUILD_MESO) == false or setup(ctx, member, 50, GUILD_HQ) == false then
        return false
    end

    if master:request(resp.guild_message, req.guild_operation { operation = GUILD.Create, guild_name = master:name() }, show_info) == false then
        return ctx:fail("길드 생성 실패")
    end
    local invite = master:request_on(member, resp.guild_invite,
        req.guild_operation { operation = GUILD.Invite, target_name = member:name() })
    if invite == false then
        return ctx:fail("길드 초대 패킷 없음")
    end
    local joined = member:request(resp.guild_message, req.guild_operation {
        operation = GUILD.AcceptInvite,
        guild_id = invite.guild_id,
        character_id = member:id(),
    }, show_info)
    if joined == false then
        return ctx:fail("길드 가입 실패")
    end

    if member:transfer(1) == false then
        return ctx:fail("길드원 채널 이동 실패")
    end
    ctx:sleep(2000)
    local level = level_up(ctx, member)
    if level == nil then
        return false
    end
    local p = wait_for(ctx, master, "guild_member_level_class_update", function(p)
        return p.character_id == member:id() and p.level == level
    end, "다른 채널 길드원 레벨 갱신")
    if member:transfer(0) == false then
        return ctx:fail("길드원 원래 채널 복귀 실패")
    end
    return p ~= false
end

local function family_level(ctx)
    local senior, junior = ctx:bot(0), ctx:bot(1)
    watch(ctx)
    if setup(ctx, senior, 50, HENESYS) == false or setup(ctx, junior, 40, HENESYS) == false then
        return false
    end
    local invite = senior:request_on(junior, resp.family_invite, req.family_invite { name = junior:name() })
    if invite == false then
        return ctx:fail("패밀리 초대가 오지 않음")
    end
    if junior:request(resp.family_joined, req.family_accept { senior_id = senior:id(), name = senior:name(), accepted = true }) == false then
        return ctx:fail("패밀리 가입 실패")
    end
    ctx:sleep(1000)

    seen = {}
    local level = level_up(ctx, senior)
    if level == nil then
        return false
    end
    if wait_for(ctx, junior, "family_info", function()
            return true
        end, "선배 레벨업 뒤 가족 정보") == false then
        return false
    end
    local chart = junior:request(resp.family_chart, req.family_pedigree { name = junior:name() })
    if chart == false then
        return ctx:fail("가계도 응답 없음")
    end
    local shown = nil
    for _, m in ipairs(chart.members) do
        if m.character_id == senior:id() then
            shown = m.level
        end
    end
    return check(ctx, shown == level, string.format("가계도의 선배 레벨 %s (기대 %d)", tostring(shown), level))
end

test_suite {
    name = "Member Level: 길드·가족 레벨 갱신",
    bot_count = 2,

    scenarios = {
        guild_level,
        family_level,
    },
}
