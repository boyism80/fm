-- Skill name (String.wz/Skill.img.xml): 트랜스폼

function on_activated_15111002(me, skill, params)
	local morph = Morph.Transform
	if me:gender() == Gender.Female then
		morph = Morph.TransformFemale
	end
	me:buff(skill, BuffFlag.Morph, morph)
end
