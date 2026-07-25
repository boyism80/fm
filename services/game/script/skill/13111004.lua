-- Skill name (String.wz/Skill.img.xml): 퍼펫

local util = require("script/lib/skill")

return {
	on_activated = function(me, skill, params)
		util.mark_archer_puppet(me, skill, params)
	end,

	on_buff = function(me, skill)
		util.apply_archer_puppet_buff(me, skill)
	end,

	on_unbuff = function(me, skill)
		util.apply_archer_puppet_unbuff(me, skill)
	end
}
