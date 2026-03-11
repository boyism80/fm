-- Skill name (String.wz/Skill.img.xml): 블리자드 차지 : 검

function on_activated(me, skill)
    me:buff(skill, BuffFlag.WkCharge, 1)
end

