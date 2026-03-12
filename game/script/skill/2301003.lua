-- Skill name (String.wz/Skill.img.xml): 인빈서블

function on_activated(me, skill)
	Skill.apply_invincible(me, skill)
end

function on_unbuff(me, skill)
end

