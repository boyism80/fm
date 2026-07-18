local combat = require("script/lib/combat")
local skill_lib = require("script/lib/skill")

function on_mob_kill(attacker, mobs)
	if attacker == nil or mobs == nil then
		return
	end
	local quest = attacker:quest(29400)
	if quest == nil or not quest:started() then
		return
	end
	local char_level = attacker:level()
	local count = 0
	for _, mob in ipairs(mobs) do
		if mob == nil then
			goto continue
		end
		local wz = mob:wz()
		if wz ~= nil and wz.level ~= nil and wz.level >= char_level then
			count = count + 1
		end
		::continue::
	end
	if count <= 0 then
		return
	end
	local mon = tonumber(quest:record_ex("mon")) or 0
	quest:record_ex("mon", tostring(mon + count))
end

function on_map_enter(me, map)
end

local function unbuff_stationary_summons_for_map(me, map_id)
	if me == nil then
		return
	end
	for _, s in ipairs(me:summons()) do
		if s == nil then
			goto continue
		end
		local sm = s:map()
		if sm == nil then
			goto continue
		end
		local swz = sm:wz()
		if swz == nil or swz.id ~= map_id then
			goto continue
		end
		local mt = s:movement_type()
		if mt ~= SummonMovementType.Stationary
			and mt ~= SummonMovementType.WalkStationary
			and mt ~= SummonMovementType.CircleStationary then
			goto continue
		end
		me:unbuff(s:skill_id())
		::continue::
	end
end

function on_map_leave(me, map)
	if me == nil or map == nil then
		return
	end
	local mob = me:homing()
	if mob ~= nil then
		mob:homing(me, nil)
	end
	local mwz = map:wz()
	if mwz == nil then
		return
	end
	unbuff_stationary_summons_for_map(me, mwz.id)
end

function on_logout(me)
	if me == nil then
		return
	end
	me:unbuff(Skill.MysticDoor)
end

function on_start(me)
	local npc = 9001000
	local text = me:dialog_input(npc,'안녕하세요')
	me:dialog(npc, string.format("반갑습니다. %s님", me:name()))
	for i=1, 10 do
		sleep(100)
	end
	if text ~= nil then
	else
	end
end

function on_script(me)
    local x, y = me:position()

    me:class(232)
    local wz_skills = class_learnable_skill_wzs(me:class())
    for _, wz in pairs(wz_skills) do
        local skill = me:skill(wz.id)
        if skill == nil then
            skill = me:add_skill(wz.id)
        end
        if skill ~= nil then
            skill:level(wz.max_level, wz.master_level)
        end
    end
    me:max_hp(20000)
    me:hp(me:max_hp())
    me:max_mp(20000)
    me:mp(me:max_mp())
    local weapon = me:mkitem('우드완드')
    if weapon ~= nil then
        me:equip(weapon)
    end
    me:mkitem('수비표창', 200)
    me:mkitem('불릿', 2000)
    me:mkitem('마법의돌', 200)
    me:base_str(80)
    me:base_dex(300)
    me:base_int(900)
    me:base_luk(900)
    me:level(200)
    me:map('커닝시티')
end

function on_damaged(me, attacker, skill, damage, params)
    local d = skill_lib.absorb_magic_guard(me, attacker, skill, damage)
    d = skill_lib.absorb_meso_guard(me, attacker, skill, d)
    d = combat.reflect_incoming_damage(me, attacker, skill, d, params)
    return d
end

function on_poison(mist, targets)
	if mist == nil or targets == nil then
		return
	end
	local effect = mist:effect()
	if effect == nil then
		return
	end
	local prop = effect.prop or 0
	if prop <= 0 then
		prop = 100
	end
	local time = effect.time or 0
	if time <= 0 then
		return
	end

	if mist:from_mob() then
		local skill_id = effect.skill_id
		local skill_level = effect.level or mist:level()
		local x = 30
		for _, ch in ipairs(targets) do
			if ch == nil then
				goto continue_char
			end
			if ch:has_debuff(DebuffFlag.Poison) then
				goto continue_char
			end
			if math.random(1, 100) > prop then
				goto continue_char
			end
			ch:debuff(DebuffFlag.Poison, time, x, skill_id, skill_level)
			ch:add_hp(-x)
			::continue_char::
		end
	else
		local wz = mist:wz()
		local level = mist:level()
		local multiplier = mist:poison_tick_multiplier() or 1.0
		if multiplier <= 0 then
			multiplier = 1.0
		end
		for _, mob in ipairs(targets) do
			if mob == nil then
				goto continue_mob
			end
			if mob:has_buff(MobBuff.Poison) then
				goto continue_mob
			end
			if math.random(1, 100) > prop then
				goto continue_mob
			end
			local value = combat.compute_poison_tick_damage_wz_level(wz, level, mob, multiplier)
			if value <= 0 then
				goto continue_mob
			end
			mob:buff(MobBuff.Poison, value, time, mist, mist:causer())
			::continue_mob::
		end
	end
