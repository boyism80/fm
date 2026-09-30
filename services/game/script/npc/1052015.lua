-- NPC name (String.wz/Npc.img.xml): 빌리

local folders = {
	{ "Bgm00", { { "FloralLife", "Bgm00/FloralLife" }, { "GoPicnic", "Bgm00/GoPicnic" }, { "Nightmare", "Bgm00/Nightmare" }, { "RestNPeace", "Bgm00/RestNPeace" }, { "SleepyWood", "Bgm00/SleepyWood" } } },
	{ "Bgm01", { { "AncientMove", "Bgm01/AncientMove" }, { "BadGuys", "Bgm01/BadGuys" }, { "CavaBien", "Bgm01/CavaBien" }, { "HighlandStar", "Bgm01/HighlandStar" }, { "MoonlightShadow", "Bgm01/MoonlightShadow" }, { "WhereTheBarlogFrom", "Bgm01/WhereTheBarlogFrom" } } },
	{ "Bgm02", { { "AboveTheTreetops", "Bgm02/AboveTheTreetops" }, { "EvilEyes", "Bgm02/EvilEyes" }, { "JungleBook", "Bgm02/JungleBook" }, { "MissingYou", "Bgm02/MissingYou" }, { "WhenTheMorningComes", "Bgm02/WhenTheMorningComes" } } },
	{ "Bgm03", { { "Beachway", "Bgm03/Beachway" }, { "BlueSky", "Bgm03/BlueSky" }, { "Elfwood", "Bgm03/Elfwood" }, { "SnowyVillage", "Bgm03/SnowyVillage" }, { "Subway", "Bgm03/Subway" } } },
	{ "Bgm04", { { "ArabPirate", "Bgm04/ArabPirate" }, { "PlayWithMe", "Bgm04/PlayWithMe" }, { "Shinin'Harbor", "Bgm04/Shinin'Harbor" }, { "UponTheSky", "Bgm04/UponTheSky" }, { "WarmRegard", "Bgm04/WarmRegard" }, { "WhiteChristmas", "Bgm04/WhiteChristmas" } } },
	{ "Bgm05", { { "AbandonedMine", "Bgm05/AbandonedMine" }, { "DownToTheCave", "Bgm05/DownToTheCave" }, { "HellGate", "Bgm05/HellGate" }, { "MineQuest", "Bgm05/MineQuest" }, { "WolfWood", "Bgm05/WolfWood" } } },
	{ "Bgm06", { { "ComeWithMe", "Bgm06/ComeWithMe" }, { "FantasticThinking", "Bgm06/FantasticThinking" }, { "FinalFight", "Bgm06/FinalFight" }, { "FlyingInABlueDream", "Bgm06/FlyingInABlueDream" }, { "WelcomeToTheHell", "Bgm06/WelcomeToTheHell" } } },
	{ "Bgm07", { { "Fantasia", "Bgm07/Fantasia" }, { "FunnyTimeMaker", "Bgm07/FunnyTimeMaker" }, { "HighEnough", "Bgm07/HighEnough" }, { "WaltzForWork", "Bgm07/WaltzForWork" }, { "WhereverYouAre", "Bgm07/WhereverYouAre" } } },
	{ "Bgm08", { { "FindingForest", "Bgm08/FindingForest" }, { "ForTheGlory", "Bgm08/ForTheGlory" }, { "LetsHuntAliens", "Bgm08/LetsHuntAliens" }, { "LetsMarch", "Bgm08/LetsMarch" }, { "PlotOfPixie", "Bgm08/PlotOfPixie" } } },
	{ "Bgm09", { { "DarkShadow", "Bgm09/DarkShadow" }, { "FairyTale", "Bgm09/FairyTale" }, { "FairyTalediffvers", "Bgm09/FairyTalediffvers" }, { "TheyMenacingYou", "Bgm09/TheyMenacingYou" }, { "TimeAttack", "Bgm09/TimeAttack" } } },
	{ "Bgm10", { { "BizarreTales", "Bgm10/BizarreTales" }, { "Eregos", "Bgm10/Eregos" }, { "TheWayGrotesque", "Bgm10/TheWayGrotesque" }, { "Timeless", "Bgm10/Timeless" }, { "TimelessB", "Bgm10/TimelessB" } } },
	{ "Bgm11", { { "Aquarium", "Bgm11/Aquarium" }, { "BlueWorld", "Bgm11/BlueWorld" }, { "DarkMountain", "Bgm11/DarkMountain" }, { "DownTown", "Bgm11/DownTown" }, { "ShiningSea", "Bgm11/ShiningSea" } } },
	{ "Bgm12", { { "AcientRemain", "Bgm12/AcientRemain" }, { "AquaCave", "Bgm12/AquaCave" }, { "DeepSee", "Bgm12/DeepSee" }, { "Dispute", "Bgm12/Dispute" }, { "RuinCastle", "Bgm12/RuinCastle" }, { "WaterWay", "Bgm12/WaterWay" } } },
	{ "Bgm13", { { "AcientForest", "Bgm13/AcientForest" }, { "CokeTown", "Bgm13/CokeTown" }, { "FightSand", "Bgm13/FightSand" }, { "Leafre", "Bgm13/Leafre" }, { "Minar'sDream", "Bgm13/Minar'sDream" }, { "TowerOfGoddess", "Bgm13/TowerOfGoddess" } } },
	{ "Bgm14", { { "Ariant", "Bgm14/Ariant" }, { "CaveOfHontale", "Bgm14/CaveOfHontale" }, { "DragonLoad", "Bgm14/DragonLoad" }, { "DragonNest", "Bgm14/DragonNest" }, { "HonTale", "Bgm14/HonTale" }, { "HotDesert", "Bgm14/HotDesert" } } },
	{ "Bgm15", { { "ElinForest", "Bgm15/ElinForest" }, { "inNautilus", "Bgm15/inNautilus" }, { "MureungForest", "Bgm15/MureungForest" }, { "MureungHill", "Bgm15/MureungHill" }, { "Nautilus", "Bgm15/Nautilus" }, { "Pirate", "Bgm15/Pirate" }, { "PoisonForest", "Bgm15/PoisonForest" }, { "SunsetDesert", "Bgm15/SunsetDesert" }, { "WhiteHerb", "Bgm15/WhiteHerb" } } },
	{ "Bgm16", { { "Duskofgod", "Bgm16/Duskofgod" }, { "FightingPinkBeen", "Bgm16/FightingPinkBeen" }, { "Forgetfulness", "Bgm16/Forgetfulness" }, { "Remembrance", "Bgm16/Remembrance" }, { "Repentance", "Bgm16/Repentance" }, { "TimeTemple", "Bgm16/TimeTemple" } } },
	{ "Bgm17", { { "MureungSchool1", "Bgm17/MureungSchool1" }, { "MureungSchool2", "Bgm17/MureungSchool2" }, { "MureungSchool3", "Bgm17/MureungSchool3" }, { "MureungSchool4", "Bgm17/MureungSchool4" } } },
	{ "Bgm18", { { "BlackWing", "Bgm18/BlackWing" }, { "DrillHall", "Bgm18/DrillHall" }, { "QueensGarden", "Bgm18/QueensGarden" }, { "RaindropFlower", "Bgm18/RaindropFlower" } } },
	{ "BgmCN", { { "GoShanghai", "BgmCN/GoShanghai" }, { "ShanghaiField", "BgmCN/ShanghaiField" } } },
	{ "BgmEvent", { { "FunnyRabbit", "BgmEvent/FunnyRabbit" }, { "FunnyRabbitFaster", "BgmEvent/FunnyRabbitFaster" }, { "wedding", "BgmEvent/wedding" }, { "weddingDance", "BgmEvent/weddingDance" } } },
	{ "BgmGL", { { "amoria", "BgmGL/amoria" }, { "cathedral", "BgmGL/cathedral" }, { "chapel", "BgmGL/chapel" }, { "HauntedHouse", "BgmGL/HauntedHouse" } } },
	{ "BgmJp", { { "Bathroom", "BgmJp/Bathroom" }, { "BattleField", "BgmJp/BattleField" }, { "BizarreForest", "BgmJp/BizarreForest" }, { "Feeling", "BgmJp/Feeling" }, { "FirstStepMaster", "BgmJp/FirstStepMaster" }, { "Hana", "BgmJp/Hana" }, { "Yume", "BgmJp/Yume" } } },
	{ "BgmTH", { { "ThaiField", "BgmTH/ThaiField" }, { "ThaiTown", "BgmTH/ThaiTown" } } },
	{ "BgmTW", { { "NightField", "BgmTW/NightField" }, { "NightMarket", "BgmTW/NightMarket" }, { "YoTaipei", "BgmTW/YoTaipei" } } },
	{ "한글", { { "1분1초", "한글/1분1초" }, { "Blue", "한글/Blue" }, { "MyLove", "한글/MyLove" }, { "RPG1", "한글/RPG1" }, { "RPG2", "한글/RPG2" }, { "곰세마리", "한글/곰세마리" }, { "구구단송", "한글/구구단송" }, { "그체", "한글/그체" }, { "나는나비", "한글/나는나비" }, { "메칸더", "한글/메칸더" }, { "바람의너를", "한글/바람의너를" }, { "벚꽃엔딩", "한글/벚꽃엔딩" }, { "별똥별", "한글/별똥별" }, { "빠빠빠", "한글/빠빠빠" }, { "새나라새주인", "한글/새나라새주인" }, { "섬집아기", "한글/섬집아기" }, { "세일러문", "한글/세일러문" }, { "슈퍼스타", "한글/슈퍼스타" }, { "아빠의얼굴", "한글/아빠의얼굴" }, { "애국가1(소리 매우 작음)", "한글/애국가1" }, { "애국가2", "한글/애국가2" }, { "어머니은혜", "한글/어머니은혜" }, { "여래아", "한글/여래아" }, { "작은별", "한글/작은별" }, { "조각나비", "한글/조각나비" }, { "질풍가도", "한글/질풍가도" }, { "천사소녀네티", "한글/천사소녀네티" }, { "체리(소리 큼)", "한글/체리" }, { "포켓몬", "한글/포켓몬" }, { "혼자가아닌나", "한글/혼자가아닌나" } } },
	{ "ef길티모노가타리", { { "Departures", "ef길티모노가타리/Departures" }, { "ef1", "ef길티모노가타리/ef1" }, { "ef2", "ef길티모노가타리/ef2" }, { "ReleaseMySoul", "ef길티모노가타리/ReleaseMySoul" }, { "고백", "ef길티모노가타리/고백" }, { "백금디스코", "ef길티모노가타리/백금디스코" }, { "세노", "ef길티모노가타리/세노" }, { "영원한유죄", "ef길티모노가타리/영원한유죄" }, { "출발", "ef길티모노가타리/출발" } } },
	{ "kiss이로하팝", { { "MarryMe", "kiss이로하팝/MarryMe" }, { "TheGreatEscape", "kiss이로하팝/TheGreatEscape" }, { "이로하1", "kiss이로하팝/이로하1" }, { "이로하2", "kiss이로하팝/이로하2" }, { "키스시스1", "kiss이로하팝/키스시스1" }, { "키스시스2", "kiss이로하팝/키스시스2" } } },
	{ "Maple추가브금", { { "Anniv1", "Maple/Anniv1" }, { "사자왕의 성", "Maple/BlizzardCastle" }, { "시그너스 정원", "Maple/CygnusGarden" }, { "팬텀", "Maple/DancingWithTheMoon" }, { "커닝스퀘어", "Maple/KerningSquare" }, { "커닝스퀘어 필드", "Maple/KerningSquareField" }, { "KnightsStronghold", "Maple/KnightsStronghold" }, { "LowGradeOre", "Maple/LowGradeOre" }, { "MapleLeaf", "Maple/MapleLeaf" }, { "PowerStation", "Maple/PowerStation" }, { "Smile", "Maple/Smile" }, { "Title_Japan", "Maple/Title_Japan" }, { "Title2010Winter", "Maple/Title2010Winter" }, { "보컬오프 Smile", "Maple/VO_Smile" }, { "에델슈타인", "Maple/에델슈타인" }, { "에우렐", "Maple/에우렐" }, { "엔젤릭버스터", "Maple/엔젤릭버스터" }, { "크리스탈가든", "Maple/크리스탈가든" } } },
	{ "supercell", { { "Aozora", "supercell/Aozora" }, { "LoveMeGimmie", "supercell/LoveMeGimmie" }, { "Palette", "supercell/Palette" }, { "TheBravery", "supercell/TheBravery" }, { "길티크라운", "supercell/길티크라운" }, { "너의모르는이야기", "supercell/너의모르는이야기" }, { "락앤롤", "supercell/락앤롤" }, { "러브앤롤", "supercell/러브앤롤" }, { "리루모아", "supercell/리루모아" }, { "마마파파", "supercell/마마파파" }, { "모시모", "supercell/모시모" }, { "별이반짝이는", "supercell/별이반짝이는" }, { "복수", "supercell/복수" }, { "사요나라메모리즈", "supercell/사요나라메모리즈" }, { "앨범op", "supercell/시작" }, { "오시에테아게루", "supercell/오시에테아게루" }, { "앨범ed", "supercell/오와리" }, { "요루(sheep)", "supercell/요루" }, { "월드와이드", "supercell/월드와이드" }, { "은색비행선", "supercell/은색비행선" }, { "첫사랑이 끝날 때", "supercell/첫사랑이끝날때" }, { "퍼펙트데이", "supercell/퍼펙트데이" }, { "필소굿", "supercell/필소굿" }, { "하나비", "supercell/하나비" }, { "히어로", "supercell/히어로" } } },
	{ "yui", { { "Gloria", "yui/Gloria" }, { "GoodByeDays", "yui/GoodByeDays" }, { "Laugh_away", "yui/Laugh_away" }, { "Life", "yui/Life" } } },
	{ "달동네", { { "Arcadia", "달동네/Arcadia" }, { "Disillusion2", "달동네/Disillusion2" }, { "구름조각", "달동네/구름조각" }, { "제로1ed", "달동네/제로1ed" }, { "제로1op", "달동네/제로1op" }, { "제로2ed", "달동네/제로2ed" }, { "제로2op", "달동네/제로2op" }, { "제로 앨범 삽입곡", "달동네/제로투투투" }, { "카니발ed", "달동네/카니발ed" }, { "카니발op", "달동네/카니발op" }, { "프리즘이리야op", "달동네/프리즘이리야op" } } },
	{ "달빛천사", { { "나의마음을담아", "달빛천사/나의마음을담아" }, { "내자신의", "달빛천사/내자신의" }, { "미소", "달빛천사/미소" }, { "사랑의기록표", "달빛천사/사랑의기록표" }, { "새로운미래", "달빛천사/새로운미래" }, { "영원한눈", "달빛천사/영원한눈" } } },
	{ "동방", { { "마법소녀들", "동방/마법소녀들" }, { "신들이사랑한", "동방/신들이사랑한" }, { "유령악단", "동방/유령악단" } } },
	{ "디지몬", { { "어드벤쳐ed", "디지몬/어드벤쳐ed" }, { "어드벤쳐op", "디지몬/어드벤쳐op" }, { "어드벤쳐ost", "디지몬/어드벤쳐ost" }, { "어드벤쳐볼레로", "디지몬/어드벤쳐볼레로" }, { "어드벤쳐서브타이틀", "디지몬/어드벤쳐서브타이틀" }, { "어드벤쳐시작", "디지몬/어드벤쳐시작" }, { "어드벤쳐악당", "디지몬/어드벤쳐악당" }, { "어드벤쳐위기", "디지몬/어드벤쳐위기" }, { "어드벤쳐전투", "디지몬/어드벤쳐전투" }, { "어드벤쳐진화", "디지몬/어드벤쳐진화" }, { "어드벤쳐파워업", "디지몬/어드벤쳐파워업" }, { "파워ed", "디지몬/파워ed" }, { "파워op", "디지몬/파워op" }, { "파워진화1", "디지몬/파워진화1" }, { "파워진화2", "디지몬/파워진화2" }, { "파워진화2한글", "디지몬/파워진화2한글" } } },
	{ "럭키하루히학생회일상", { { "세라복2(기본 빠른 템포)", "럭키하루히학생회일상/세라복2" }, { "일상ed", "럭키하루히학생회일상/일상ed" }, { "일상op1", "럭키하루히학생회일상/일상op1" }, { "일상op2", "럭키하루히학생회일상/일상op2" }, { "하루히ost1", "럭키하루히학생회일상/하루히ost1" }, { "하루히ost2", "럭키하루히학생회일상/하루히ost2" }, { "학생회ed", "럭키하루히학생회일상/학생회ed" } } },
	{ "몬무스", { { "izumi", "몬무스/izumi" }, { "댄스", "몬무스/댄스" }, { "만담", "몬무스/만담" }, { "반성", "몬무스/반성" }, { "슬픔1", "몬무스/슬픔1" }, { "슬픔2", "몬무스/슬픔2" }, { "시르후", "몬무스/시르후" }, { "아쿠아", "몬무스/아쿠아" }, { "야영", "몬무스/야영" }, { "에덴", "몬무스/에덴" }, { "연구실", "몬무스/연구실" }, { "이리아스", "몬무스/이리아스" }, { "최종보스", "몬무스/최종보스" }, { "크롬", "몬무스/크롬" }, { "타마모", "몬무스/타마모" }, { "플랜섹트", "몬무스/플랜섹트" } } },
	{ "보컬오프", { { "Aozora", "보컬오프/Aozora" }, { "AsitanoSora", "보컬오프/AsitanoSora" }, { "Everything", "보컬오프/Everything" }, { "GentleJena", "보컬오프/GentleJena" }, { "GrowSlowly", "보컬오프/GrowSlowly" }, { "HappySong", "보컬오프/HappySong" }, { "LoveMeGimmie", "보컬오프/LoveMeGimmie" }, { "MyDearest", "보컬오프/MyDearest" }, { "Oshieteageru", "보컬오프/Oshieteageru" }, { "Palette", "보컬오프/Palette" }, { "TheBravery", "보컬오프/TheBravery" }, { "경단대가족", "보컬오프/경단대가족" }, { "꽃의댄스", "보컬오프/꽃의댄스" }, { "나친적op", "보컬오프/나친적op" }, { "네가모르는이야기", "보컬오프/네가모르는이야기" }, { "놀라운그레이스", "보컬오프/놀라운그레이스" }, { "달묘전설", "보컬오프/달묘전설" }, { "마비노기", "보컬오프/마비노기" }, { "메모리아", "보컬오프/메모리아" }, { "백금디스코", "보컬오프/백금디스코" }, { "별빛쪼개기", "보컬오프/별빛쪼개기" }, { "세라샵", "보컬오프/세라샵" }, { "센과치히로", "보컬오프/센과치히로" }, { "소원이이루어지는장소", "보컬오프/소원이이루어지는장소" }, { "시대를초월한마음", "보컬오프/시대를초월한마음" }, { "아즈망가대왕", "보컬오프/아즈망가대왕" }, { "언덕아래의이별", "보컬오프/언덕아래의이별" }, { "익시온사가게임", "보컬오프/익시온사가게임" }, { "캐롤샵", "보컬오프/캐롤샵" }, { "크로아티아랩소디", "보컬오프/크로아티아랩소디" }, { "클라나드일렉톤", "보컬오프/클라나드일렉톤" }, { "테일즈위버ost1", "보컬오프/테일즈위버ost1" }, { "테일즈위버ost2", "보컬오프/테일즈위버ost2" }, { "테일즈위버ost3", "보컬오프/테일즈위버ost3" }, { "테일즈위버ost4", "보컬오프/테일즈위버ost4" }, { "포켓몬도로", "보컬오프/포켓몬도로" }, { "하늘에빛나다", "보컬오프/하늘에빛나다" }, { "하울", "보컬오프/하울" }, { "해적", "보컬오프/해적" }, { "히로시의회상", "보컬오프/히로시의회상" } } },
	{ "신만세", { { "사랑의예감", "신만세/사랑의예감" }, { "신만세1op", "신만세/신만세1op" }, { "신만세1op2", "신만세/신만세1op2" }, { "신만세2op", "신만세/신만세2op" }, { "신만세3op", "신만세/신만세3op" }, { "집적회로", "신만세/집적회로" }, { "집적회로2", "신만세/집적회로2" } } },
	{ "약", { { "Electric_Six_Gay_Bar", "약/Electric_Six_Gay_Bar" }, { "I_just_had_sex", "약/I_just_had_sex" }, { "You_spin_me_round", "약/You_spin_me_round" }, { "가그린기자", "약/가그린기자" }, { "고자라니", "약/고자라니" }, { "기자리믹스1", "약/기자리믹스1" }, { "기자리믹스2", "약/기자리믹스2" }, { "남극탐험", "약/남극탐험" }, { "뚤훍뚤", "약/뚤훍뚤" }, { "라면한박스", "약/라면한박스" }, { "마제윤주제가", "약/마제윤주제가" }, { "문학소녀", "약/문학소녀" }, { "발랄라이카", "약/발랄라이카" }, { "비둘기야", "약/비둘기야" }, { "빠삐놈1", "약/빠삐놈1" }, { "빠삐놈2", "약/빠삐놈2" }, { "빠삐코기자", "약/빠삐코기자" }, { "숨겨왔던", "약/숨겨왔던" }, { "스머프 양념통닭", "약/스머프" }, { "스폰지밥op", "약/스폰지밥op" }, { "스폰지밥크리스마스", "약/스폰지밥크리스마스" }, { "시유천", "약/시유천" }, { "심영", "약/심영" }, { "앵그리기자", "약/앵그리기자" }, { "야라나이카", "약/야라나이카" }, { "에어맨", "약/에어맨" }, { "오리온좌", "약/오리온좌" }, { "우마우마", "약/우마우마" }, { "운지천", "약/운지천" }, { "이루마고자라니", "약/이루마고자라니" }, { "케이크하우스", "약/케이크하우스" }, { "파돌리기송", "약/파돌리기송" }, { "피카츄기자", "약/피카츄기자" } } },
	{ "어떤세상", { { "Dear My Friend", "어과초어마금/DearMyFriend" }, { "Grow Slowly", "어과초어마금/GrowSlowly" }, { "Late in autumn", "어과초어마금/Late_in_autumn" }, { "Memory of snow", "어과초어마금/Memory_of_snow" }, { "Only my lailgun", "어과초어마금/OnlyMyLailgun" }, { "Smile", "어과초어마금/Smile" }, { "Special one", "어과초어마금/Special_One" }, { "질주감", "어과초어마금/질주감" } } },
	{ "엔젤비트", { { "Alchemy", "엔젤비트/Alchemy" }, { "BraveSong", "엔젤비트/BraveSong" }, { "CrowSong", "엔젤비트/CrowSong" }, { "God_bless_you", "엔젤비트/God_bless_you" }, { "HighestSong", "엔젤비트/HighestSong" }, { "HotMeal", "엔젤비트/HotMeal" }, { "LastSong", "엔젤비트/LastSong" }, { "MySong", "엔젤비트/MySong" }, { "ShineDays", "엔젤비트/ShineDays" }, { "ThousandEnemies", "엔젤비트/ThousandEnemies" }, { "엔비op", "엔젤비트/엔비op" }, { "최고의보물ori", "엔젤비트/최고의보물ori" }, { "최고의보물yui", "엔젤비트/최고의보물yui" } } },
	{ "짜투리", { { "AliceMagic", "잡/AliceMagic" }, { "FirstKiss", "잡/FirstKiss" }, { "You", "잡/You" }, { "나친적ed", "잡/나친적ed" }, { "남고생op", "잡/남고생op" }, { "네기마op", "잡/네기마op" }, { "늑향ed", "잡/늑향ed" }, { "베토벤 바이러스", "잡/베토벤" }, { "브툼op", "잡/브툼op" }, { "블리치13op", "잡/블리치13op" }, { "빈보가미op", "잡/빈보가미op" }, { "시달소ost", "잡/시달소ost" }, { "프리징op", "잡/프리징op" } } },
	{ "전파케이온아이들", { { "Heart Goes Boom!!", "전파케이온아이들/Boom" }, { "Don't say 'lazy'", "전파케이온아이들/Lazy" }, { "Listen!", "전파케이온아이들/Listen" }, { "NO,Thank You!", "전파케이온아이들/Thank" }, { "카제히쿠노", "전파케이온아이들/전파녀ed" }, { "하나마루센세이션", "전파케이온아이들/코도모ed" } } },
	{ "클라나드", { { "경단대가족", "클라나드/경단대가족" }, { "기쁨의 섬", "클라나드/기쁨의섬" }, { "사쿠라 사쿠라 아이타이요", "클라나드/사쿠라" }, { "시간을 새기는 노래", "클라나드/시간을새기는노래" }, { "작은 손바닥", "클라나드/작은손바닥" } } },
	{ "토라도라 진격거 아노하나", { { "Vanilla_Salt", "토라도라진격거아노하나/바닐라소금" }, { "아노하나op", "토라도라진격거아노하나/아노하나op" }, { "아노하나ost(secret base)", "토라도라진격거아노하나/아노하나ost" }, { "Orange", "토라도라진격거아노하나/오렌지1" }, { "Orange2", "토라도라진격거아노하나/오렌지2" }, { "진격거ed", "토라도라진격거아노하나/진격거ed" }, { "진격거op", "토라도라진격거아노하나/진격거op" }, { "진격거ost", "토라도라진격거아노하나/진격거ost" }, { "HolyNight", "토라도라진격거아노하나/홀리나이트" } } },
	{ "하야테", { { "유카1", "하야테/유카1" }, { "유카2", "하야테/유카2" }, { "하야테1-2ed", "하야테/하야테1-2ed" }, { "하야테1-3ed", "하야테/하야테1-3ed" }, { "하야테3ed", "하야테/하야테3ed" }, { "하야테3op", "하야테/하야테3op" }, { "하야테 극장판 ed", "하야테/하야테극장판ed" }, { "Princess is me!", "하야테/하야테이즈미" }, { "축하", "하야테/축하" }, { "히나기쿠", nil, "히나 히나 히나 히나 히나 히나 히나!!", { { "Do my best!", "하야테/히나기쿠/Do_my_best" }, { "Honjitsu Mankai", "하야테/히나기쿠/Honjitsu" }, { "I miss you", "하야테/히나기쿠/I_miss_you" }, { "Power of flower", "하야테/히나기쿠/PowerOfFlower" }, { "Tensi", "하야테/히나기쿠/Tensi" } } } } },
}

local function play(me, npc, tracks, prompt)
	local names = {}
	for i, track in ipairs(tracks) do
		names[i] = track[1]
	end
	local sel = me:dialog_list(npc, prompt or "뭐 켤꺼야?", names)
	if sel == nil then
		return
	end
	local track = tracks[sel]
	if track == nil then
		return
	end
	if track[4] ~= nil then
		play(me, npc, track[4], track[3])
		return
	end
	local map = me:map()
	if map == nil then
		return
	end
	map:music(track[2])
end

return {
	on_click = function(me, npc)
		local map = me:map()
		if map ~= nil and map:template_id() == 193000000 then
			return
		end
		local names = {}
		for i, folder in ipairs(folders) do
			names[i] = folder[1]
		end
		local sel = me:dialog_list(npc, "안녕하신가! 힘세고 강한아침, 만일 내게 물어보면 나는 DJ빌리.", names)
		if sel == nil then
			return
		end
		local folder = folders[sel]
		if folder == nil then
			return
		end
		play(me, npc, folder[2])
	end
}
