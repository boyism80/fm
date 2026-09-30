-- NPC name (String.wz/Npc.img.xml): 에버

local COUPON = 5152006

return {
	on_click = function(me, npc)
		if me:dialog(npc, "싼 요금으로 성형수술을 받고 싶어? 그렇다면 잘 찾아왔군.. #b#t" .. COUPON .. "##k 을 가져온다면 무작위로 얼굴을 바꿔줄게.", false, true) == false then
			return
		end
		if me:dialog_yes_no(npc, "정말로 무작위로 얼굴을 바꿔보고 싶어? 신중하게 결정하라구.") == false then
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
			me:dialog(npc, "쿠폰이 없으면 얼굴을 바꿔줄 수 없어.")
			return
		end
		me:face(styles[math.random(1, #styles)])
		me:dialog(npc, "시술이 끝났어. 마음에 들었으면 좋겠는데.")
	end,
}
