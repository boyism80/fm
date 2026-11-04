-- NPC name (String.wz/Npc.img.xml): 헤네시스 중형택시

local taxi = require("script/lib/taxi")

local CFG = {
    town = "헤네시스",
    prompt = "목적지를 선택해 주세요. 마을마다 요금이 다릅니다.#b",
    prompt_beginner = "저희 택시는 초보자 분들은 90% 할인을 해드립니다, 목적지를 선택해주세요. 마을마다 요금이 다릅니다.#b",
    confirm_tail = "입니다.",
    routes = {
        { 104000000, 800, 80 },
        { 102000000, 1000, 100 },
        { 101000000, 1000, 100 },
        { 103000000, 1200, 120 },
        { 120000000, 800, 80 },
    },
}

return {
    on_click = function(me, npc)
        taxi.run(me, npc, CFG)
    end
}
