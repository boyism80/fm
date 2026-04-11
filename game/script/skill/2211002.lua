-- Skill name (String.wz/Skill.img.xml): 아이스 스트라이크

function on_attack_2211002(me, skill, damages)
	apply_prob_status_on_skill_hit(me, skill, damages, MobBuff.Freeze)
end

function on_activated_2211002(me, skill, params)
end
