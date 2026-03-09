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
    me:max_mp(20000)
end

local function handle_combo_attack(me, targets, skill)
    if targets == nil or #targets == 0 then
        return
    end

    local current = me:buff_value(BuffFlag.Combo)
    if current == nil then
        return
    end

    -- Shout should not increase combo.
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

    -- Resolve combo/advanced-combo pair from learned skills.
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

function on_attack(me, targets, skill)
    handle_combo_attack(me, targets, skill)
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

function on_damaged(me, attacker, skill, damage)
    return handle_magic_guard(me, attacker, skill, damage)
end