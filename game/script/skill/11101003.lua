-- Skill name (String.wz/Skill.img.xml): 분노

function on_activated_11101003(me, skill, params)
    local effect = skill:effect()
    if effect == nil then
        return
    end

    me:buff(skill, BuffFlag.WeaponAtk, effect.pad)
end
