-- NPC name (String.wz/Npc.img.xml): 앨

local salon = require("script/lib/salon")

local CONFIG = {
    coupon = 5152007,
    intro = "어서오세요! 루디브리엄 성형외과에 잘 오셨어요~ ",
    outro = "을 가져오신다면 손님의 얼굴을 멋지고 아름답게 바꿔드릴 수 있답니다!",
    prompt = "원하시는 얼굴을 선택해주세요~",
    no_coupon = "쿠폰이 없으시다면 얼굴 성형을 해드릴 수 없답니다.",
    done = "시술이 잘 된것 같아요. 얼굴이 마음에 드셨으면 좋겠어요~",
}

return {
    on_click = function(me, npc)
        salon.face(me, npc, CONFIG)
    end,
}
