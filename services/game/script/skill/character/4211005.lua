-- Skill name (String.wz/Skill.img.xml): 메소 가드

function on_activated_4211005(me, skill, params)
    local effect = skill:effect()
    if effect == nil then
        return
    end
    me:buff(skill, BuffFlag.MesoGuard, effect.x)
end

