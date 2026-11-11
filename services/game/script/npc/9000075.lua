-- NPC name (String.wz/Npc.img.xml): 찬

local pq = require("script/lib/party_quest")
local item_count = pq.item_count

local SCRIPT = "script/npc/9000075.lua"
local BUSY = "누군가 맵에 들어가있습니다."

local function enter(me, target, ticket, free_ticket)
    if ticket >= 1 and free_ticket < 1 then
        local code = me:exchange({ item = { [4001431] = 1 } }, nil)
        if code ~= ExchangeResult.OK then
            return "골든티켓이 없으시면 입장하실 수 없습니다."
        end
    end
    me:map(target, 0)
    return nil
end

return {

    on_click = function(me, npc)
        local sel = me:dialog_list(npc,
            "정말 꼭 입장하고 싶으신 것 같군요. 좀 곤라합니다만, 골든티켓을 갖고 계시다면 들어가실 수도 있을 것 같군요. 어느 곳에 가길 희망하시는 겁니까? #k", {
                "도깨비 동굴1(Lv.43 파란 도깨비)",
                "도깨비 동굴2(Lv.54 빨간 도깨비",
                "도깨비 동굴3(Lv.66 힘센 돌 도깨비)",
            })
        if sel == nil then
            return
        end
        local ticket = item_count(me, 4001431)
        local free_ticket = item_count(me, 4001432)
        if ticket < 1 and free_ticket < 1 then
            me:dialog(npc, "골든티켓이 없으시면 입장하실 수 없습니다.")
            return
        end
        local target = me:map():wz():id() + 500 + (sel - 1) * 100
        local refused = pq.enter_rooms(me:map(), SCRIPT, { { target } }, BUSY, function()
            return enter(me, target, ticket, free_ticket)
        end)
        if refused ~= nil then
            me:dialog(npc, refused)
        end
    end
}
