-- NPC name (String.wz/Npc.img.xml): 다오

local pq = require("script/lib/party_quest")
local item_count = pq.item_count

local SCRIPT = "script/npc/9000080.lua"
local BUSY = "누군가 이미 맵에 들어가있습니다."

local function enter(me, target, ticket, free_ticket)
    if ticket >= 1 and free_ticket < 1 then
        local code = me:exchange({ item = { [4001431] = 1 } }, nil)
        if code ~= ExchangeResult.OK then
            return "골드 티켓이 필요합니다."
        end
    end
    me:map(target, 0)
    return nil
end

return {

    on_click = function(me, npc)
        local sel = me:dialog_list(npc, "어느 사원에 들어가고 싶으세요?\r\n#r입장을 하기 위해서는 골든티켓이 필요합니다.#k", {
            "원숭이 사원1(Lv.15 사나운 원숭이)",
            "원숭이 사원2(Lv.21 어미 원숭이)",
            "원숭이 사원3(Lv.27 흰털 아기 원숭이)",
            "원숭이 사원4(Lv.34 흰털 어미 원숭이)",
        })
        if sel == nil then
            return
        end
        local ticket = item_count(me, 4001431)
        local free_ticket = item_count(me, 4001432)
        if ticket < 1 and free_ticket < 1 then
            me:dialog(npc, "골드 티켓이 필요합니다.")
            return
        end
        local target = me:map():wz():id() + 100 + (sel - 1) * 100
        local refused = pq.enter_rooms(me:map(), SCRIPT, { { target } }, BUSY, function()
            return enter(me, target, ticket, free_ticket)
        end)
        if refused ~= nil then
            me:dialog(npc, refused)
        end
    end
}
