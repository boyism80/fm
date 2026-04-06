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
	local effect = skill:effect()
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
	if callback == nil then
		return
	end
	local map = me:map()
	if map == nil then
		return
	end
	local effect = skill:effect()
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

function apply_shadow_web_skill(me, skill)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	local duration_ms = effect.time or 0
	if duration_ms <= 0 then
		return
	end
	for_each_mob_in_skill_area(me, skill, function(mob)
		local mwz = mob:wz()
		if mwz ~= nil and mwz.boss then
			return false
		end
		mob:buff(MobBuff.ShadowWeb, 1, duration_ms, skill, me)
		return true
	end)
end

function compute_ninja_ambush_tick_damage(me, skill)
	local effect = skill:effect()
	if effect == nil then
		return 0
	end
	local skill_level = skill:level() or 1
	local dmg_field = tonumber(effect.damage) or 0
	local str = math.floor(tonumber(me:base_str()) or 0) + math.floor(tonumber(me:bonus_str()) or 0)
	local luk = math.floor(tonumber(me:base_luk()) or 0) + math.floor(tonumber(me:bonus_luk()) or 0)
	local raw = (skill_level + 30) * dmg_field * (str + luk) / 2000
	local pdam = math.max(1, math.floor(raw))
	if pdam > 30000 then
		pdam = 30000
	end
	return pdam
end

function apply_ninja_ambush_skill(me, skill)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	local duration_ms = effect.time or 0
	if duration_ms <= 0 then
		return
	end
	for_each_mob_in_skill_area(me, skill, function(mob)
		local value = compute_ninja_ambush_tick_damage(me, skill)
		if value <= 0 then
			return true
		end
		mob:buff(MobBuff.NinjaAmbush, value, duration_ms, skill, me)
		return true
	end)
end

function for_each_character_in_skill_area(me, skill, callback)
	if callback == nil then
		return
	end
	local map = me:map()
	if map == nil then
		return
	end
	local effect = skill:effect()
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
			local be = boost:effect()
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
	return compute_poison_tick_damage_wz_level(skill:wz(), skill:level(), mob, multiplier)
end

function compute_poison_tick_damage_wz_level(wz, level, mob, multiplier)
	if mob == nil then
		return 0
	end
	local mwz = mob:wz()
	local max_hp = mwz.max_hp or 0
	if max_hp <= 0 then
		return 0
	end
	local skill_level = level or 1
	local denom = poison_denominator_base - skill_level
	if denom <= 0 then
		denom = 1
	end
	local quotient = math.floor(max_hp / denom)
	local base = quotient + 0.999
	local weak = element_weak_multiplier(wz, mwz)
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
	local effect = skill:effect()
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
				if status == MobBuff.Poison then
					local multiplier = compute_poison_tick_multiplier(me, skill)
					value = compute_poison_tick_damage(skill, mob, multiplier)
				end
				mob:buff(status, value, duration_ms, skill, me)
			end
		end
	end
end

function apply_showdown_on_attack(me, skill, damages)
	local effect = skill:effect()
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
	local raw_x = effect.x
	if raw_x == nil then
		raw_x = 0
	end
	local val = math.floor(tonumber(raw_x) or 0)
	if val <= 0 then
		return
	end
	for mob, hits in pairs(damages) do
		if total_damage_to_mob(hits) <= 0 then
			goto continue_showdown
		end
		if math.random(1, 100) > prop then
			goto continue_showdown
		end
		mob:buff({
			[MobBuff.Showdown] = val,
			[MobBuff.Mdef] = val,
			[MobBuff.Wdef] = val,
		}, duration_ms, skill, me)
		::continue_showdown::
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

    local ceffect = ceffect_skill:effect()
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
        for _, amount in ipairs(hits) do
            if (amount or 0) <= 0 then
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
    end
end

function damages_to_targets(damages)
    if damages == nil then
        return {}
    end
    local targets = {}
    for mob, _ in pairs(damages) do
        targets[#targets + 1] = mob
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
    if wz == nil then
        return
    end
    local sid = wz.id
    if sid ~= Skill.IceChargeSword and sid ~= Skill.BlizzardChargeBw then
        return
    end
    local effect = buff:effect()
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
            mob:buff(MobBuff.Freeze, 1, duration_ms, buff)
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
	local effect = skill:effect()
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

-- v83 MapleMonster VENOM: matk * (dex + 5*v43) / 49, v43 = floor((rand[0,v55-1] + v55*0.8)), v55 = str+luk (min 1).
function roll_venom_tick_damage(me, venom_skill)
	if venom_skill == nil then
		return 1
	end
	local effect = venom_skill:effect()
	if effect == nil then
		return 1
	end
	local matk = math.floor(tonumber(effect.mad) or 0)
	if matk < 1 then
		matk = 1
	end
	local str = math.floor(tonumber(me:base_str()) or 0) + math.floor(tonumber(me:bonus_str()) or 0)
	local dex = math.floor(tonumber(me:base_dex()) or 0) + math.floor(tonumber(me:bonus_dex()) or 0)
	local luk = math.floor(tonumber(me:base_luk()) or 0) + math.floor(tonumber(me:bonus_luk()) or 0)
	local v55 = str + luk
	if v55 < 1 then
		v55 = 1
	end
	local v56 = v55 * 0.8
	local mod = math.random(0, v55 - 1)
	local v43 = math.floor(mod + v56)
	local v44 = math.floor(matk * (dex + 5 * v43) / 49)
	if v44 < 1 then
		v44 = 1
	end
	if v44 > 30000 then
		v44 = 30000
	end
	return v44
end

-- Passive venom: WZ prop, stack 1..3, tick = sum of rolls capped 30000, immediate hit, buff with stack (Lua drives value like poison mist).
function apply_venom(me, damages, passive_skill_id)
	if damages == nil or passive_skill_id == nil then
		return
	end
	local venom_skill = me:skill(passive_skill_id)
	if venom_skill == nil then
		return
	end
	local venom_level = venom_skill:level()
	if venom_level == nil or venom_level <= 0 then
		return
	end
	local effect = venom_skill:effect()
	if effect == nil then
		return
	end
	local chance = tonumber(effect.prop) or 0
	if chance <= 0 then
		return
	end
	local duration_ms = tonumber(effect.time) or 0
	if duration_ms <= 0 then
		return
	end
	for mob, hits in pairs(damages) do
		if mob == nil or hits == nil then
			goto venom_passive_continue
		end
		if total_damage_to_mob(hits) <= 0 then
			goto venom_passive_continue
		end
		if math.random(0, 99) >= chance then
			goto venom_passive_continue
		end
		local old_stack = math.floor(tonumber(mob:buff_stack(MobBuff.Venom)) or 0)
		if old_stack >= 3 then
			goto venom_passive_continue
		end
		local old_tick = 0
		if mob:has_buff(MobBuff.Venom) then
			old_tick = math.floor(tonumber(mob:buff_value(MobBuff.Venom)) or 0)
		end
		local roll = roll_venom_tick_damage(me, venom_skill)
		local new_tick = old_tick + roll
		if new_tick < 1 then
			new_tick = 1
		end
		if new_tick > 30000 then
			new_tick = 30000
		end
		local new_stack = old_stack + 1
		mob:buff(MobBuff.Venom, new_tick, duration_ms, venom_skill, me, new_stack)
		::venom_passive_continue::
	end
end