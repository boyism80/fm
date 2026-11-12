-- NPC name (String.wz/Npc.img.xml): 마즈라

local salon = require("script/lib/salon")

local CONFIG = {
    hair_coupon = 5150027,
    color_coupon = 5151022,
    intro = "원하는 헤어라도 있는가~? 아리안트 미용실에 잘 오셨네~ ",
    outro = "만 있다면 멋진 헤어스타일로 바꿔드리지~",
    male = { 30030, 30020, 30000, 30470, 30490, 30230, 30200, 30260, 30290, 30160, 30050, 30630 },
    female = { 31000, 31040, 31050, 31490, 31130, 31420, 31260, 31110, 31160, 31300, 31610, 31030 },
    style_prompt = "어떤 머리로 손질해주길 원하나~?",
    colors = { 0, 3, 2, 4 },
    color_prompt = "어떤 머리 색으로 원하나?",
    no_coupon = "미안하지만 쿠폰 없이는 머리를 손질해줄 수 없다네.",
    done = "머리 손질이 끝났다네!",
}

return {
    on_click = function(me, npc)
        salon.hair(me, npc, CONFIG)
    end,
}
