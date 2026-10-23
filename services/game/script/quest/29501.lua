-- 칭호 도전 - 혼테일 슬레이어 (Quest.wz/QuestData/29501.img): 칭호 도전 - 혼테일 슬레이어

local medal = require("script/lib/medal")

local quest_id = 29501

return {
    on_start = function(me, npc)
        medal.start_slayer(me, npc, quest_id, medal.HORNTAIL_SLAYER, "혼테일")
    end,

    on_end = function(me, npc)
        medal.end_slayer(me, npc, quest_id, medal.HORNTAIL_SLAYER)
    end,
}
