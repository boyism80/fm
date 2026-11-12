-- NPC name (String.wz/Npc.img.xml): 덴마

local salon = require("script/lib/salon")

local CONFIG = {
    coupon = 5152001,
    intro = "어서오게. 헤네시스 성형외과에 온걸 환영하네. ",
    outro = "을 가져온다면 자네의 얼굴을 다른 얼굴로 바꿔줄수 있다네.",
    done = "자아~ 다 됐다네. 어떤가? 맘에 들지? 다음에 또 이 덴마 님을 찾아주게나~",
}

return {
    on_click = function(me, npc)
        salon.face(me, npc, CONFIG)
    end,
}
