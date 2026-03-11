-- Skill name (String.wz/Skill.img.xml): 소울 차지

function on_activated(me, skill)
    me:buff(skill, BuffFlag.WkCharge, 1)
end

