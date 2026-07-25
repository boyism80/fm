-- Skill name (String.wz/Skill.img.xml): 숨기

return {
	on_activated = function(me, skill, params)
		me:hidden(not me:hidden())
	end
}
