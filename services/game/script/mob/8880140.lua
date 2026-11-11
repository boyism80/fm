-- Mob name (String.wz/Mob.img.xml): 꿈속의 루시드

local boss_entry = require("script/lib/boss_entry")

return {
    on_mob_die = function(mob, attacker, map)
        boss_entry.on_phase_cleared(map, attacker, "lucid")
    end
}
