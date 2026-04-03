-- Attack handlers: combo, pickpocket, ice charge, ammo consume

local function skill_effect_mob_limit(effect)
	if effect == nil then
		return nil
	end
	local c = effect.mob_count
	if c == nil then
		return nil
	end
	c = tonumber(c) or 0
	if c <= 0 then
		return nil
	end
	return c
end

local function total_damage_to_mob(hits)
	if hits == nil then
		return 0
	end
	local total = 0
	for _, amount in ipairs(hits) do
		if amount and amount > 0 then
			total = total + amount
		end
	end
	return total
end

function apply_skill_drain_on_attack(me, skill, damages)
	if me == nil or skill == nil or damages == nil then
		return
	end
	local effect = get_skill_effect(skill)
	if effect == nil then
		return
	end
	local pct = tonumber(effect.x) or 0
	if pct <= 0 then
		return
	end
	local player_max = me:max_hp()
	if player_max == nil or player_max <= 0 then
		return
	end
	local cap_half = math.floor(player_max / 2)
	local total_heal = 0
	for mob, hits in pairs(damages) do
		if mob == nil or hits == nil then
			goto continue_drain
		end
		local tot_damage = total_damage_to_mob(hits)
		if tot_damage <= 0 then
			goto continue_drain
		end
		local mob_max_hp = 0
		local mwz = mob:wz()
		if mwz ~= nil and mwz.max_hp ~= nil then
			mob_max_hp = tonumber(mwz.max_hp) or 0
		end
		local raw = math.floor(tot_damage * pct / 100.0)
		local heal = math.min(mob_max_hp, math.min(raw, cap_half))
		if heal > 0 then
			total_heal = total_heal + heal
		end
		::continue_drain::
	end
	if total_heal > 0 then
		me:add_hp(total_heal)
	end
end

function for_each_mob_in_skill_area(me, skill, callback)
	if me == nil or skill == nil or callback == nil then
		return
	end
	local map = me:map()
	if map == nil then
		return
	end
	local effect = get_skill_effect(skill)
	if effect == nil then
		return
	end
	local prop = effect.prop or 0
	if prop <= 0 then
		prop = 100
	end
	local mob_limit = skill_effect_mob_limit(effect)
	local pos_x, pos_y = me:position()
	local lt = effect.lt
	local rb = effect.rb
	if lt == nil or rb == nil then
		return
	end
	local min_x = pos_x + math.min(lt.x, rb.x)
	local max_x = pos_x + math.max(lt.x, rb.x)
	local min_y = pos_y + math.min(lt.y, rb.y)
	local max_y = pos_y + math.max(lt.y, rb.y)
	local mobs = map:objects(ObjectType.Mob, { area = { minX = min_x, minY = min_y, maxX = max_x, maxY = max_y } })
	local n = 0
	for _, mob in ipairs(mobs) do
		if mob_limit ~= nil and n >= mob_limit then
			break
		end
		if math.random(1, 100) <= prop then
			if callback(mob) ~= false then
				n = n + 1
			end
		else
			n = n + 1
		end
	end
end

function for_each_character_in_skill_area(me, skill, callback)
	if me == nil or skill == nil or callback == nil then
		return
	end
	local map = me:map()
	if map == nil then
		return
	end
	local effect = get_skill_effect(skill)
	local pos_x, pos_y = me:position()
	local lt, rb
	if effect ~= nil then
		lt = effect.lt
		rb = effect.rb
	end
	if lt == nil or rb == nil then
		lt = { x = -400, y = -350 }
		rb = { x = 400, y = 250 }
	end
	local min_x = pos_x + math.min(lt.x, rb.x)
	local max_x = pos_x + math.max(lt.x, rb.x)
	local min_y = pos_y + math.min(lt.y, rb.y)
	local max_y = pos_y + math.max(lt.y, rb.y)
	local chars = map:objects(ObjectType.Character, { area = { minX = min_x, minY = min_y, maxX = max_x, maxY = max_y } })
	for _, ch in ipairs(chars) do
		if ch ~= nil then
			callback(ch)
		end
	end
