-- NPC name (String.wz/Npc.img.xml): 다오

local item_count = require("script/lib/party_quest").item_count

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
        local ok, count = run_on_map(target, "script/lib/party_quest.lua", "character_count")
        if ok and count ~= nil and count > 0 then
            me:dialog(npc, "누군가 이미 맵에 들어가있습니다.")
            return
        end

        if ticket >= 1 and free_ticket < 1 then
            local code = me:exchange({ item = { [4001431] = 1 } }, nil)
            if code ~= ExchangeResult.OK then
                me:dialog(npc, "골드 티켓이 필요합니다.")
                return
            end
        end
        me:map(target, 0)
    end
}
