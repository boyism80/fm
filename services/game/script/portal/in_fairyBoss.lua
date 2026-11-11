local pq = require("script/lib/party_quest")

local SCRIPT = "script/portal/in_fairyBoss.lua"
local ROOM = 300030310

local REFUSED = "이미 누군가 안에 입장해있거나 퀘스트를 진행중이시지 않습니다."

return {

    on_enter = function(me)
        if me:quest(31229):started() == false then
            me:message(REFUSED, Msg.PinkText)
            return
        end
        local refused = pq.enter_rooms(me:map(), SCRIPT, { { ROOM } }, REFUSED, function()
            me:map(ROOM)
        end)
        if refused ~= nil then
            me:message(refused, Msg.PinkText)
        end
    end
}