end

local poison_pdam_cap = 30000
local poison_denominator_base = 70

function compute_poison_tick_multiplier(me, skill)
	local mul = 1.0
	if me ~= nil then
		mul = mul * element_amp_from_class(me)
	end
	local swz = skill and skill:wz()
	if swz ~= nil and swz.id == Skill.Flamethrower and me ~= nil then
		local boost = me:skill(Skill.ElementBoost)
		if boost ~= nil then
			local be = get_skill_effect(boost)
			local x = 0
			if be ~= nil and be.x ~= nil then
				x = be.x
			end
			if x > 0 then
				mul = mul * (x / 100.0 + 1.0)
			end
		end
	end
	return mul
end

function compute_poison_tick_damage(skill, mob, multiplier)
	if mob == nil or skill == nil then
		return 0
	end
	local effect = get_skill_effect(skill)
	if effect == nil then
		return 0
	end
	local mwz = mob:wz()
	if mwz == nil then
		return 0
	end
	local max_hp = mwz.max_hp or 0
	if max_hp <= 0 then
		return 0
	end
	local skill_level = skill:level() or 1
	local denom = poison_denominator_base - skill_level
	if denom <= 0 then
		denom = 1
	end
	local quotient = math.floor(max_hp / denom)
	local base = quotient + 0.999
	local weak = element_weak_multiplier(skill, mwz)
	local mul = multiplier or 1.0
	if mul <= 0 then
		mul = 1.0
	end
	local raw = base * weak * mul
	local clamped = math.max(1, math.min(raw, poison_pdam_cap))
	return math.floor(clamped)
end

function apply_prob_status_on_skill_hit(me, skill, damages, status)
	if damages == nil or skill == nil or status == nil then
		return
	end
	local effect = get_skill_effect(skill)
	if effect == nil then
		return
	end
	local prop = effect.prop or 0
	if prop <= 0 then
		prop = 100
	end
	local duration_ms = effect.time or 0
	if duration_ms <= 0 then
		return
	end
	for mob, hits in pairs(damages) do
		if mob and hits and total_damage_to_mob(hits) > 0 then
			if math.random(1, 100) <= prop then
				local value = 1
				if status == MobStatus.Poison then
					local multiplier = compute_poison_tick_multiplier(me, skill)
					value = compute_poison_tick_damage(skill, mob, multiplier)
				end
				mob:set_status(status, value, duration_ms, skill, me)
			end
		end
	end
end

local PICKPOCKET_SKILL_IDS = { [0] = true }
do
    local function add(id)
        if id ~= nil then
            PICKPOCKET_SKILL_IDS[id] = true
        end
    end
    add(Skill.DoubleStab)
    add(Skill.SavageBlow)
    add(Skill.Assaulter)
    add(Skill.BandOfThieves)
    add(Skill.Showdown4221003)
    add(Skill.BoomerangStep)
end

function handle_combo_attack(me, targets, skill)
    if targets == nil or #targets == 0 then
        return
    end

    local current = me:buff_value(BuffFlag.Combo)
    if current == nil then
        return
    end

    local shout_hero = Skill.Shout
    local shout_dw = Skill.DawnWarriorShout
    if skill ~= nil then
        local wz = skill:wz()
        if wz ~= nil then
            local sid = wz.id
            if sid == shout_hero or sid == shout_dw then
                return
            end
        end
    end

    local combo_skill = nil
    local adv_skill = nil
    if me:class_of(Class.Crusader) then
        combo_skill = Skill.ComboAttack
        adv_skill = Skill.AdvancedCombo
    elseif me:class_of(Class.DawnWarrior3) then
        combo_skill = Skill.DawnWarriorComboAttack
        adv_skill = Skill.DawnWarriorAdvancedCombo
    else
        return
    end

    local combo = me:skill(combo_skill)
    local adv = me:skill(adv_skill)
    if combo == nil then
        return
    end

    local ceffect_skill = combo
    if adv ~= nil and adv:level() > 0 then
        ceffect_skill = adv
    end

    local ceffect = nil
    local wz = ceffect_skill:wz()
    if wz ~= nil and wz.effects ~= nil then
        ceffect = wz.effects[ceffect_skill:level()]
    end
    if ceffect == nil then
        return
    end

    local max_orbs = (ceffect.x or 0) + 1
    if max_orbs <= 0 then
        return
    end
    if current >= max_orbs then
        return
    end

    local new_orbs = current + 1
    if adv ~= nil and adv:level() > 0 then
        local prop = ceffect.prop or 0
        if prop > 0 and math.random(1, 100) <= prop and new_orbs < max_orbs then
            new_orbs = new_orbs + 1
        end
    end

    if new_orbs > max_orbs then
        new_orbs = max_orbs
    end
    me:buff_value(BuffFlag.Combo, new_orbs)
