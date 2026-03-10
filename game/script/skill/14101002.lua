-- Skill name (String.wz/Skill.img.xml): 자벨린 부스터

function on_activated(me, skill)
    local wz = skill:wz()
    if wz == nil or wz.effects == nil then
        return
    end

    local effect = wz.effects[skill:level()]
    if effect == nil then
        return
    end

    local x = effect.x or 0
    me:add_buff(skill, BuffFlag.Booster, x)
end

function on_deactivated(me, skill)
end
