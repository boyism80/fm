-- Damage handlers: magic guard, meso guard (on_damaged)

function handle_magic_guard(me, attacker, skill, damage)
    if damage == nil or damage <= 0 then
        return damage
    end

    if me:buff_value(BuffFlag.MagicGuard) == nil then
        return damage
    end

    local effect = get_skill_effect(me:skill(Skill.MagicGuard))
    if effect == nil then
        effect = get_skill_effect(me:skill(Skill.MagicGuardCygnus))
    end
    local guard_percent = effect and effect.x
    if guard_percent == nil or guard_percent <= 0 then
        return damage
    end

    local current_mp = me:mp()
    if current_mp == nil or current_mp <= 0 then
        return damage
    end

    local mp_loss = math.floor(damage * (guard_percent / 100.0))
    if mp_loss < 0 then
        mp_loss = 0
    end
    if mp_loss > current_mp then
        mp_loss = current_mp
    end

    local hp_loss = damage - mp_loss
    if hp_loss < 0 then
        hp_loss = 0
    end

    if mp_loss > 0 then
        me:add_mp(-mp_loss)
    end

    return hp_loss
end

function handle_meso_guard(me, attacker, skill, damage)
    if damage == nil or damage <= 0 then
        return damage
    end
    if me:buff_value(BuffFlag.MesoGuard) == nil then
        return damage
    end

    local effect = get_skill_effect(me:skill(Skill.MesoGuard))
    local guard_percent = effect and effect.x
    if guard_percent == nil or guard_percent <= 0 then
        return damage
    end
    local current_meso = me:meso()
    if current_meso == nil or current_meso <= 0 then
        return damage
    end
    local meso_loss = math.floor(damage * (guard_percent / 100.0))
    if meso_loss < 0 then
        meso_loss = 0
    end
    if meso_loss > current_meso then
        meso_loss = current_meso
    end
    local hp_loss = damage - meso_loss
    if hp_loss < 0 then
        hp_loss = 0
    end
    if meso_loss > 0 then
        me:meso(current_meso - meso_loss)
    end
    return hp_loss
end

local function mob_reflect_cap(max_hp, divisor)
	if max_hp == nil or max_hp <= 0 or divisor == nil or divisor <= 0 then
		return 0
	end
	return math.floor(max_hp / divisor)
end

local function handle_power_guard_reflect(me, attacker, damage, reflect_ratio)
	if me:buff_value(BuffFlag.Powerguard) == nil then
		return damage
	end
	local reduce = math.floor(damage * (reflect_ratio / 100.0))
	if reduce < 0 then
		reduce = 0
	end
	if reduce > damage then
		reduce = damage
	end
	local hp_loss = damage - reduce
	if hp_loss < 0 then
		hp_loss = 0
	end
	if attacker ~= nil and attacker:is(ObjectType.Mob) then
		local bounced = reduce
		local cap = mob_reflect_cap(attacker:max_hp(), 10)
		bounced = math.min(bounced, cap)
		if bounced > 0 then
			attacker:damage(me, bounced)
		end
	end
	return hp_loss
end

local function handle_mana_reflection_reflect(me, attacker, damage, reflect_ratio)
	if me:buff_value(BuffFlag.ManaReflection) == nil then
		return damage
	end
	if attacker ~= nil and attacker:is(ObjectType.Mob) then
		local bounce = math.floor(damage * (reflect_ratio / 100.0))
		if bounce < 0 then
			bounce = 0
		end
        local cap = mob_reflect_cap(attacker:max_hp(), 20)
		bounce = math.min(bounce, cap)
		if bounce > 0 then
			attacker:damage(me, bounce)
		end
	end
	return damage
end

function handle_reflect_damage(me, attacker, skill, damage, params)
	if damage == nil or damage <= 0 or params == nil then
		return damage
	end
	local reflect_ratio = params.reflect_ratio
	if reflect_ratio == nil or reflect_ratio <= 0 then
		return damage
	end
	local hit_type = params.hit_type
	if hit_type == IncomingHit.Collide then
		return handle_power_guard_reflect(me, attacker, damage, reflect_ratio)
	end
	return handle_mana_reflection_reflect(me, attacker, damage, reflect_ratio)
end
