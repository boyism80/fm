-- NPC name (String.wz/Npc.img.xml): 바드르

local salon = require("script/lib/salon")

local CONFIG = {
    coupon = 5152030,
    intro = "성형수술이 필요한가? ",
    outro = "을 가져온다면 자네의 얼굴을 바꿔주겠네.",
    prompt = "원하는 얼굴을 골라보게.",
    done = "자, 다 되었다네.",
}

return {
    on_click = function(me, npc)
        salon.face(me, npc, CONFIG)
    end,
}
