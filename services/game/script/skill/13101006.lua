-- Skill name (String.wz/Skill.img.xml): 윈드워크

local util = require("script/lib/skill")

return {
	on_activated = function(me, skill, params)
		util.apply_buff_fixed(me, skill, BuffFlag.Darksight, 1)
	end
}
