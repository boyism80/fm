package wz

import (
	"time"

	"github.com/boyism80/fm/types"
)

// fieldBossSpawnID keys a field boss among a map's life spawns, past any WZ life index.
const fieldBossSpawnID = ^uint32(0)

// fieldBoss is a boss the original server spawned from its own table rather than from WZ life: one at a time, at
// one of three points, mobTime after the last one died, with a notice to the map. Two of its entries are left out
// because their mob is already in the map's WZ life.
type fieldBoss struct {
	mobID     uint32
	mobTime   time.Duration
	positions [3]types.Vector2[int16]
	notice    string
}

var fieldBosses = map[uint32]fieldBoss{
	104000400: {2220000, 500 * time.Second, [3]types.Vector2[int16]{{X: 439, Y: 185}, {X: 301, Y: -85}, {X: 107, Y: -355}}, "서늘한 기운이 감돌며 마노가 나타났습니다."},
	229010100: {9880150, 2000 * time.Second, [3]types.Vector2[int16]{{X: 439, Y: 185}, {X: 301, Y: -85}, {X: 107, Y: -355}}, "장난끼 가득한 웃음소리와 함깨 루시드가 찾아왔습니다."},
	123456795: {3000008, 7000 * time.Second, [3]types.Vector2[int16]{{X: -223, Y: 323}, {X: -223, Y: 323}, {X: -223, Y: 323}}, "스산한 기운과 함깨 벨룸이 나타났다."},
	271040210: {8850011, 58000 * time.Second, [3]types.Vector2[int16]{{X: -400, Y: 113}, {X: -400, Y: 113}, {X: -400, Y: 113}}, "어두운 구석 한켠에서 시그너스 여제가 나타났다."},
	802000611: {9400271, 20000 * time.Second, [3]types.Vector2[int16]{{X: 383, Y: 20}, {X: 383, Y: 20}, {X: 383, Y: 20}}, "천상의 지배 전함, 니벨룽겐이 나타났다."},
	802000711: {9400295, 2000 * time.Second, [3]types.Vector2[int16]{{X: -830, Y: 43}, {X: -830, Y: 43}, {X: -830, Y: 43}}, "시간의 지배자 듀나스가 나타났다."},
	802000821: {9400376, 2000 * time.Second, [3]types.Vector2[int16]{{X: 93, Y: 335}, {X: 93, Y: 335}, {X: 93, Y: 335}}, "타락한 천사가 나를 주시하고있다."},
	802000824: {9400289, 2000 * time.Second, [3]types.Vector2[int16]{{X: 93, Y: 335}, {X: 93, Y: 335}, {X: 93, Y: 335}}, "아름다운 파괴의 천사 아우프 헤븐이 강림했다."},
	802000411: {9400270, 2000 * time.Second, [3]types.Vector2[int16]{{X: 2434, Y: 114}, {X: 2434, Y: 114}, {X: 2434, Y: 114}}, "의문의 남자가 나타났다."},
	802000211: {9400263, 5000 * time.Second, [3]types.Vector2[int16]{{X: 456, Y: 545}, {X: 456, Y: 545}, {X: 456, Y: 545}}, "지상 제어 병기, 베르가모트가 나타났다."},
	910811280: {9803059, 200 * time.Second, [3]types.Vector2[int16]{{X: -291, Y: 146}, {X: -291, Y: 146}, {X: -291, Y: 146}}, "음흉한 미소를 띄며 미로쿠가 나타났다."},
	910811320: {9803060, 200 * time.Second, [3]types.Vector2[int16]{{X: 187, Y: 62}, {X: 187, Y: 62}, {X: 187, Y: 62}}, "미카가미가 나타났다."},
	910811340: {9803061, 200 * time.Second, [3]types.Vector2[int16]{{X: 234, Y: 187}, {X: 234, Y: 187}, {X: 234, Y: 187}}, "렌게가 당신의 길을 가로막는다."},
	910811360: {9803062, 200 * time.Second, [3]types.Vector2[int16]{{X: 24, Y: 187}, {X: 24, Y: 187}, {X: 24, Y: 187}}, "오우라 쇼군이 나타났다."},
	910811380: {9803075, 200 * time.Second, [3]types.Vector2[int16]{{X: 44, Y: -205}, {X: 44, Y: -205}, {X: 44, Y: -205}}, "모든 일의 원흉, 노히메가 나타났다."},
	40000:     {8930000, 100 * time.Second, [3]types.Vector2[int16]{{X: 370, Y: 275}, {X: 370, Y: 275}, {X: 370, Y: 275}}, "위험! 메이플 아일랜드에 벨룸이 나타났다!"},
	800000000: {9421519, 100 * time.Second, [3]types.Vector2[int16]{{X: 171, Y: 35}, {X: 171, Y: 35}, {X: 171, Y: 35}}, "동물을 해치는 나쁜 악인이 나타났다."},
	800020000: {9460048, 100 * time.Second, [3]types.Vector2[int16]{{X: 798, Y: -1105}, {X: 798, Y: -1105}, {X: 798, Y: -1105}}, "요괴들의 왕이 모습을 드러냈다."},
	800030000: {9421547, 100 * time.Second, [3]types.Vector2[int16]{{X: 1247, Y: -1015}, {X: 1247, Y: -1015}, {X: 1247, Y: -1015}}, "어둠신의 애완동물, 멍왕이 나타났다."},
	800020120: {9421562, 100 * time.Second, [3]types.Vector2[int16]{{X: -47, Y: 155}, {X: -47, Y: 155}, {X: -47, Y: 155}}, "귀가 찢어질 듯 한 벼락 소리와 함깨, 텐구의 추종자, 뇌왕이 나타났다."},
	800020130: {9400014, 8000 * time.Second, [3]types.Vector2[int16]{{X: 439, Y: 185}, {X: 301, Y: -85}, {X: 107, Y: -355}}, "텐구가 선풍을 머금으며, 나타났습니다."},
	101073201: {3501008, 2000 * time.Second, [3]types.Vector2[int16]{{X: -89, Y: 245}, {X: -89, Y: 245}, {X: -89, Y: 245}}, "몰킹이 나타났다."},
	450001200: {8641010, 2000 * time.Second, [3]types.Vector2[int16]{{X: 900, Y: 177}, {X: 899, Y: 177}, {X: 898, Y: 177}}, "아르마가 나타났다."},
	450002250: {8642016, 2000 * time.Second, [3]types.Vector2[int16]{{X: 963, Y: 138}, {X: 1044, Y: 138}, {X: 741, Y: 138}}, "츄릅 나무가 나타났다."},
	101030404: {3220000, 500 * time.Second, [3]types.Vector2[int16]{{X: 867, Y: 1282}, {X: 810, Y: 1570}, {X: 838, Y: 2197}}, "바위산을 울리는 발걸음 소리와 함께 스텀피가 나타났습니다."},
	110040000: {5220001, 1000 * time.Second, [3]types.Vector2[int16]{{X: -355, Y: 179}, {X: -1283, Y: -113}, {X: -571, Y: -593}}, "모래사장에 수상한 소라껍질이 나타났습니다."},
	250010304: {7220000, 2100 * time.Second, [3]types.Vector2[int16]{{X: -210, Y: 33}, {X: -234, Y: 393}, {X: -654, Y: 33}}, "나즈막한 휘파람 소리와 함께 태륜이 나타났습니다."},
	200010300: {8220000, 1200 * time.Second, [3]types.Vector2[int16]{{X: 665, Y: 83}, {X: 672, Y: -217}, {X: -123, Y: -217}}, "검은 돌개바람과 함께 엘리쟈가 나타났습니다."},
	250010503: {7220002, 1800 * time.Second, [3]types.Vector2[int16]{{X: -303, Y: 543}, {X: 227, Y: 543}, {X: 719, Y: 543}}, "주변을 흐르는 요기가 강해졌습니다. 기분나쁜 고양이 울음소리가 들립니다."},
	222010310: {7220001, 1500 * time.Second, [3]types.Vector2[int16]{{X: -169, Y: -147}, {X: -517, Y: 93}, {X: 247, Y: 93}}, "달빛이 흐려지면서 긴 여우 울음소리와 함께 구미호의 기운이 느껴집니다."},
	107000300: {6220000, 1800 * time.Second, [3]types.Vector2[int16]{{X: 710, Y: 118}, {X: 95, Y: 119}, {X: -535, Y: 120}}, "늪 속에서 거대한 악어가 올라왔습니다."},
	100040105: {5220002, 1800 * time.Second, [3]types.Vector2[int16]{{X: 1000, Y: 278}, {X: 557, Y: 278}, {X: 95, Y: 278}}, "푸른 안개가 짙어지면서 파우스트가 나타났습니다."},
	100040106: {5220002, 1800 * time.Second, [3]types.Vector2[int16]{{X: 1000, Y: 278}, {X: 557, Y: 278}, {X: 95, Y: 278}}, "푸른 안개가 짙어지면서 파우스트가 나타났습니다."},
	220050100: {5220003, 1500 * time.Second, [3]types.Vector2[int16]{{X: -467, Y: 1032}, {X: 532, Y: 1032}, {X: -47, Y: 1032}}, "짹각짹각! 불규칙한 시계 소리와 함께 타이머가 나타났습니다."},
	221040301: {6220001, 2400 * time.Second, [3]types.Vector2[int16]{{X: -4134, Y: 416}, {X: -4283, Y: 776}, {X: -3292, Y: 776}}, "묵직한 기계음을 울리며 제노가 나타났습니다."},
	240040401: {8220003, 7200 * time.Second, [3]types.Vector2[int16]{{X: -15, Y: 2481}, {X: 127, Y: 1634}, {X: 159, Y: 1142}}, "차가운 바람과 함께 협곡 너머에서 레비아탄이 모습을 드러냈습니다."},
	260010201: {3220001, 1600 * time.Second, [3]types.Vector2[int16]{{X: -215, Y: 275}, {X: 298, Y: 275}, {X: 592, Y: 275}}, "모래먼지 속에서 데우가 천천히 모습을 드러냈습니다."},
	261030000: {8220002, 2700 * time.Second, [3]types.Vector2[int16]{{X: -1094, Y: -405}, {X: -772, Y: -116}, {X: -108, Y: 181}}, "지하의 어둠 속에서 안광을 빛내며 키메라가 모습을 드러냈습니다."},
	230020100: {4220000, 1200 * time.Second, [3]types.Vector2[int16]{{X: -291, Y: -20}, {X: -272, Y: -500}, {X: -462, Y: 640}}, "해초 수풀 사이로 이상한 조개가 나타났습니다."},
	209000000: {9500317, 300 * time.Second, [3]types.Vector2[int16]{{X: -115, Y: 154}, {X: -115, Y: 154}, {X: -115, Y: 154}}, "스노우맨이 나타났습니다!"},
}

func (m *Map) addFieldBoss() {
	boss, ok := fieldBosses[m.ID]
	if !ok {
		return
	}
	m.MobSpawns[fieldBossSpawnID] = MobSpawn{
		BaseSpawn: &BaseSpawn{ID: boss.mobID, Position: boss.positions[0], MobTime: boss.mobTime, Team: -1},
		Positions: boss.positions[:],
		Notice:    boss.notice,
	}
}
