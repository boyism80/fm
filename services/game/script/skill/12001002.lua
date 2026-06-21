-- Skill name (String.wz/Skill.img.xml): 매직 아머

local util = require("script/lib/skill")

function on_activated_12001002(me, skill, params)
	util.apply_iron_body(me, skill)
end
