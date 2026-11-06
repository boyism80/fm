-- Skill name (String.wz/Skill.img.xml): 민첩한 몸놀림

return {
    on_activated = function(me, skill, params)
        local effect = skill:effect()
        if effect == nil then
            return
        end
        me:buff(skill, { [BuffFlag.Speed] = effect.speed })
    end
}
