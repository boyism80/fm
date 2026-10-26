local pq = require("script/integration/lib/party_quest")

local M = {}

local LITH_HARBOR = 104000000
local GUIDE = 9000001
local DIALOG_LIST = 4
local JOIN = 2

-- Schedules the event, joins it through the guide NPC and checks the bot lands on the event's first map.
function M.join(ctx, bot, key, name, first_map)
    if pq.command(bot, "/봇초기화 30 0 0 - -", "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    if bot:instance_move(LITH_HARBOR) == false then
        return ctx:fail("리스항구 이동 실패")
    end
    local announced = bot:request(resp.notice, req.normal_chat { message = "/일정이벤트 " .. key }, function(p)
        return p.message:find(name, 1, true) ~= nil and p.message:find("잠시 후", 1, true) ~= nil
    end, 15000)
    if announced == false then
        return ctx:fail(name .. " 모집 공지가 오지 않음")
    end
    local oid = pq.npc(ctx, bot, GUIDE)
    if oid == false then
        return false
    end
    if bot:npc_click(oid) == false then
        return ctx:fail("도우미 대화가 열리지 않음")
    end
    if bot:request(resp.warp, req.dialog { dialog_type = DIALOG_LIST, next = true, selected = JOIN }, nil, 5000) == false then
        return ctx:fail("이벤트 참가로 이동하지 않음")
    end
    if bot:map() ~= first_map then
        return ctx:fail(name .. " 첫 맵으로 이동하지 않음: " .. bot:map())
    end
    return true
end

-- Ends recruiting at once and waits for the countdown, matching the start with the game's own first packet.
function M.start(ctx, bot, name, packet, match)
    local started = bot:request(packet, req.normal_chat { message = "/이벤트시작" }, match, 60000)
    if started == false then
        return ctx:fail(name .. " 카운트다운 뒤 이벤트가 시작되지 않음")
    end
    return started
end

function M.cancel(bot)
    pq.command(bot, "/이벤트취소", "이벤트")
end

return M
