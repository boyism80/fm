local pq = require("script/lib/party_quest")

local M = {}

local ROOMS = 39
local ENTRY_KEY = "mini_dungeon_entry"
local BUSY = "모든 미니던전 인스턴스가 사용중입니다. 나중에 다시 시도하세요."

local function enter(me, party, first_room, rooms)
    for i = 0, (rooms or ROOMS) - 1 do
        local room = first_room + i
        local ok, count = run_on_map(room, "script/lib/party_quest.lua", "character_count")
        if ok and count == 0 then
            local pid = party:id()
            for _, ch in pairs(me:map():characters()) do
                local p = ch:party()
                if p ~= nil and p:id() == pid then
                    ch:map(room, 0)
                end
            end
            return true
        end
    end
    return false
end

function M.enter(me, base, first_room, rooms)
    if me:map():wz():id() ~= base then
        me:play_portal_sound()
        me:map(base, "MD00")
        return
    end

    local party = me:party()
    if party == nil then
        me:message("파티상태가 아니여서 들어갈 수 없습니다.", Msg.PinkText)
        return
    end
    if party:leader_id() ~= me:id() then
        me:message("파티장이 아닙니다.", Msg.PinkText)
        return
    end

    local lobby = me:map()
    if not pq.hold_entry(lobby, ENTRY_KEY) then
        me:message(BUSY, Msg.PinkText)
        return
    end
    local entered = enter(me, party, first_room, rooms)
    pq.release_entry(lobby, ENTRY_KEY)
    if not entered then
        me:message(BUSY, Msg.PinkText)
    end
end

return M
