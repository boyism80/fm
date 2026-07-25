-- Skill name (String.wz/Skill.img.xml): 알바트로스

return {
	on_activated = function(me, skill, params)
		local morph = Morph.Albatross
		if me:gender() == Gender.Female then
			morph = Morph.AlbatrossFemale
		end
		me:buff(skill, BuffFlag.Morph, morph)
	end
}
