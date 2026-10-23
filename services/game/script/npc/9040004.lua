-- NPC name (String.wz/Npc.img.xml): 명예로운 바위

local class_names = require("script/lib/class")

local COUNT = 30

return {
    on_click = function(me, npc)
        local entries = me:overall_ranking(COUNT)
        if entries == nil then
            me:dialog(npc, "바위에 새겨진 글자가 흐릿해 읽을 수가 없다. 잠시 후 다시 살펴보자.")
            return
        end
        if #entries == 0 then
            me:dialog(npc, "아직 바위에 새겨진 이름이 없다.")
            return
        end
        local lines = {}
        for _, entry in ipairs(entries) do
            lines[#lines + 1] = entry.rank ..
                "위 #e#b" .. entry.name .. "#k#n  Lv. #r" .. entry.level .. "#k  " .. class_names.name(entry.class)
        end
        me:dialog(npc,
            "이 월드에서 가장 뛰어난 모험가들의 이름이 새겨진 신비로운 바위다.\r\n랭킹은 #r" ..
            COUNT .. "#k위까지만 표기된다.\r\n\r\n" .. table.concat(lines, "\r\n"))
    end,
}
