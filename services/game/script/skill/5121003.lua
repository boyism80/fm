-- Skill name (String.wz/Skill.img.xml): 슈퍼트랜스폼

local util = require("script/lib/skill")

return {
    on_activated = function(me, skill, params)
        local morph = Morph.SuperTransform
        if me:gender() == Gender.Female then
            morph = Morph.SuperTransformFemale
        end
        util.apply_morph(me, skill, morph)
    end
}
