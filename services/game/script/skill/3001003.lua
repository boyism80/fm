-- Skill name (String.wz/Skill.img.xml): 포커스

return {
    on_activated = function(me, skill, params)
        me:buff(skill)
    end
}
