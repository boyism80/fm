local pq = require("script/lib/party_quest")

local SCRIPT = "script/portal/skyrom.lua"
local ROOM = 926000010

local BUSY = "이미 이 안에 다른 누군가가 들어가서 스카이롬을 훔치는 중입니다."

return {

    shuffle_reactors = function(map)
        pq.shuffle_reactors(map)
    end,

    on_enter = function(me)
        if me:quest(3935):started() == false or pq.has_item(me, 4031574) then
            return
        end
        local refused = pq.enter_rooms(me:map(), SCRIPT, { { ROOM } }, BUSY, function()
            run_on_map(ROOM, SCRIPT, "shuffle_reactors")
            me:play_portal_sound()
            me:map(ROOM)
        end)
        if refused ~= nil then
            me:message(refused, Msg.PinkText)
        end
    end
}
