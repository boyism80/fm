-- Skill name (String.wz/Skill.img.xml): 숨기

function on_activated(me, skill, params)
	me:hidden(not me:hidden())
end
