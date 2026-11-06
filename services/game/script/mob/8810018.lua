-- Mob name (String.wz/Mob.img.xml): 혼테일

local horntail = require("script/lib/horntail")
local medal = require("script/lib/medal")

-- The WZ kill buff (info/buff), whose effect everyone it reached shows.
local NINE_SPIRIT_BUFF = 2022108

return {
    on_mob_die = function(mob, attacker, map)
        for _, ch in pairs(map:characters()) do
            if ch:is_alive() then
                ch:show_nine_spirit_effect(NINE_SPIRIT_BUFF)
            end
        end
        medal.count_boss_kill(map, 29501)
        horntail.clear_map_after(map, 3000)
    end
}
