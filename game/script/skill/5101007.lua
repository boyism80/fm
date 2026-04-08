-- Skill name (String.wz/Skill.img.xml): 오크통

function on_activated_5101007(me, skill, params)
	me:buff(skill, BuffFlag.Morph, Morph.Barrel)
end
