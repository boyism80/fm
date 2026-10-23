return {
    on_enter = function(me)
        local q3169 = me:quest(3169)
        if me:quest(3141):completed() == false and q3169:started() == false and q3169:completed() == false then
            me:message("길이 막혀있습니다.", Msg.PinkText)
            return
        end
        me:map(211060620, "in00")
    end
}
