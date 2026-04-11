-- Skill name (String.wz/Skill.img.xml): 쉐도우 파트너

function on_activated_14111000(me, skill, params)
    local effect = skill:effect()
    if effect == nil then
        return
    end
    me:buff(skill, BuffFlag.ShadowPartner, effect.x)
end

