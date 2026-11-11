local pq = require("script/lib/party_quest")

local SCRIPT = "script/portal/s4hitman.lua"
local ROOM = 910200000
local EXIT = 101030104

local function enter(me)
    me:play_portal_sound()
    local function on_arrive(me)
        me:clock(1200, function(me)
            me:map(EXIT)
        end)
    end
    me:map(ROOM, 0, { callback = on_arrive })
end

return {

    on_enter = function(me)
        if me:quest(6201):started() == false then
            me:message("알 수 없는 힘으로 봉인되어 있습니다.")
            return
        end
        if pq.has_item(me, 4031452) then
            me:message("샨이 요청한 것을 끝냈습니다. 더 이상 입장할 필요가 없습니다.")
            return
        end
        local refused = pq.enter_rooms(me:map(), SCRIPT, { { ROOM } }, "이미 다른 누군가가 퀘스트에 도전중입니다.", function()
            enter(me)
        end)
        if refused ~= nil then
            me:message(refused, Msg.PinkText)
        end
    end
}
