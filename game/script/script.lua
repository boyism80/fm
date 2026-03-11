function on_start(me)
	local npc = 9001000
	-- me:dialog_list("안녕하세요", {"hello1", "hello2", "hello3"})
	-- me:dialog_accept('안녕하세요', true)
	-- me:dialog_yes_no('안녕하세요', true, true)
	local text = me:dialog_input(npc,'안녕하세요')
	me:dialog(npc, string.format("반갑습니다. %s님", me:name()))
	for i=1, 10 do
		me:chat(tostring(i))
		sleep(100)
	end
	if text ~= nil then
		me:chat(text)
	else
		me:chat('cancel')
	end
end

function on_script(me)
    local x, y = me:position()
    me:chat(string.format("position: %d, %d", x, y))

    me:class(412)
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
    me:mkitem('후루츠 대거')
    me:mkitem('소환의 돌', 200)
    me:mkitem('화비표창', 2000)
    me:mkitem('메바')
    me:base_dex(128)
    me:base_luk(128)
    me:level(200)
end

local function handle_combo_attack(me, targets, skill)
    if targets == nil or #targets == 0 then
        return
    end

    local current = me:buff_value(BuffFlag.Combo)
    if current == nil then
        return
    end

    local shout_hero = SKILL.SHOUT
    local shout_dw = SKILL.DAWN_WARRIOR_SHOUT
    if skill ~= nil then
        local wz = skill:wz()
        if wz ~= nil then
            local sid = wz.id
            if sid == shout_hero or sid == shout_dw then
                return
            end
        end
    end

    local combo = me:skill(SKILL.COMBO_ATTACK)
    local adv = me:skill(SKILL.ADVANCED_COMBO)
    if combo == nil then
        combo = me:skill(SKILL.DAWN_WARRIOR_COMBO_ATTACK)
        adv = me:skill(SKILL.DAWN_WARRIOR_ADVANCED_COMBO)
    end
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
        if prop > 0 and math.random(100) <= prop and new_orbs < max_orbs then
            new_orbs = new_orbs + 1
        end
    end

    if new_orbs > max_orbs then
        new_orbs = max_orbs
    end
    me:buff_value(BuffFlag.Combo, new_orbs)
end

local PICKPOCKET_SKILL_IDS = {
    [0] = true,
    [SKILL.DOUBLE_STAB] = true,
    [SKILL.SAVAGE_BLOW] = true,
    [SKILL.ASSAULTER] = true,
    [SKILL.BAND_OF_THIEVES] = true,
    [SKILL.SHOWDOWN_4221003] = true,
    [SKILL.BOOMERANG_STEP] = true,
}

local function handle_pickpocket(me, skill, damages)
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

local function damages_to_targets(damages)
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

local ICE_CHARGE_SWORD = 1211005
local BLIZZARD_CHARGE_BW = 1211006
local HERO_JOB = 121
local PALADIN_JOB = 122

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

local function handle_ice_charge_freeze(me, damages)
    if damages == nil then
        return
    end
    local job = me:class()
    if job ~= HERO_JOB and job ~= PALADIN_JOB then
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
    if sid ~= ICE_CHARGE_SWORD and sid ~= BLIZZARD_CHARGE_BW then
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
            mob:set_debuff(Debuff.Freeze, 1, duration_ms, buff)
        end
    end
end

function on_attack(me, skill, damages, attack_info)
    local targets = damages_to_targets(damages)
    handle_combo_attack(me, targets, skill)
    handle_pickpocket(me, skill, damages)
    handle_ice_charge_freeze(me, damages)

    if attack_info and attack_info.ranged and attack_info.consume_slot and attack_info.consume_slot > 0 then
        local weapon = me:equipped(EquipmentPart.Weapon)
        if weapon == nil then
            return
        end
        local wz_weapon = weapon:wz()
        if wz_weapon == nil then
            return
        end
        local weapon_type = wz_weapon:weapon_type()
        if weapon_type == nil then
            return
        end

        me:chat(string.format("consume_slot: %d", attack_info.consume_slot))
        local consume_item = me:item(InventoryType.Use, attack_info.consume_slot)
        if consume_item == nil then
            return
        end
        local wz_consume = consume_item:wz()
        if wz_consume == nil then
            return
        end
        local consume_type = wz_consume:consume_type()
        if consume_type == nil then
            return
        end

        local valid_ammo = (weapon_type == WeaponType.Bow and consume_type == ConsumeType.ArrowBow)
            or (weapon_type == WeaponType.Crossbow and consume_type == ConsumeType.ArrowCrossBow)
            or (weapon_type == WeaponType.Claw and consume_type == ConsumeType.Shuriken)
            or (weapon_type == WeaponType.Gun and consume_type == ConsumeType.Bullet)
        if not valid_ammo then
            return
        end

        local skip_consume = false
        if weapon_type == WeaponType.Bow or weapon_type == WeaponType.Crossbow then
            if me:buff_value(BuffFlag.SoulArrow) ~= nil then
                skip_consume = true
            end
        elseif weapon_type == WeaponType.Claw then
            if me:buff_value(BuffFlag.SpiritClaw) ~= nil then
                skip_consume = true
            end
        end

        if not skip_consume then
            local count = 1
            if skill ~= nil then
                local wz = skill:wz()
                if wz ~= nil and wz.effects ~= nil then
                    local lv = skill:level()
                    local effect = wz.effects[lv]
                    if effect ~= nil then
                        local bc = (effect.bullet_count and effect.bullet_count > 0) and effect.bullet_count or 1
                        local ac = (effect.attack_count and effect.attack_count > 0) and effect.attack_count or 1
                        count = math.max(bc, ac)
                    end
                end
            end

            if me:buff_value(BuffFlag.ShadowPartner) ~= nil then
                count = count * 2
            end

            local ok = me:rmitem(InventoryType.Use, attack_info.consume_slot, count)
            if not ok then
            end
        end
    end
end

local function get_skill_effect_x(skill_entry)
    if skill_entry == nil then
        return nil
    end
    local wz = skill_entry:wz()
    if wz == nil or wz.effects == nil then
        return nil
    end
    local effect = wz.effects[skill_entry:level()]
    if effect == nil then
        return nil
    end
    return effect.x
end

local function handle_magic_guard(me, attacker, skill, damage)
    if damage == nil or damage <= 0 then
        return damage
    end

    if me:buff_value(BuffFlag.MagicGuard) == nil then
        return damage
    end

    local guard_percent = get_skill_effect_x(me:skill(SKILL.MAGIC_GUARD))
    if guard_percent == nil or guard_percent <= 0 then
        guard_percent = get_skill_effect_x(me:skill(SKILL.MAGIC_GUARD_CYGNUS))
    end
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

local function handle_meso_guard(me, attacker, skill, damage)
    if damage == nil or damage <= 0 then
        return damage
    end
    if me:buff_value(BuffFlag.MesoGuard) == nil then
        return damage
    end

    local guard_percent = get_skill_effect_x(me:skill(SKILL.MESO_GUARD))
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
        me:meso(-meso_loss)
    end
    return hp_loss
end

function on_damaged(me, attacker, skill, damage)
    local d = handle_magic_guard(me, attacker, skill, damage)
    return handle_meso_guard(me, attacker, skill, d)
end

function on_equipment_changed(me, part, before, after)
    if part == EquipmentPart.Weapon then
        me:unbuff(BuffFlag.WkCharge)
    end
end