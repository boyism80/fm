-- NPC name (String.wz/Npc.img.xml): 레고르

local fourth_class = require("script/lib/fourth_class")

local INTRO = "자유로움을 누리고 다니는 궁수는 항상 자유롭게 활동하는 직업입니다."

local ADVANCED = {
    [Class.Bowmaster] = "자네는 얼마 전에 #b보우마스터#k가 되었던 #r#h ##k이군, 수행은 잘 되는가? 자유로운 존재가 되기 위하여 더 정진하길 바라네..",
    [Class.Crossbowmaster] = "자네는 얼마 전에 #b신궁#k이 되었던 #r#h ##k이군, 수행은 잘 되는가? 자유로운 존재가 되기 위하여 더 정진하길 바라네..",
}

local ADVANCE = {
    [Class.Ranger] = {
        class = Class.Bowmaster,
        option = "#r#h # : #b보우마스터로 전직하고 싶습니다.",
        done =
        "이제 자네는 최고의 궁수인 #b보우마스터#k이 되었다네. 축하 선물로 자네에게 약간의 SP와 AP와 주었으니 확인해 보게나. 이로써 자네는 최고의 강한 궁수가 되었네. 하지만 앞으로 더 진보하기 위해서는 끊임없는 수련이 필요하다네...",
    },
    [Class.Sniper] = {
        class = Class.Crossbowmaster,
        option = "#r#h # : #b신궁으로 전직하고 싶습니다.",
        done =
        "이제 자네는 최고의 궁수인 #b신궁#k이 되었다네. 축하 선물로 자네에게 약간의 SP와 AP와 주었으니 확인해 보게나. 이로써 자네는 최고의 강한 궁수가 되었네. 하지만 앞으로 더 진보하기 위해서는 끊임없는 수련이 필요하다네...",
    },
}

local CFG = {
    quest = 6924,
    advanced = ADVANCED,
    advance = ADVANCE,
    text = {
        intro = INTRO,
        choose = "어느 클레스로 전직하기를 원하는가?\r\n그대의 영웅으로써의 자질은 증명되었다.\r\n",
        wait = "그대는 그대의 영웅으로서의 자질을 증명했습니다. \r\n이제 남은 것은 궁수 궁극의 길로 향하는 것 뿐.\r\n전직을 할 준비가 되면 다시 말을 걸어주세요.",
        sp_left = "흠... 당신은 너무 많은 #bSP#k를 갖고 있는 것 같군요. 최소한 레벨 120 이전에 얻은 스킬포인트를 모두 사용하셔야 4차 전직을 하실 수 있습니다.",
        no_space = "흠.. 소비 인벤토리 공간이 부족한 건 아닌가요? 다시 한번 확인해 보세요.",
        farewell = "난 자네가 더욱 더 수련하여 지상 최고의 궁수가 되는 것을 뒤에서 지켜봐주도록 하겠네.",
    },
}

return {
    on_click = function(me, npc)
        fourth_class.advance(me, npc, CFG)
    end
}
