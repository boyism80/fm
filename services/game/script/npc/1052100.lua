-- NPC name (String.wz/Npc.img.xml): 돈 지오바네

local salon = require("script/lib/salon")

local CONFIG = {
    hair_coupon = 5150003,
    color_coupon = 5151003,
    intro = "커닝시티 헤어샵에 온걸 환영하네! 하하하! 나는 이 헤어샵의 원장, 돈 지오바네 라고 한다네! ",
    outro = "을 가져온다면 머리를 아주 아름~답게 손질해주지!",
    male = { 30030, 30020, 30000, 30360, 30370, 30450, 30350, 30190, 30180, 30160, 30130, 30110, 30050, 30040 },
    female = { 31050, 31040, 31000, 31010, 31510, 31170, 31180, 31330, 31140, 31130, 31120, 31090, 31060, 31020, 31010 },
    style_prompt = "어떤 머리 스타일을 원하나?",
    colors = { 0, 2, 3, 7, 5 },
    color_prompt = "어떤 머리 색깔을 원하나?",
    no_coupon = "미안하지만 쿠폰 없이는 머리를 손질해줄 수 없다네.",
    done = "머리 손질이 끝났다네! 후후. 정말 멋지군! 역시 내 솜씨는 훌륭하다니까!",
}

return {
    on_click = function(me, npc)
        salon.hair(me, npc, CONFIG)
    end,
}
