-- NPC name (String.wz/Npc.img.xml): 프란츠

local salon = require("script/lib/salon")

local CONFIG = {
    coupon = 5152005,
    intro = "어서오시게. 오르비스 성형외과에 온걸 환영하네. ",
    outro = "을 가져온다면 자네의 얼굴을 멋지게 바꿔주겠네.",
    done = "자아~ 다 됐다네. 어떤가? 멋지지 않은가? 하하하. 나중에 성형이 또 하고 싶다면 이곳을 다시 찾아주게나!",
}

return {
    on_click = function(me, npc)
        salon.face(me, npc, CONFIG)
    end,
}
