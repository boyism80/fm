local KERNING_PQ_STAGE = 103000800
local ROMEO_STAGE = 926100000
local NEO_CITY = 240070100
local INSTANCE_KEY = 91

local function check(ctx, ok, message)
    if ok == false then
        return ctx:fail(message)
    end
    return true
end

-- The first user to enter an instance gets the stage hint the map's WZ names.
local function stage_hint(ctx)
    local bot = ctx:bot(0)
    local hint = bot:request(resp.weather,
        req.normal_chat { message = string.format("/인스턴스이동 %d %d", KERNING_PQ_STAGE, INSTANCE_KEY) }, function(p)
        return p.message:find("쿠폰", 1, true) ~= nil
    end, 5000)
    if check(ctx, hint ~= false, "첫 입장에 단계 안내 문구가 오지 않음") == false then
        return false
    end

    hint = bot:request(resp.weather,
        req.normal_chat { message = string.format("/인스턴스이동 %d %d", ROMEO_STAGE, INSTANCE_KEY) }, function(p)
        return p.message:find("스위치", 1, true) ~= nil
    end, 5000)
    return check(ctx, hint ~= false, "로미오 단계 안내 문구가 오지 않음")
end

-- onUserEnter names that send a title effect to the one who enters.
local function title_effect(ctx)
    local bot = ctx:bot(0)
    local effect = bot:request(resp.environment_change,
        req.normal_chat { message = string.format("/인스턴스이동 %d %d", NEO_CITY, INSTANCE_KEY) }, function(p)
        return p.env == "temaD/enter/neoCity1"
    end, 5000)
    return check(ctx, effect ~= false, "입장 제목 효과가 오지 않음")
end

test_suite {
    name = "MapEnter: WZ가 이름 붙인 맵 입장 스크립트 (단계 안내·제목 효과)",
    bot_count = 1,

    scenarios = {
        stage_hint,
        title_effect,
    },
}
