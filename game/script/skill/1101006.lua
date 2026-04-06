-- Skill name (String.wz/Skill.img.xml): 분노

function on_activated_1101006(me, skill, params)
    local wz = skill:wz()
    if wz == nil or wz.effects == nil then
        return
    end

    local effect = skill:effect()
    if effect == nil then
        return
    end

    local pad = effect.pad or 0
    me:buff(skill, BuffFlag.WeaponAtk, pad)
    -- TODO: Apply this buff to nearby party members when party system is implemented.
end

