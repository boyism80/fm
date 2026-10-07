local function hp_recover_check(me)
    if me:class_of(Class.Warrior) then
        local s = me:skill(Skill.ImprovingHpRecovery)
        if s ~= nil then
            local effect = s:effect()
            if effect.hp > 0 then
                return effect.hp
            end
        end
    end
    if me:class_of(Class.Assassin) then
        local s = me:skill(Skill.Endure4100002)
        if s ~= nil then
            local effect = s:effect()
            if effect.hp > 0 then
                return effect.hp
            end
        end
    end
    if me:class_of(Class.Bandit) then
        local s = me:skill(Skill.Endure4200001)
        if s ~= nil then
            local effect = s:effect()
            if effect.hp > 0 then
                return effect.hp
            end
        end
    end
    return 0
end

local function mp_recover_check(me)
    if me:class_of(Class.Magician) then
        local s = me:skill(Skill.ImprovingMpRecovery2000000)
        if s ~= nil then
            local effect = s:effect()
            if effect.mp > 0 then
                return effect.mp
            end
        end
        return 0
    end
    if me:class_of(Class.Assassin) then
        local s = me:skill(Skill.Endure4100002)
        if s ~= nil then
            local effect = s:effect()
            if effect.mp > 0 then
                return effect.mp
            end
        end
    elseif me:class_of(Class.Bandit) then
        local s = me:skill(Skill.Endure4200001)
        if s ~= nil then
            local effect = s:effect()
            if effect.mp > 0 then
                return effect.mp
            end
        end
    elseif me:class_of(Class.Crusader) then
        local s = me:skill(Skill.ImprovingMpRecovery)
        if s ~= nil then
            local effect = s:effect()
            if effect.mp > 0 then
                return effect.mp
            end
        end
    elseif me:class_of(Class.WhiteKnight) then
        local s = me:skill(Skill.ImprovingMpRecovery1210000)
        if s ~= nil then
            local effect = s:effect()
            if effect.mp > 0 then
                return effect.mp
            end
        end
    elseif me:class_of(Class.DawnWarrior3) then
        local s = me:skill(Skill.ImprovingMpRecoveryCygnus)
        if s ~= nil then
            local effect = s:effect()
            if effect ~= nil and effect.mp > 0 then
                return effect.mp
            end
        end
    end
    return 0
end

local function ap_to_hp_base(me)
    if me:class_of(Class.Beginner) or me:class_of(Class.Noblesse) or me:class_of(Class.Legend) then
        return math.random(8, 12)
    end
    if me:class_of(Class.Warrior) then
        return math.random(12, 20)
    end
    if me:class_of(Class.Magician) then
        return math.random(6, 11)
    end
    if me:class_of(Class.Bowman) or me:class_of(Class.Thief) then
        return math.random(14, 18)
    end
    return math.random(50, 100)
end

local function ap_to_mp_base(me)
    if me:class_of(Class.Beginner) or me:class_of(Class.Noblesse) or me:class_of(Class.Legend) then
        return math.random(6, 8)
    end
    if me:class_of(Class.Magician) then
        return math.random(10, 20)
    end
    if me:class_of(Class.Bowman) or me:class_of(Class.Thief) then
        return math.random(8, 12)
    end
    if me:class_of(Class.Warrior) then
        return math.random(4, 7)
    end
    return math.random(50, 100)
end

local function endure_hp_interval(me)
    if me:class_of(Class.Warrior) then
        local s = me:skill(Skill.Endure)
        if s ~= nil then
            local effect = s:effect()
            if effect.time > 0 then
                return effect.time / 1000
            end
        end
    end
    if me:class_of(Class.Assassin) then
        local s = me:skill(Skill.Endure4100002)
        if s ~= nil then
            local effect = s:effect()
            if effect.time > 0 then
                return effect.time / 1000
            end
        end
    end
    if me:class_of(Class.Bandit) then
        local s = me:skill(Skill.Endure4200001)
        if s ~= nil then
            local effect = s:effect()
            if effect.time > 0 then
                return effect.time / 1000
            end
        end
    end
    return 0
end

local function return_scroll_continent_allowed(from_map_id, to_map_id)
    local from = math.floor(from_map_id / 10000000)
    local to = math.floor(to_map_id / 10000000)
    if to == 60 or from == 61 then
        return true
    end
    if to == 21 or from == 20 then
        return true
    end
    if to == 12 or from == 10 then
        return true
    end
    if to == 10 or from == 12 then
        return true
    end
    return from == to
end

local teleport_stone_blocked_maps = {
    [109050000] = true,
    [180000001] = true,
    [180000002] = true,
    [240060200] = true,
    [240060201] = true,
    [280030000] = true,
    [280030001] = true,
    [280090000] = true,
    [950101010] = true,
    [950101100] = true,
}

