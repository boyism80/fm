-- Skill name (String.wz/Skill.img.xml): 매직 가드

function on_activated(me, skill, params)
    local wz = skill:wz()
    if wz == nil or wz.effects == nil then
        return
    end

    local effect = wz.effects[skill:level()]
    if effect == nil then
        return
    end

    local guard_percent = effect.x or 0
    me:buff(skill, BuffFlag.MagicGuard, guard_percent)
end

function on_unbuff(me, skill)
end

