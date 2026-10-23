return {
    on_click = function(me, npc)
        local lines = npc:dialogues()
        if #lines == 0 then
            return
        end
        me:dialog(npc, lines[math.random(#lines)])
    end
}
