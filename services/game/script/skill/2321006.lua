-- Skill name (String.wz/Skill.img.xml): 리저렉션

local util = require("script/lib/skill")

function on_activated_2321006(me, skill, params)
	util.apply_resurrection(me, skill)
end
