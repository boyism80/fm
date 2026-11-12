-- NPC name (String.wz/Npc.img.xml): 나탈리

local salon = require("script/lib/salon")

local CONFIG = {
    hair_coupon = 5150001,
    color_coupon = 5151001,
    intro = "어디가도 멋진 헤어스타일을 만들 곳은 여기 뿐이란 말이지. 서비스를 이용하려면 ",
    outro = "아이템을 가져와야 해요~",
    male = { 30030, 30020, 30000, 30480, 30440, 30470, 30410, 30330, 30310, 30210, 30200, 30150, 30140, 30120, 30060, 30670 },
    female = { 31050, 31040, 31000, 31070, 31080, 31030, 31340, 31480, 31490, 31410, 31310, 31300, 31160, 31100, 31150, 31640 },
    style_prompt = "오호호~ 어떤 헤어스타일을 원하시나요 고객님?",
    colors = { 1, 2, 4, 7 },
    color_prompt = "머리를 염색하고 싶으시군요? 원하시는 색을 골라보세요.",
    no_coupon = "죄송하지만 쿠폰을 가져오시지 않으면 머리 손질을 해드릴 수 없답니다.",
    done = "자~ 다 되었답니다. 어떠세요? 저희 미용실만의 최고의 솜씨를 발휘해 보았답니다.",
}

return {
    on_click = function(me, npc)
        salon.hair(me, npc, CONFIG)
    end,
}
