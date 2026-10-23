-- Map script (onFirstUserEnter boss_Ravana): 라바나 방 알림

return {
    on_activate = function(me, map)
        map:message("라바나가 나타났습니다!")
    end,
}
