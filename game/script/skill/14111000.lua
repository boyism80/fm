-- Skill name (String.wz/Skill.img.xml): 쉐도우 파트너 (4th job)

function on_activated(me, skill)
    local wz = skill:wz()
    if wz == nil or wz.effects == nil then
        return
    end
    local effect = wz.effects[skill:level()]
    if effect == nil then
        return
    end
    local percent = effect.x or 0
    me:add_buff(skill, BuffFlag.ShadowPartner, percent)
end
