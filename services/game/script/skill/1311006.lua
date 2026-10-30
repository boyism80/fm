-- Skill name (String.wz/Skill.img.xml): 드래곤 로어

-- MP 16/24/30 소비, HP의 59%~30% 소비, 15마리 공격, 발동 시 자기 스턴 (4/3/2초). HP 50% 이상일 때만 사용 가능.

return {
    on_activating = function(me, skill, params)
        return me:hp() * 100 >= me:max_hp() * 50
    end,

    on_activated = function(me, skill, params)
        local effect = skill:effect()
        if effect == nil then
            return
        end
        local max_hp = me:max_hp()
        if effect.x > 0 then
            local hp_loss = math.floor(max_hp * effect.x / 100)
            if hp_loss > 0 then
                me:hp(math.max(me:hp() - hp_loss, 1), true)
            end
        end
        if effect.y > 0 then
            me:debuff(DebuffFlag.Stun, effect.y * 1000)
        end
    end
}
