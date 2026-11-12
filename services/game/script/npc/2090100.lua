-- NPC name (String.wz/Npc.img.xml): 루오 할아범

local salon = require("script/lib/salon")

local HAIR_COUPON = 5150025
local COLOR_COUPON = 5151020

return {
    on_click = function(me, npc)
        local sel = me:dialog_list(npc,
            "홀홀~ 이 늙은이의 머리 손질 솜씨를 보러 온겐가? #b#i" ..
            HAIR_COUPON ..
            "# #t" ..
            HAIR_COUPON .. "##k 또는 #b#i" .. COLOR_COUPON .. "# #t" .. COLOR_COUPON .. "##k 만 있다면 왕년의 실력을 한껏 뽐내보겠네!", {
                "머리 스타일 바꾸기",
                "머리 색깔 염색하기",
            })
        if sel == nil then
            return
        end

        local coupon = HAIR_COUPON
        local styles = nil
        local prompt = nil
        if sel == 1 then
            local color = me:hair() % 10
            if me:gender() == 0 then
                styles = { 30030, 30020, 30000, 30220, 30460, 30490, 30330, 30420, 30240, 30310, 30180, 30150 }
            else
                styles = { 31000, 31040, 31050, 31470, 31320, 31540, 31120, 31310, 31140, 31280, 31490, 31030, 31010 }
            end
            for i = 1, #styles do
                styles[i] = styles[i] + color
            end
            prompt = "어떤 머리로 손질해줄까?"
        else
            local base = math.floor(me:hair() / 10) * 10
            coupon = COLOR_COUPON
            styles = { base, base + 1, base + 3, base + 6, base + 5 }
            prompt = "어떤 머리 색으로 원하나?"
        end

        local style, paid = salon.pick(me, npc, prompt, styles, coupon)
        if style == nil then
            return
        end
        if not paid then
            me:dialog(npc, "미안하지만 쿠폰 없이는 머리를 손질해줄 수 없다네.")
            return
        end
        me:hair(style)
        me:dialog(npc, "머리 손질이 끝났다네!")
    end,
}