local function teleport_stone_blocked(map_id)
    if map_id < 100000000 then
        return true
    end
    if math.floor(map_id / 1000000) % 100 == 9 then
        return true
    end
    if map_id >= 680000210 and map_id <= 680000502 then
        return true
    end
    if map_id >= 270010200 and map_id <= 270030500 then
        return true
    end
    if map_id >= 925020100 and map_id <= 925039999 then
        return true
    end
    return teleport_stone_blocked_maps[map_id] == true
end

local function teleport_stone_region_allowed(item_id, from_map_id, to_map_id, by_name)
    if item_id ~= 5040000 or by_name == false then
        return true
    end
    local region = math.floor(to_map_id / 10000000)
    if region >= 20 and region <= 23 then
        return math.floor(from_map_id / 100000000) == math.floor(to_map_id / 100000000)
    end
    if region >= 24 and region <= 26 then
        return math.floor(from_map_id / 10000000) == region
    end
    return true
end

return {
	can_use_return_scroll = function(me, item_id, move_to)
		if me:hp() <= 0 then
		    return false
		end
		local map = me:map()
		if map == nil then
		    return false
		end
		local wz = map:wz()
		if wz == nil or wz:blocks_potion_use() then
		    return false
		end
		local current = map:template_id()
		if current == 749040100 then
		    return false
		end
		if move_to == current and item_id ~= 2031010 and item_id ~= 2030021 then
		    return false
		end
		if move_to == 999999999 then
		    return true
		end
		return return_scroll_continent_allowed(current, move_to)
	end,

	can_register_teleport_stone = function(me, map_id)
		return teleport_stone_blocked(map_id) == false
	end,

	can_teleport_stone = function(me, item_id, from_map_id, to_map_id, by_name)
		if teleport_stone_blocked(from_map_id) or teleport_stone_blocked(to_map_id) then
		    return false
		end
		return teleport_stone_region_allowed(item_id, from_map_id, to_map_id, by_name)
	end,


	get_heal_over_time_cap = function(me)
		local recovery_rate = 1.0
		local map = me:map()
		if map ~= nil then
		    recovery_rate = map:recovery_rate()
		end
		local stance_mult = me:stance_of(Stance.Sit) and 1.5 or 1.0
		local chair = me:chair()

		local hp_check = hp_recover_check(me) + 10 * recovery_rate
		hp_check = hp_check * stance_mult
		if chair >= 3000000 then
		    if chair == 3010000 then
		        hp_check = hp_check + 50
		    elseif chair == 3010001 then
		        hp_check = hp_check + 35
		    elseif chair == 3010007 then
		        hp_check = hp_check + 60
		    elseif chair == 3010009 then
		        hp_check = hp_check + 20
		    end
		end

		local mp_check = mp_recover_check(me) + 3 * recovery_rate
		mp_check = mp_check * stance_mult
		if chair >= 3000000 then
		    if chair == 3010008 then
		        mp_check = mp_check + 60
		    elseif chair == 3010009 then
		        mp_check = mp_check + 20
		    end
		end

		return {
		    max_hp = math.floor(hp_check + 0.5),
		    max_mp = math.floor(mp_check + 0.5),
		    endure_hp_interval = endure_hp_interval(me),
		}
	end,

	get_ap_to_hp = function(me)
		local base = ap_to_hp_base(me)
		local bonus = 0
		if me:class_of(Class.Warrior) then
		    local s = me:skill(Skill.ImprovingMaxhpIncrease)
		    if s ~= nil then
		        local effect = s:effect()
		        bonus = effect.y
		    end
		elseif me:class_of(Class.Pirate) then
		    local s = me:skill(Skill.HpIncrease)
		    if s ~= nil then
		        local effect = s:effect()
		        if effect ~= nil then
		            bonus = effect.y
		        end
		    end
		elseif me:class_of(Class.ThunderBreaker1) then
		    local s = me:skill(Skill.HpIncreaseCygnus)
		    if s ~= nil then
		        local effect = s:effect()
		        if effect ~= nil then
		            bonus = effect.y
		        end
		    end
		elseif me:class_of(Class.DawnWarrior1) then
		    local s = me:skill(Skill.ImprovingMaxhpIncreaseCygnus)
		    if s ~= nil then
		        local effect = s:effect()
		        if effect ~= nil then
		            bonus = effect.y
		        end
		    end
		end
		return base + bonus
	end,

	get_ap_to_mp = function(me)
		local base = ap_to_mp_base(me)

		local bonus = 0
		if me:class_of(Class.Magician) then
		    local s = me:skill(Skill.ImprovingMaxMpIncrease)
		    if s ~= nil then
		        local effect = s:effect()
		        bonus = effect.y
		    end
		elseif me:class_of(Class.BlazeWizard1) then
		    local s = me:skill(Skill.ImprovingMaxMpIncreaseCygnus)
		    if s ~= nil then
		        local effect = s:effect()
		        if effect ~= nil then
		            bonus = effect.y
		        end
		    end
		end
		return base + bonus
	end
}
