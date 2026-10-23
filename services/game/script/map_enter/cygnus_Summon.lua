-- Map script (onFirstUserEnter cygnus_Summon): 시그너스 방 첫 입장 문구

return {
    on_activate = function(me, map)
        map:start_weather(5120043, "손님이 찾아온 게 얼마 만인지... 살아서 돌아간 손님은 더욱 오래전 일이지.")
    end,
}
