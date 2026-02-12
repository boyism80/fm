-- Skill name (String.wz/Skill.img.xml): 숨기

function on_active(me, skill)
	me:set_hidden(not me:is_hidden())
end
