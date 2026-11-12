-- NPC name (String.wz/Npc.img.xml): 허치희

local salon = require("script/lib/salon")

local CONFIG = {
    coupon = 5152024,
    intro = "시술을 받으러 여기까지 온건가요? ",
    outro = "을 가져온다면 아주 멋진 얼굴로 바꿔드릴수 있답니다.",
    done = "자아~ 다 됐다네. 어떤가? 맘에 들지? 다음에 또 이곳을 찾아주게나~",
}

return {
    on_click = function(me, npc)
        salon.face(me, npc, CONFIG)
    end,
}
