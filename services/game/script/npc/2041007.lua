-- NPC name (String.wz/Npc.img.xml): 미유

local salon = require("script/lib/salon")

local CONFIG = {
    hair_coupon = 5150007,
    color_coupon = 5151007,
    intro = "안녕하세요~ 루디브리엄 헤어숍입니다! 어머나~ 머리 스타일이 너무 촌스러우신걸요? 여기서 머리 손질을 좀 받으셔야겠어요~ ",
    outro = "아이템을 가져와야 해요~",
    male = { 30030, 30020, 30000, 30510, 30340, 30710, 30300, 30050, 30160, 30190, 30280, 30240, 30150, 30650 },
    female = { 31040, 31050, 31000, 31520, 31460, 31290, 31280, 31270, 31230, 31160, 31120, 31150, 31010, 31030, 31650 },
    style_prompt = "원하시는 헤어 스타일을 선택해 주세요.",
    colors = { 0, 2, 3, 4, 5 },
    color_prompt = "헤어 컬러를 변경하시고 싶으신가요? 원하시는 컬러를 선택해주세요.",
    no_coupon = "죄송하지만 쿠폰을 가져오시지 않으면 머리 손질을 해드릴 수 없답니다.",
    done = "와아~ 너무 멋져요! 원하시는 헤어 스타일이 또 생기실 경우 다시 찾아와 주세요~",
}

return {
    on_click = function(me, npc)
        salon.hair(me, npc, CONFIG)
    end,
}
