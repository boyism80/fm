-- Mob skill name (Skill.wz/MobSkill.img.xml): mob skill 127

local function dispel_skill_id(target)
    if target:hidden() then
        return
    end

    local skip_skill_id = {}

    local morph_buff = target:buff(BuffFlag.Morph)
    if morph_buff ~= nil then
        skip_skill_id[morph_buff:wz():id()] = true
    end

    local riding_buff = target:buff(BuffFlag.MonsterRiding)
    if riding_buff ~= nil then
        skip_skill_id[riding_buff:wz():id()] = true
    end

    for _, buff in ipairs(target:buffs(BuffType.Skill)) do
        local skill_id = buff:wz():id()
        if not skip_skill_id[skill_id] then
            target:unbuff(skill_id)
        end
    end
end

return {
    on_mob_skill_choose = function(mob, skill)
        return true
    end,

    on_mob_skill = function(mob, controller, skill)
        for _, target in ipairs(mob:skill_characters(skill, controller)) do
            dispel_skill_id(target)
        end

        return true
    end
}
