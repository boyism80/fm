-- Skill name (String.wz/Skill.img.xml): 파워 가드

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
    me:buff(skill, BuffFlag.Powerguard, x)
end

function on_unbuff(me, skill)
end

