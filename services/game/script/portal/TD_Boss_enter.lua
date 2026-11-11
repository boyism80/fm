local pq = require("script/lib/party_quest")

local SCRIPT = "script/portal/TD_Boss_enter.lua"
local ROOMS = 6

local BUSY = "이미 다른 누군가가 들어가 있는 것 같다. 지금은 들어갈 수 없을 것 같다."

return {

    on_enter = function(me)
        local rooms = {}
        for i = 1, ROOMS do
            rooms[#rooms + 1] = { me:map():wz():id() + i }
        end
        local refused = pq.enter_rooms(me:map(), SCRIPT, rooms, BUSY, function(index)
            me:play_portal_sound()
            me:map(rooms[index][1], 1)
        end)
        if refused ~= nil then
            me:message(refused, Msg.PinkText)
        end
    end
}
