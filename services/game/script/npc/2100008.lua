-- NPC name (String.wz/Npc.img.xml): 바드르

local salon = require("script/lib/salon")

local COUPON = 5152030

return {
    on_click = function(me, npc)
        if me:dialog(npc, "성형수술이 필요한가? #b#t" .. COUPON .. "##k 을 가져온다면 자네의 얼굴을 바꿔주겠네.", false, true) == false then
            return
        end

        local face, paid = salon.pick(me, npc, "원하는 얼굴을 골라보게.", salon.faces(me), COUPON)
        if face == nil then
            return
        end
        if not paid then
            me:dialog(npc, "미안하지만 쿠폰 없이는 성형수술을 해 줄수 없다네.")
            return
        end
        me:face(face)
        me:dialog(npc, "자, 다 되었다네.")
    end,
}
