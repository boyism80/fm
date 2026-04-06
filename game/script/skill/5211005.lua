-- Skill name (String.wz/Skill.img.xml): 쿨링 이펙트

function on_attack_5211005(me, skill, damages)
	apply_prob_status_on_skill_hit(me, skill, damages, MobBuff.Freeze)
end

function on_activated_5211005(me, skill, params)
end
