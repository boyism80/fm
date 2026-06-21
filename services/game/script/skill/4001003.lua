-- Skill name (String.wz/Skill.img.xml): 다크 사이트
local util = require("script/lib/skill")

function on_activated_4001003(me, skill, params)
	util.apply_buff_fixed(me, skill, BuffFlag.Darksight, 1)
end
