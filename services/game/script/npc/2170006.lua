-- NPC name (String.wz/Npc.img.xml): 콜로세움

local pq = require("script/lib/party_quest")

local SCRIPT = "script/npc/2170006.lua"
local BASE_MAP = 200101500
local ROOMS = 10
local BUSY = "모든 지역에 다른 파티가 입장 중입니다."

local function enter(me, party, map_id)
    run_on_map(map_id, SCRIPT, "spawn_boss")
    local pid = party:id()
    for _, ch in pairs(me:map():characters()) do
        local p = ch:party()
        if p ~= nil and p:id() == pid then
            ch:map(map_id, 1)
        end
    end
end

return {

    spawn_boss = function(map)
        map:spawn_mob(6160003, 642, 227)
    end,

    on_click = function(me, npc)
        local sel = me:dialog_list(npc, "위험지역! 크세르크세스가 있는 곳입니다. 입장하시겠습니까?\r\n\r\n입장조건 : 레벨 50이상#b", {
            "입장 한다.",
            "되돌아 간다.",
        })
        if sel ~= 1 then
            return
        end
        local party = me:party()
        if party == nil then
            me:dialog(npc, "파티를 맺어야만 입장하실 수 있습니다.")
            return
        end
        if party:leader_id() ~= me:id() then
            me:dialog(npc, "파티장만 입장하실 수 있습니다.")
            return
        end
        local rooms = {}
        for i = 0, ROOMS - 1 do
            if id2map(BASE_MAP + i) ~= nil then
                rooms[#rooms + 1] = { BASE_MAP + i }
            end
        end
        local refused = pq.enter_rooms(me:map(), SCRIPT, rooms, BUSY, function(index)
            enter(me, party, rooms[index][1])
        end)
        if refused ~= nil then
            me:dialog(npc, refused)
        end
    end
}
