-- Mob name (String.wz/Mob.img.xml): 검은 마법사

local boss_entry = require("script/lib/boss_entry")

return {
    on_mob_die = function(mob, attacker, map)
        boss_entry.on_phase_cleared(map, attacker, "black_mage")
    end
}
