-- Skill name (String.wz/Skill.img.xml): 스피릿 자벨린

function on_activated(me, skill)
    local wz = skill:wz()
    if wz == nil or wz.effects == nil then
        return
    end
    local effect = wz.effects[skill:level()]
    local value = 1
    if effect ~= nil and effect.x ~= nil then
        value = effect.x
    end
    me:buff(skill, BuffFlag.SpiritClaw, value)
end
