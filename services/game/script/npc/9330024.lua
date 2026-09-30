-- NPC name (String.wz/Npc.img.xml): 곽이랑

local COUPON = 5152023

return {
	on_click = function(me, npc)
		if me:dialog(npc, "음.. 성형수술.. 받으러 오신거죠? #b#t" .. COUPON .. "##k 을 가져오시면 무작위로 얼굴을 바꿔줄 수 있어요..", false, true) == false then
			return
		end
		if me:dialog_yes_no(npc, "정말로 무작위로 얼굴을 바꿔보고 싶은가요..? 저한테 책임 전가 하시면 안돼요..?") == false then
			return
		end

		local face = me:face()
		local styles = nil
		if me:gender() == 0 then
			styles = { 20000, 20001, 20002, 20003, 20004, 20005, 20006, 20007, 20008, 20012, 20014, 20016, 20020, 20017 }
		else
			styles = { 21000, 21001, 21002, 21003, 21004, 21005, 21006, 21007, 21008, 21012, 21014, 21016, 21020, 21017 }
		end
		for i = 1, #styles do
			styles[i] = styles[i] + face % 1000 - face % 100
		end

		if me:exchange({ item = { [COUPON] = 1 } }) ~= ExchangeResult.OK then
			me:dialog(npc, "쿠폰.. 어디에..?")
			return
		end
		me:face(styles[math.random(1, #styles)])
		me:dialog(npc, "끝..났어요")
	end,
}