end

function handle_pickpocket(me, skill, damages)
    if damages == nil then
        return
    end
    local maxmeso = me:buff_value(BuffFlag.Pickpocket)
    if maxmeso == nil or maxmeso < 1 then
        return
    end
    local wz = skill and skill:wz()
    local skill_id = (wz and wz.id) or 0
    if not PICKPOCKET_SKILL_IDS[skill_id] then
        return
    end
    local map = me:map()
    if map == nil then
        return
    end
    for mob, hits in pairs(damages) do
        if not mob or not hits then
            goto continue_mob
        end
        for _, amount in ipairs(hits) do
            if not amount or amount <= 0 then
                goto continue_hit
            end
            local meso = math.floor((amount / 12300) * maxmeso)
            if meso < 1 then
                meso = 1
            end
            if meso > maxmeso then
                meso = maxmeso
            end
            if math.random(100) >= 100 then
                goto continue_hit
            end
            local x, y = mob:position()
            local offset = math.random(-20, 20)
            map:spawn_meso(meso, { x + offset, y }, me)
            ::continue_hit::
        end
        ::continue_mob::
    end
end

function damages_to_targets(damages)
    if damages == nil then
        return {}
    end
    local targets = {}
    for mob, _ in pairs(damages) do
        if mob ~= nil then
            targets[#targets + 1] = mob
        end
    end
    return targets
end

function handle_ice_charge_freeze(me, damages)
    if damages == nil then
        return
    end
    if not me:class_of(Class.WhiteKnight) and not me:class_of(Class.Paladin) then
        return
    end
    local buff = me:buff(BuffFlag.WkCharge)
    if buff == nil then
        return
    end
    local wz = buff:wz()
    if wz == nil or wz.effects == nil then
        return
    end
    local sid = wz.id
    if sid ~= Skill.IceChargeSword and sid ~= Skill.BlizzardChargeBw then
        return
    end
    local effect = wz.effects[buff:level()]
    if effect == nil then
        return
    end
    local y = effect.y or 0
    local duration_ms = y * 2000
    if duration_ms <= 0 then
        return
    end
    for mob, hits in pairs(damages) do
        if mob and hits and total_damage_to_mob(hits) > 0 then
            mob:set_status(MobStatus.Freeze, 1, duration_ms, buff)
        end
    end
end

function handle_mortal_blow(me, damages)
    if damages == nil then
        return
    end

	local skill = nil
	if me:class_of(Class.Ranger) then
		skill = me:skill(Skill.MortalBlow)
	elseif me:class_of(Class.Crossbowman) then
		skill = me:skill(Skill.MortalBlow3210001)
	else
		return
	end
	if skill == nil then
		return
	end
	local wz = skill:wz()
	if wz == nil then
		return
	end
	local effect = wz.effects[skill:level()]
	if effect == nil then
		return
	end

	for mob, hits in pairs(damages) do
		local mob_hp_percent = mob:hp() / mob:max_hp()
		if mob_hp_percent <= effect.x / 100 then
			if math.random(1, 100) <= effect.y then
				for i = 1, #hits do
					hits[i] = mob:max_hp()
				end
			end
		end
	end
end