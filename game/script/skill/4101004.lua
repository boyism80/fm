-- Skill name (String.wz/Skill.img.xml): 헤이스트

function on_activated(me, skill)
    local wz = skill:wz()
    if wz == nil or wz.effects == nil then
        return
    end

    local effect = wz.effects[skill:level()]
    if effect == nil then
        return
    end

    local speed = effect.speed or 0
    local jump = effect.jump or 0
    me:add_buff(skill, {
        [BuffFlag.Speed] = speed,
        [BuffFlag.Jump] = jump,
    })
end
