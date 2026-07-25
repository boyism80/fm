-- Skill name (String.wz/Skill.img.xml): 인빈서블

local util = require("script/lib/skill")

return {
	on_activated = function(me, skill, params)
		util.apply_buff_from_effect(me, skill, BuffFlag.Invincible, "x")
	end
}
