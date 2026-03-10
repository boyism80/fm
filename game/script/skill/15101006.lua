-- Skill name (String.wz/Skill.img.xml): 라이트닝 차지

function on_activated(me, skill)
    me:add_buff(skill, BuffFlag.WkCharge, 1)
end
