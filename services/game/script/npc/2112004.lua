-- NPC name (String.wz/Npc.img.xml): 로미오

return require("script/lib/romeo_juliet").npc({
    group = "romeo_party_quest",
    entry_map = 261000011,
    base = 926100000,
    other = 926110000,
    story = "로미오의 이야기를 듣는다",
    letter = 4001131,
    letter_before_leader = false,
    letter_notice = "로미오는 줄리엣이 쓴 편지를 보고 생각에 잠겼다.",
    letter_reply = "이.. 이것은 줄리엣이 제게 쓴 편지..?",
    stage4_exp = 5000,
    door = "rnj3_out3",
    beloved = "줄리엣을",
    marble = 4001159,
    lab_hint = "이 연구실에서 가끔 수상한 소리가 들린다는 소문이 있었어요. 분명 이 근처 어딘가에 수상한 소리의 근원이 있을겁니다.",
    mobs_hint = "모든 몬스터를 물리쳐 주십시오.",
    beaker_hint = "이곳의 비커들은 깨져서 내용물이 새고 있습니다. 비커에 내용물이 모두 새기 전에 빠르게 수상한 액체를 채워 넣어주십시오.",
})
