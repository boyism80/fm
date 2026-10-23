-- Map script (onFirstUserEnter StageMsg_davy): 단계 안내 문구

local HINTS = {
    [925100000] = "몬스터를 모두 없애주세요!",
    [925100100] = "몬스터를 잡고 해적의 증표를 모아오세요!",
    [925100200] = "몬스터를 모두 없애주세요!",
    [925100300] = "몬스터를 모두 없애주세요!",
    [925100400] = "열쇠로 문을 잠가주세요!",
    [925100500] = "해적왕을 물리쳐주세요!",
}

return {
    on_activate = function(me, map)
        local hint = HINTS[map:wz():id()]
        if hint ~= nil then
            map:start_weather(5120020, hint)
        end
    end,
}
