-- Skill name (String.wz/Skill.img.xml): 파이널 어택

local util = require("script/lib/skill")

return {
    follow_up_attack = true,

    on_activated = function(me, skill, params)
        util.apply_buff_from_effect(me, skill, BuffFlag.FinalAttack, "x")
    end
}
