local pq = require("script/lib/party_quest")

local SCRIPT = "script/portal/TD_neo_inTree.lua"
local EXIT = 240070000
local QUESTS = {
    { quest = 3719, map = 240070010 },
    { quest = 3724, map = 240070020 },
    { quest = 3730, map = 240070030 },
    { quest = 3736, map = 240070040 },
    { quest = 3742, map = 240070050 },
    { quest = 3748, map = 240070060 },
}
local ROOMS = 10

local BUSY = "이미 다른 누군가가 들어가 있는 것 같다. 지금은 들어갈 수 없을 것 같다."

local function enter(me, room)
    me:play_portal_sound()
    local function on_arrive(me)
        me:clock(180, function(me)
            me:map(EXIT)
        end)
    end
    me:map(room, 0, { callback = on_arrive })
end

return {

    on_enter = function(me)
        local base = nil
        for _, entry in ipairs(QUESTS) do
            if me:quest(entry.quest):started() then
                base = entry.map
            end
        end
        if base == nil then
            me:message("퀘스트를 진행 중일 때만 들어갈 수 있습니다.", Msg.PinkText)
            return
        end
        local rooms = {}
        for i = 0, ROOMS - 1 do
            rooms[#rooms + 1] = { base + i }
        end
        local refused = pq.enter_rooms(me:map(), SCRIPT, rooms, BUSY, function(index)
            enter(me, rooms[index][1])
        end)
        if refused ~= nil then
            me:message(refused, Msg.PinkText)
        end
    end
}
