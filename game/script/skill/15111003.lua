-- Skill name (String.wz/Skill.img.xml): 쇼크웨이브

function on_activating_15111003(me, skill, params)
	local morph = me:buff_value(BuffFlag.Morph)
	if morph == Morph.Transform or morph == Morph.TransformFemale then
		return true
	end
	return false
end

function on_activated_15111003(me, skill, params)
end
