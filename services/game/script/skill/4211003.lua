-- Skill name (String.wz/Skill.img.xml): 픽파킷

function on_activated_4211003(me, skill, params)
    local effect = skill:effect()
    if effect == nil then
        return
    end

    me:buff(skill, BuffFlag.Pickpocket, effect.x)
end

function on_unbuff_4211003(me, skill)
end

