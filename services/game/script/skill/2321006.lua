-- Skill name (String.wz/Skill.img.xml): 리저렉션

local util = require("script/lib/skill")

return {
	on_activated = function(me, skill, params)
		util.apply_resurrection(me, skill)
	end
}
