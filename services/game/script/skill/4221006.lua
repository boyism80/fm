-- Skill name (String.wz/Skill.img.xml): 연막탄

return {
    on_activated = function(me, skill, params)
        local effect = skill:effect()
        if effect == nil then
            return
        end
        if effect.time <= 0 then
            return
        end
        me:create_mist(skill, effect.time, MistType.Smoke, me:skill_area(skill), 0, 1.0)
    end
}
