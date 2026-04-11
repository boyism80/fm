-- Skill name (String.wz/Skill.img.xml): 헤이스트

function on_activated_4101004(me, skill, params)
    local effect = skill:effect()
    if effect == nil then
        return
    end

    me:buff(skill, {
        [BuffFlag.Speed] = effect.speed,
        [BuffFlag.Jump] = effect.jump,
    })
end

-- TODO: Party buff : when party system exists, apply same buff to party members
-- in range (lt/rb bounding box) with primary=false; each needs giveBuff-equivalent + registerEffect
-- + showOwnBuffEffect/showBuffeffect. Not just buff on other Character without running
-- on_activated per target / duration sync.

