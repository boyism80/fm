-- Skill name (String.wz/Skill.img.xml): 파이어 차지 : 둔기

function on_activated(me, skill)
    me:buff(skill, BuffFlag.WkCharge, 1)
end

