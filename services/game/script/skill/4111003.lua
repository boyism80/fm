-- Skill name (String.wz/Skill.img.xml): 쉐도우 웹

local combat = require("script/lib/combat")

function on_activated_4111003(me, skill, params)
	combat.apply_shadow_web_skill(me, skill)
end
