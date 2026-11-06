-- Skill name (String.wz/Skill.img.xml): 힐

local util = require("script/lib/skill")
local combat = require("script/lib/combat")

return {
    on_activated = function(me, skill, params)
        local amount = util.get_heal_recovery_amount(me, skill)
        if amount == 0 then
            return
        end
        combat.for_each_near_party_member(me, skill, function(ch)
            if ch:has_debuff(DebuffFlag.Zombify) then
                ch:add_hp(-amount)
            else
                -- Healing another member pays the caster for the HP actually restored.
                if ch ~= me then
                    local healed = math.min(ch:max_hp() - ch:hp(), amount)
                    if healed > 0 then
                        local max_mp = math.floor(me:max_mp() / 256)
                        me:exp(me:exp() + math.floor(20 * healed / (8 * max_mp + 190)))
                    end
                end
                ch:add_hp(amount)
            end
            if ch ~= me then
                ch:show_skill_effect(skill, SkillEffectType.Affected)
            end
        end)
    end
}
