-- Map script (onFirstUserEnter StageMsg_juliet): 단계 안내 문구

local HINTS = {
    [926110000] = "연구실로 가는 문을 여는 수상한 스위치를 찾아주세요.",
    [926110001] = "어둠 속의 몬스터를 모두 물리쳐 주세요.",
    [926110100] = "비커에 액체를 가득 채워주세요.",
    [926110200] = "각 문을 통해 연구 자료를 찾아와 주세요!",
    [926110203] = "모든 몬스터를 없애주세요!",
    [926110300] = "연구실 꼭대기로 가는 길을 찾으세요!",
    [926110401] = "제 사랑 로미오를 보호해주세요!",
}

return {
    on_activate = function(me, map)
        local hint = HINTS[map:wz():id()]
        if hint ~= nil then
            map:start_weather(5120022, hint)
        end
    end,
}
