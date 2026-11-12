-- NPC name (String.wz/Npc.img.xml): 프란츠

local salon = require("script/lib/salon")

local COUPON = 5152005

return {
    on_click = function(me, npc)
        if me:dialog(npc, "어서오시게. 오르비스 성형외과에 온걸 환영하네. #b#t" .. COUPON .. "##k 을 가져온다면 자네의 얼굴을 멋지게 바꿔주겠네.", false, true) == false then
            return
        end

        local face, paid = salon.pick(me, npc, "어떤 얼굴로 바꿔보고 싶은가?", salon.faces(me), COUPON)
        if face == nil then
            return
        end
        if not paid then
            me:dialog(npc, "미안하지만 쿠폰 없이는 성형수술을 해 줄수 없다네.")
            return
        end
        me:face(face)
        me:dialog(npc, "자아~ 다 됐다네. 어떤가? 멋지지 않은가? 하하하. 나중에 성형이 또 하고 싶다면 이곳을 다시 찾아주게나!")
    end,
}
