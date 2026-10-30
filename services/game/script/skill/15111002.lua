-- Skill name (String.wz/Skill.img.xml): 트랜스폼

local util = require("script/lib/skill")

return {
    on_activated = function(me, skill, params)
        local morph = Morph.Transform
        if me:gender() == Gender.Female then
            morph = Morph.TransformFemale
        end
        util.apply_morph(me, skill, morph)
    end
}
