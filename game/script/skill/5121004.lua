-- Skill name (String.wz/Skill.img.xml): 데몰리션

function on_activating_5121004(me, skill, params)
	local morph = me:buff_value(BuffFlag.Morph)
	if morph == Morph.SuperTransform or morph == Morph.SuperTransformFemale then
		return true
	end
	return false
end

function on_activated_5121004(me, skill, params)
end
