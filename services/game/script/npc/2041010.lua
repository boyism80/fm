-- NPC name (String.wz/Npc.img.xml): 앨

local salon = require("script/lib/salon")

local COUPON = 5152007

return {
    on_click = function(me, npc)
        if me:dialog(npc, "어서오세요! 루디브리엄 성형외과에 잘 오셨어요~ #b#t" .. COUPON .. "##k 을 가져오신다면 손님의 얼굴을 멋지고 아름답게 바꿔드릴 수 있답니다!", false, true) == false then
            return
        end

        local face, paid = salon.pick(me, npc, "원하시는 얼굴을 선택해주세요~", salon.faces(me), COUPON)
        if face == nil then
            return
        end
        if not paid then
            me:dialog(npc, "쿠폰이 없으시다면 얼굴 성형을 해드릴 수 없답니다.")
            return
        end
        me:face(face)
        me:dialog(npc, "시술이 잘 된것 같아요. 얼굴이 마음에 드셨으면 좋겠어요~")
    end,
}
