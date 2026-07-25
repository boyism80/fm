-- Skill name (String.wz/Skill.img.xml): 인피니티

local util = require("script/lib/skill")

return {
	on_activated = function(me, skill, params)
		util.apply_buff_fixed(me, skill, BuffFlag.Infinity, 1)
	end
}
