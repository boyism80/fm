-- Skill name (String.wz/Skill.img.xml): 픽파킷

function on_activated_4211003(me, skill, params)
    local wz = skill:wz()
    if wz == nil or wz.effects == nil then
        return
    end

    local effect = wz.effects[skill:level()]
    if effect == nil then
        return
    end

    local x = effect.x or 0
    me:buff(skill, BuffFlag.Pickpocket, x)
end

function on_unbuff_4211003(me, skill)
end

