-- Skill name (String.wz/Skill.img.xml): 쇼크웨이브

return {
	on_activating = function(me, skill, params)
		local morph = me:buff_value(BuffFlag.Morph)
		if morph == Morph.Transform or morph == Morph.TransformFemale then
			return true
		end
		return false
	end
}
