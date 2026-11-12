-- NPC name (String.wz/Npc.img.xml): 니키

local salon = require("script/lib/salon")

local CONFIG = {
    hair_coupon = 5150022,
    color_coupon = 5151015,
    intro = "하하! 야시장의 멋진 헤어샵에 잘 오셨습니다~ ",
    outro = "만 있다면 멋진 헤어스타일로 바꿔드리지~",
    male = { 30030, 30020, 30000, 30560, 30420, 30310, 30270, 30260, 30240, 30230, 30200, 30150 },
    female = { 31000, 31040, 31050, 31560, 31320, 31240, 31230, 31150, 31140, 31110, 31100, 31030, 31010 },
    style_prompt = "어떤 머리로 손질해주길 원하나요~?",
    colors = { 0, 3, 4, 7, 5 },
    color_prompt = "어떤 머리 색으로 원하나요~?",
    no_coupon = "미안하지만 쿠폰 없이는 머리를 손질해줄 수 없습니다.",
    done = "머리 손질이 끝났답니다~",
}

return {
    on_click = function(me, npc)
        salon.hair(me, npc, CONFIG)
    end,
}
