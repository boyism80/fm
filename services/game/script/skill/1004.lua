-- Skill name (String.wz/Skill.img.xml): 몬스터 라이딩

return {
	on_activated = function(me, skill, params)
		return me:ride(skill)
	end
}
