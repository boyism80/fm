-- Skill name (String.wz/Skill.img.xml): 매직 가드

function on_activated_12001001(me, skill, params)
    local effect = skill:effect()
    if effect == nil then
        return
    end

    me:buff(skill, BuffFlag.MagicGuard, effect.x)
end
