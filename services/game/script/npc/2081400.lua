-- NPC name (String.wz/Npc.img.xml): 헬린

local fourth_class = require("script/lib/fourth_class")

local INTRO = "어둠 속에서 항상 빛을 찾으며 빠르고 신속하게 활주하는 직업. 도적이 되어볼 생각은 없습니까?"

local ADVANCED = {
    [Class.Nightlord] = "당신은 얼마 전에 #b나이트로드#k가 되었던 #r#h ##k이군요, 수행은 잘 되시나요? 더 강한 도적이 되기 위해서 수련함으로 인해 더 정진하길 바랍니다..",
    [Class.Shadower] = "당신은 얼마 전에 #b섀도어#k가 되었던 #r#h ##k이군요, 수행은 잘 되시나요? 더 강한 도적이 되기 위해서 수련함으로 인해 더 정진하길 바랍니다..",
}

local ADVANCE = {
    [Class.Hermit] = {
        class = Class.Nightlord,
        option = "#r#h # : #b나이트로드로 전직하고 싶습니다.",
        done =
        "최고의 도적 #b나이트로드#k가 되신 것을 축하드립니다. 선물로 당신에게 약간의 SP와 AP와 주었으니 확인해 보세요. 이로써 당신은 최고의 강한 도적이 되었지만, 앞으로 더 진보하기 위해서는 끊임없는 수련이 필요합니다...",
    },
    [Class.ChiefBandit] = {
        class = Class.Shadower,
        option = "#r#h # : #b섀도어로 전직하고 싶습니다.",
        done =
        "최고의 도적 #b섀도어#k가 되신 것을 축하드립니다. 선물로 당신에게 약간의 SP와 AP와 주었으니 확인해 보세요. 이로써 당신은 최고의 강한 도적이 되었지만, 앞으로 더 진보하기 위해서는 끊임없는 수련이 필요합니다... ",
    },
}

local CFG = {
    quest = 6934,
    advanced = ADVANCED,
    advance = ADVANCE,
    text = {
        intro = INTRO,
        choose = "어느 클레스로 전직하기를 원하나요?\r\n그대의 영웅으로써의 자질은 증명되었습니다.\r\n",
        wait = "그대는 그대의 영웅으로서의 자질을 증명했소. \r\n이제 남은 것은 도적 궁극의 길로 향하는 것 뿐.\r\n전직을 할 준비가 되면 다시 말을 걸어주시오.",
        sp_left = "흠... 당신은 너무 많은 #bSP#k를 갖고 있는 것 같소. 최소한 레벨 120 이전에 얻은 스킬포인트를 모두 사용하셔야 4차 전직을 하실 수 있소.",
        no_space = "흠.. 소비 인벤토리 공간이 부족한 건 아닌가? 다시 한번 확인해 보게나.",
        farewell = "난 그대가 더욱 더 수련하여 지상 최고의 도적이 되는 것을 뒤에서 지켜봐주도록 하겠습니다.",
    },
}

return {
    on_click = function(me, npc)
        fourth_class.advance(me, npc, CFG)
    end
}
