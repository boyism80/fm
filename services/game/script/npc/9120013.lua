-- NPC name (String.wz/Npc.img.xml): 웃고있는고양이

local MAPS = {
    { map = 801000100, name = "남탕", tail = "#k으" },
    { map = 801000200, name = "여탕", tail = "#k으" },
}

return {
    on_click = function(me, npc)
        local names = {}
        for i, m in ipairs(MAPS) do
            names[i] = m.name
        end
        local sel = me:dialog_list(npc, "탈의실 엿보고싶냥?", names)
        if sel == nil then
            return
        end
        local selected = MAPS[sel]
        local price = 3000
        if me:class() == Class.Beginner then
            price = 300
        end
        if not me:dialog_yes_no(npc, "정말 #b" .. selected.name .. selected.tail .. "로 가볼테냥? 공짜는 아냥 요금은 #b" .. price .. "#k 메소이며, 구타 당해도 원망하기없기다냥") then
            me:dialog(npc, "여기까지와서 포기할거냥? 실망이다냥.")
            return
        end

        local code = me:exchange({ meso = price }, nil)
        if code ~= ExchangeResult.OK then
            me:dialog(npc, "흐음. 메소가 부족하신 것 같은데요? 요금이 없으시다면 선택하신 곳으로 보내드릴 수 없답니다.")
            return
        end
        me:map(selected.map)
    end
}
