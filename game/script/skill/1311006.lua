-- Skill name (String.wz/Skill.img.xml): 드래곤 로어
-- MP 16/24/30 소비, HP의 59%~30% 소비, 15마리 공격, 발동 시 자기 스턴 (4/3/2초). HP 50% 이상일 때만 사용 가능.

function on_activated_1311006(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	local cur_hp = me:hp()
	local max_hp = me:max_hp()
	if max_hp <= 0 then
		return
	end
	if cur_hp * 100 < max_hp * 50 then
		return
	end
	local x = effect.x or 0
	if x > 0 then
		local hp_loss = math.floor(cur_hp * x / 100)
		if hp_loss > 0 then
			local new_hp = math.max(cur_hp - hp_loss, 1)
			me:hp(new_hp, true)
		end
	end
	local stun_sec = effect.y or 0
	if stun_sec > 0 then
		local duration_ms = stun_sec * 1000
		me:debuff(DebuffFlag.Stun, duration_ms)
	end
end
