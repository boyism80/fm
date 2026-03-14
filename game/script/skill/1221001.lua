-- Skill name (String.wz/Skill.img.xml): 몬스터 마그넷

function on_preactivated(me, skill, params)
	return true
end

function on_activated(me, skill, params)
	apply_monster_magnet(me, skill, params)
end
