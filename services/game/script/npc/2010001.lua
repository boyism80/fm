-- NPC name (String.wz/Npc.img.xml): 미노

local salon = require("script/lib/salon")

local MESSY_COUPON = 5154000
local MESSY_HAIR = 31240

local CONFIG = {
    hair_coupon = 5150005,
    color_coupon = 5151005,
    intro = "오르비스의 멋진 헤어숍 입니다. 여러분의 헤어 스타일을 책임져 드린답니다~ ",
    outro = "아이템을 가져와야 해요~",
    male = { 30030, 30020, 30000, 30520, 30480, 30490, 30460, 30420, 30340, 30290, 30280, 30270, 30260, 30240, 30230 },
    female = { 31040, 31000, 31050, 31440, 31540, 31420, 31320, 31270, 31260, 31250, 31240, 31230, 31220, 31110, 31030, 31530 },
    style_prompt = "어떤 헤어 스타일을 원하시나요?",
    colors = { 0, 1, 3, 4, 5 },
    color_prompt = "어떤 헤어 색깔을 원하시나요?",
    no_coupon = "죄송하지만 쿠폰을 가져오시지 않으면 저희 헤어숍을 이용하실 수 없답니다.",
    done = "헤어 손질이 모두 끝났답니다. 다음에 또 이용해 주세요.",
}

-- The messy hair coupon is this salon's own, for women only.
local function messy_hair(me, npc)
    if math.floor(me:hair() / 10) * 10 == MESSY_HAIR then
        me:dialog(npc, "이미 머리 손질이 끝나신 것 같은데요~? 다른 헤어가 필요하시다면 다시 제게 말을 걸어주세요~")
        return
    end
    if me:gender() ~= 1 then
        me:dialog(npc, "죄송하지만 산발머리는 #r여성#k분만 이용이 가능해요. 다른 헤어를 원하신다면 다시 제게 말을 걸어주세요~")
        return
    end
    if me:dialog_yes_no(npc, "정말 #i5154000# #t5154000# 아이템을 사용하시겠어요?") == false then
        return
    end
    if me:exchange({ item = { [MESSY_COUPON] = 1 } }) ~= ExchangeResult.OK then
        me:dialog(npc, CONFIG.no_coupon)
        return
    end
    me:hair(MESSY_HAIR + me:hair() % 10)
    me:dialog(npc, CONFIG.done)
end

return {
    on_click = function(me, npc)
        local sel = me:dialog_list(npc, salon.hair_greeting(CONFIG), {
            "머리 스타일 바꾸기",
            "머리 색깔 염색하기",
            "산발머리 쿠폰 사용하기",
        })
        if sel == 1 then
            salon.hair_style(me, npc, CONFIG)
        elseif sel == 2 then
            salon.hair_color(me, npc, CONFIG)
        elseif sel == 3 then
            messy_hair(me, npc)
        end
    end,
}
