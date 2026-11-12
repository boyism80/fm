-- NPC name (String.wz/Npc.img.xml): 덴마

local salon = require("script/lib/salon")

local COUPON = 5152001

return {
    on_click = function(me, npc)
        if me:dialog(npc, "어서오게. 헤네시스 성형외과에 온걸 환영하네. #b#t" .. COUPON .. "##k 을 가져온다면 자네의 얼굴을 다른 얼굴로 바꿔줄수 있다네.", false, true) == false then
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
        me:dialog(npc, "자아~ 다 됐다네. 어떤가? 맘에 들지? 다음에 또 이 덴마 님을 찾아주게나~")
    end,
}
