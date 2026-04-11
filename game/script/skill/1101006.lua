-- Skill name (String.wz/Skill.img.xml): 분노

function on_activated_1101006(me, skill, params)
    local effect = skill:effect()
    if effect == nil then
        return
    end

    me:buff(skill, BuffFlag.WeaponAtk, effect.pad)
    -- TODO: Apply this buff to nearby party members when party system is implemented.
end