end

function on_blocked(me, attacker)
    if me == nil or attacker == nil then
        return
    end
    local class = me:class()
    local block_skill_id
    if class == Class.Hero then
        block_skill_id = Skill.Guardian
    elseif class == Class.Paladin then
        block_skill_id = Skill.Guardian1220006
    else
        return
    end
    local skill = me:skill(block_skill_id)
    if skill == nil then
        return
    end
    local effect = skill:effect()
    local prop = effect.prop or 0
    local time = effect.time or 0
    if prop <= 0 or time <= 0 then
        return
    end
    if math.random(1, 100) <= math.min(100, prop) then
        attacker:buff(MobBuff.Stun, 1, time, skill, me)
    end
end

function on_equipment_changed(me, part, before, after)
    if part == EquipmentPart.Weapon then
        me:unbuff(BuffFlag.WkCharge)
    end
end

function on_level_up(me, old_level, new_level)
    if me == nil or new_level <= old_level then
        return
    end
    local diff = new_level - old_level
    me:ability_point(me:ability_point() + diff * 5, false)
    if not me:class_of(Class.Beginner) then
        me:skill_point(me:skill_point() + diff * 3, false)
    end
    local total_hp = 0
    local total_mp = 0
    for level = old_level + 1, new_level do
        total_hp = total_hp + (20 + level * 2)
        total_mp = total_mp + (10 + level)
    end
    local bonus_hp = 0
    local bonus_mp = 0
    if me:class_of(Class.Warrior) then
        local s = me:skill(Skill.ImprovingMaxhpIncrease)
        if s ~= nil then
            local effect = s:effect()
            bonus_hp = bonus_hp + diff * effect.x
        end
    elseif me:class_of(Class.Pirate) then
        local s = me:skill(Skill.HpIncrease)
        if s ~= nil then
            local effect = s:effect()
            if effect ~= nil then
                bonus_hp = bonus_hp + diff * effect.x
            end
        end
    elseif me:class_of(Class.ThunderBreaker1) then
        local s = me:skill(Skill.HpIncreaseCygnus)
        if s ~= nil then
            local effect = s:effect()
            if effect ~= nil then
                bonus_hp = bonus_hp + diff * effect.x
            end
        end
    elseif me:class_of(Class.DawnWarrior1) then
        local s = me:skill(Skill.ImprovingMaxhpIncreaseCygnus)
        if s ~= nil then
            local effect = s:effect()
            if effect ~= nil then
                bonus_hp = bonus_hp + diff * effect.x
            end
        end
    elseif me:class_of(Class.Magician) then
        local s = me:skill(Skill.ImprovingMaxMpIncrease)
        if s ~= nil then
            local effect = s:effect()
            bonus_mp = bonus_mp + diff * effect.x
        end
    elseif me:class_of(Class.BlazeWizard1) then
        local s = me:skill(Skill.ImprovingMaxMpIncreaseCygnus)
        if s ~= nil then
            local effect = s:effect()
            if effect ~= nil then
                bonus_mp = bonus_mp + diff * effect.x
            end
        end
    end
    me:base_hp(me:base_hp() + total_hp + bonus_hp, false)
    me:base_mp(me:base_mp() + total_mp + bonus_mp, false)
    me:hp(me:max_hp(), false)
    me:mp(me:max_mp(), false)
    me:update_stats({
        STAT.Level,
        STAT.Exp,
        STAT.MaxHp,
        STAT.MaxMp,
        STAT.Hp,
        STAT.Mp,
        STAT.AvailableAP,
        STAT.AvailableSP,
    })
end
