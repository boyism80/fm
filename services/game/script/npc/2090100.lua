-- NPC name (String.wz/Npc.img.xml): 루오 할아범

local salon = require("script/lib/salon")

local CONFIG = {
    hair_coupon = 5150025,
    color_coupon = 5151020,
    intro = "홀홀~ 이 늙은이의 머리 손질 솜씨를 보러 온겐가? ",
    outro = "만 있다면 왕년의 실력을 한껏 뽐내보겠네!",
    male = { 30030, 30020, 30000, 30220, 30460, 30490, 30330, 30420, 30240, 30310, 30180, 30150 },
    female = { 31000, 31040, 31050, 31470, 31320, 31540, 31120, 31310, 31140, 31280, 31490, 31030, 31010 },
    style_prompt = "어떤 머리로 손질해줄까?",
    colors = { 0, 1, 3, 6, 5 },
    color_prompt = "어떤 머리 색으로 원하나?",
    no_coupon = "미안하지만 쿠폰 없이는 머리를 손질해줄 수 없다네.",
    done = "머리 손질이 끝났다네!",
}

return {
    on_click = function(me, npc)
        salon.hair(me, npc, CONFIG)
    end,
}
