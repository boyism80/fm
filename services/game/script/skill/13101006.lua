-- Skill name (String.wz/Skill.img.xml): 윈드워크

local util = require("script/lib/skill")

return {
    on_activated = function(me, skill, params)
        util.apply_dark_sight(me, skill)
    end
}
