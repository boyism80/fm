-- Skill name (String.wz/Skill.img.xml): 디바인 차지 : 둔기

function on_activated(me, skill)
    me:buff(skill, BuffFlag.WkCharge, 1)
end

