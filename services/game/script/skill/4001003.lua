-- Skill name (String.wz/Skill.img.xml): 다크 사이트

local util = require("script/lib/skill")

return {
    on_activated = function(me, skill, params)
        util.apply_dark_sight(me, skill)
    end
}
