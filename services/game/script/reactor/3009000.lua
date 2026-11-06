-- Reactor name (Reactor.wz/3009000.img.xml): 메이플아일랜드 범용리엑터

return {
    on_reactor = function(reactor)
        if reactor:state() ~= 4 then
            return
        end
        local map = reactor:map()
        if map ~= nil then
            map:clear_effect()
        end
    end
}
