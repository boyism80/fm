-- Skill name (String.wz/Skill.img.xml): 알바트로스

local util = require("script/lib/skill")

return {
    on_activated = function(me, skill, params)
        local morph = Morph.Albatross
        if me:gender() == Gender.Female then
            morph = Morph.AlbatrossFemale
        end
        util.apply_morph(me, skill, morph)
    end
}
