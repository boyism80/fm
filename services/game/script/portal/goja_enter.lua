local GROUP_NAME = "gojarani"

return {
    on_enter = function(me)
        local group = state_machine(GROUP_NAME)
        if group == nil then
            me:message("지금은 입장할 수 없습니다.", Msg.PinkText)
            return
        end
        me:play_portal_sound()
        local sm, err = group:start_solo(me)
        if sm == nil then
            if err == StateMachineError.Full then
                me:message("이미 다른 사람이 안에 있습니다.", Msg.PinkText)
            else
                me:message("지금은 입장할 수 없습니다.", Msg.PinkText)
                if err ~= nil then
                    log("gojarani start_solo:", err)
                end
            end
        end
    end,
}
