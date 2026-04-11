-- Skill name (String.wz/Skill.img.xml): 스내치

function on_activating_5121005(me, skill, params)
	local morph = me:buff_value(BuffFlag.Morph)
	if morph == Morph.SuperTransform or morph == Morph.SuperTransformFemale then
		return true
	end
	return false
end

function on_activated_5121005(me, skill, params)
end
