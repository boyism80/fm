-- Mob name (String.wz/Mob.img.xml): 핑크빈

local pinkbean = require("script/lib/pinkbean")
local medal = require("script/lib/medal")

return {
    on_mob_die = function(mob, attacker, map)
        medal.count_boss_kill(map, 29502)
        pinkbean.clear_map_after(map, 3000)
    end
}
