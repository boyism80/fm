-- NPC name (String.wz/Npc.img.xml): 달리어

local medal = require("script/lib/medal")

return {
    on_click = function(me, npc)
        local titles = {}
        for i, kind in ipairs(medal.BOARDS) do
            titles[i] = medal.title(kind)
        end
        local sel = me:dialog_list(npc, "보고싶은 랭킹을 선택하시오.", titles)
        if sel == nil then
            return
        end
        me:dialog(npc, medal.ranking_text(medal.BOARDS[sel]))
    end,
}
