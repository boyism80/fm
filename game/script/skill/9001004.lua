-- Skill name (String.wz/Skill.img.xml): 숨기

function on_activated_9001004(me, skill, params)
	me:hidden(not me:hidden())
end
