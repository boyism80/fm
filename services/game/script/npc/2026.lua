-- NPC name (String.wz/Npc.img.xml): 할로 캣

local boss_entry = require("script/lib/boss_entry")
local bosses = require("script/lib/bosses")

return {
    on_click = function(me, npc)
        boss_entry.on_click(me, npc, bosses.black_mage)
    end
}
