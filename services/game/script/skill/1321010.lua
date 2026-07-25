-- Skill name (String.wz/Skill.img.xml): 용사의 의지

local util = require("script/lib/skill")

return {
	on_activated = function(me, skill, params)
		util.apply_hero_will(me)
	end
}
