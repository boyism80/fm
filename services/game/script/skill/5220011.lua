-- Skill name (String.wz/Skill.img.xml): 어드밴스드 호밍

local util = require("script/lib/skill")

return {
    on_attack = function(me, skill, damages)
        util.mark_homing(me, skill, damages)
    end
}
