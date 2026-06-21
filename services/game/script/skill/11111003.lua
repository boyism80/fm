-- Skill name (String.wz/Skill.img.xml): 코마

local util = require("script/lib/skill")
local combat = require("script/lib/combat")

function on_attack_11111003(me, skill, damages)
	combat.apply_prob_status(me, skill, damages, MobBuff.Stun)
end

function on_activated_11111003(me, skill, params)
	-- Dawn Warrior Coma consumes combo orbs the same way as Hero Coma.
	util.consume_combo_orbs(me)
end
