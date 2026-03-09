-- Skill name (String.wz/Skill.img.xml): 콤보 어택

function on_activated(me, skill)
    local previous_value = me:buff_value(BuffFlag.Combo)
    me:add_buff(skill, BuffFlag.Combo, 1)
    if previous_value ~= nil then
        me:buff_value(BuffFlag.Combo, previous_value)
    end
end

function on_deactivated(me, skill)
end