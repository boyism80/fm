-- Skill name (String.wz/Skill.img.xml): 헤이스트

function on_activated_9001001(me, skill, params)
    local effect = skill:effect()
    if effect == nil then
        return
    end

    me:buff(skill, {
        [BuffFlag.Speed] = effect.speed,
        [BuffFlag.Jump] = effect.jump,
    })
end
