local pq = require("script/integration/lib/party_quest")
local fault = require("script/integration/lib/fault")
local settle = require("script/integration/lib/settle")

local LITH_HARBOR = 104000000
local GRANDPA = 1002003
local FEE = 50000
local MESO = 100000
local DEFAULT = 20
local STEP = 5
local DIALOG_YES_NO = 1
local METHOD = "IncreaseBuddyCapacity"

local function check(ctx, ok, message)
    if ok == false then
        return ctx:fail(message)
    end
    return true
end

local function setup(ctx, bot)
    if pq.command(bot, "/봇초기화 30 0 " .. MESO .. " -", "봇초기화 완료") == false then
        return ctx:fail(bot:name() .. " 봇 초기화 실패")
    end
    if bot:instance_move(LITH_HARBOR) == false then
        return ctx:fail(bot:name() .. " 리스항구 이동 실패")
    end
    return true
end

-- Answers yes twice and returns the grandpa's last line.
local function expand(ctx, bot)
    local oid = pq.npc(ctx, bot, GRANDPA)
    if oid == false then
        return false
    end
    if bot:npc_click(oid) == false then
        return ctx:fail("몽롱 할배 대화가 오지 않음")
    end
    if bot:request(resp.dialog_yes_no, req.dialog { dialog_type = DIALOG_YES_NO, next = true }, nil, 5000) == false then
        return ctx:fail("가격 안내가 오지 않음")
    end
    local last = bot:request(resp.dialog, req.dialog { dialog_type = DIALOG_YES_NO, next = true }, nil, 10000)
    bot:dialog(false)
    if last == false then
        return ctx:fail("확장 결과 대사가 오지 않음")
    end
    return last.text
end

local function holds(ctx, bot, capacity, meso, what)
    if check(ctx, bot:buddy_capacity() == capacity,
            string.format("%s 친구 목록 %d (기대 %d)", what, bot:buddy_capacity(), capacity)) == false then
        return false
    end
    return check(ctx, bot:meso() == meso, string.format("%s 메소 %d (기대 %d)", what, bot:meso(), meso))
end

local function relog(ctx, bot)
    if bot:relog({ passwords = { "" } }) == false then
        return ctx:fail("재접속 실패")
    end
    return true
end

-- The fault is server wide, so it is turned off right after the one request it targets; the timer stays paused.
local function expand_unknown(ctx, bot, mode)
    if settle.pause(bot) == false then
        return ctx:fail("정산 타이머 정지 실패")
    end
    if fault.rpc(bot, METHOD, mode) == false then
        return ctx:fail(METHOD .. " 장애 주입 실패")
    end
    local text = expand(ctx, bot)
    fault.rpc(bot, METHOD, "off")
    if text == false then
        return false
    end
    return check(ctx, text:find("확실히 가지고 있는 건가", 1, true) ~= nil, "결과를 모르는데 확장 성공 대사가 나옴: " .. text)
end

-- The raise is kept by internal, so it survives a relogin.
local function expanded(ctx)
    local bot = ctx:bot(0)
    if setup(ctx, bot) == false then
        return false
    end
    local text = expand(ctx, bot)
    if text == false then
        return false
    end
    if check(ctx, text:find("5개 늘어나", 1, true) ~= nil, "확장 성공 대사가 아님: " .. text) == false then
        return false
    end
    ctx:sleep(500)
    if holds(ctx, bot, DEFAULT + STEP, MESO - FEE, "확장 직후") == false then
        return false
    end
    if relog(ctx, bot) == false then
        return false
    end
    return holds(ctx, bot, DEFAULT + STEP, MESO - FEE, "재접속 뒤")
end

-- A reply lost after internal raised it: the fee stays held, and the timer's resend applies the same raise once.
local function lost_timer(ctx)
    local bot = ctx:bot(1)
    if setup(ctx, bot) == false then
        return false
    end
    if expand_unknown(ctx, bot, "lost") == false then
        return false
    end
    if holds(ctx, bot, DEFAULT, MESO - FEE, "응답 유실 직후 보류") == false then
        return false
    end
    if settle.resume(bot) == false then
        return ctx:fail("정산 타이머 재개 실패")
    end
    if settle.wait(ctx, function()
            return bot:buddy_capacity() == DEFAULT + STEP
        end) == false then
        return ctx:fail("주기 정리 뒤 친구 목록이 늘지 않음: " .. bot:buddy_capacity())
    end
    settle.tick(ctx)
    if holds(ctx, bot, DEFAULT + STEP, MESO - FEE, "주기 정리 뒤 (중복 확장·환불 없음)") == false then
        return false
    end
    if relog(ctx, bot) == false then
        return false
    end
    return holds(ctx, bot, DEFAULT + STEP, MESO - FEE, "재접속 뒤")
end

-- A request that never reached internal survives the relogin and is applied there, once.
local function unsent_relogin(ctx)
    local bot = ctx:bot(1)
    if expand_unknown(ctx, bot, "on") == false then
        return false
    end
    if holds(ctx, bot, DEFAULT + STEP, MESO - 2 * FEE, "요청 미전송 직후 보류") == false then
        return false
    end
    if relog(ctx, bot) == false then
        return false
    end
    if settle.wait(ctx, function()
            return bot:buddy_capacity() == DEFAULT + 2 * STEP
        end) == false then
        return ctx:fail("재접속 정리 뒤 친구 목록이 늘지 않음: " .. bot:buddy_capacity())
    end
    if settle.resume(bot) == false then
        return ctx:fail("정산 타이머 재개 실패")
    end
    settle.tick(ctx)
    return holds(ctx, bot, DEFAULT + 2 * STEP, MESO - 2 * FEE, "재접속 정리 뒤 (중복 확장·환불 없음)")
end

test_suite {
    name = "Buddy: 친구 목록 확장과 결과 유실 정리",
    bot_count = 2,

    on_scenario_finished = function(ctx)
        settle.resume_all(ctx, ctx:bot_count())
    end,

    scenarios = {
        expanded,
        lost_timer,
        unsent_relogin,
    },
}
