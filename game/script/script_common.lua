-- Shared helpers for script modules (skill effect, etc.)

function element_amp_from_class(me)
    local amp = 1.0
    if me == nil then
        return amp
    end
    local skill_id
    if me:class_of(Class.FpWizard) then
        skill_id = Skill.ElementAmplification
    elseif me:class_of(Class.IlWizard) then
        skill_id = Skill.ElementAmplification2210001
    elseif me:class_of(Class.BlazeWizard1) then
        skill_id = Skill.ElementAmplificationCygnus
    end
    if skill_id == nil then
        return amp
    end
    local s = me:skill(skill_id)
    if s == nil then
        return amp
    end
    local e = s:effect()
    amp = e.y / 100.0
    return amp
end

function element_weak_multiplier(swz, mwz)
    local weak = 1.0
    if mwz == nil or mwz.elem_resist == nil then
        return weak
    end
    local attr = swz.elem_attr
    if attr == nil or attr == "" then
        return weak
    end
    local letter = string.lower(string.sub(tostring(attr), 1, 1))
    local v = mwz.elem_resist[letter]
    if v == nil then
        return weak
    end
    if v == 1 or v == 2 then
        return 0.5
    end
    if v == 3 then
        return 1.5
    end
    return weak
end
