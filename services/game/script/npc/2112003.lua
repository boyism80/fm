-- NPC name (String.wz/Npc.img.xml): 줄리엣

return require("script/lib/romeo_juliet").npc({
    group = "juliet_party_quest",
    entry_map = 261000021,
    base = 926110000,
    other = 926100000,
    story = "줄리엣의 이야기를 듣는다",
    letter = 4001130,
    letter_before_leader = true,
    letter_notice = "줄리엣은 로미오가 쓴 편지를 보고 생각에 잠겼다.",
    letter_reply = "이.. 이것은 로미오가 쓴 편지..?",
    stage4_exp = 10000,
    door = "jnr3_out3",
    beloved = "로미오를",
    marble = 4001160,
    lab_hint = "이 연구실에서 가끔 수상한 소리가 들린다는 소문이 있었어요. 분명 이 근처 어딘가에 수상한 소리의 근원이 있을 거에요.",
    mobs_hint = "모든 몬스터를 물리쳐 주세요!",
    beaker_hint = "이곳의 비커들은 깨져서 내용물이 새고 있답니다. 비커에 내용물이 모두 새기 전에 빠르게 수상한 액체를 채워 넣어주세요.",
})
