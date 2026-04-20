-- Skill name (String.wz/Skill.img.xml): 닌자 스톰

function on_attack_4121008(me, skill, damages)
	apply_prob_status_on_skill_hit(me, skill, damages, MobBuff.Stun)
end

function on_activated_4121008(me, skill, params)
end
