-- NPC name (String.wz/Npc.img.xml): 커닝시티 중형택시

local taxi = require("script/lib/taxi")

local CFG = {
    town = "커닝시티",
    prompt = "목적지를 선택해주세요. 지역마다 요금이 다릅니다.#b",
    prompt_beginner = "저희 택시는 초보자 분들에게는 90% 할인된 요금을 받습니다. 목적지를 선택해주세요.#b",
    confirm_tail = ". 입니다.",
    routes = {
        { 104000000, 1000, 100 },
        { 102000000, 800, 80 },
        { 101000000, 1200, 120 },
        { 100000000, 1000, 100 },
        { 120000000, 1200, 80 },
    },
}

return {
    on_click = function(me, npc)
        taxi.run(me, npc, CFG)
    end
}
