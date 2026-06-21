-- Skill name (String.wz/Skill.img.xml): 아이언 바디

local util = require("script/lib/skill")

function on_activated_1001003(me, skill, params)
	util.apply_iron_body(me, skill)
end
