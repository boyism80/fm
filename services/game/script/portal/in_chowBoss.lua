local pq = require("script/lib/party_quest")

local SCRIPT = "script/portal/in_chowBoss.lua"
local ROOM = 300010420

local BUSY = "이미 누군가 안에 입장해있습니다."

return {

    on_enter = function(me)
        local refused = pq.enter_rooms(me:map(), SCRIPT, { { ROOM } }, BUSY, function()
            me:map(ROOM)
        end)
        if refused ~= nil then
            me:message(refused, Msg.PinkText)
        end
    end
}
