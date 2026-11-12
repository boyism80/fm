-- Skill name (String.wz/Skill.img.xml): 민첩한 몸놀림

return {
    on_activated = function(me, skill, params)
        me:buff(skill)
    end
}
