package storage

import (
	"leaderboard/internal/model"
)

// Holds all Players and Competitions in memory
var (
	Players      = map[string]*model.Player{}
	Competitions = map[string]*model.Competition{}
)

// TODO: Define an interface
func LoadDummyPlayers() {
	AddPlayers(dummyPlayers[:])
}

func AddPlayers(players []NewPlayer) {
	for _, dummy := range players {
		player := model.NewPlayer(dummy.Id, dummy.Level, dummy.CountryCode)
		Players[player.Id()] = player
	}
}

type NewPlayer struct {
	Id          string
	CountryCode string
	Level       uint
}

var dummyPlayers = [...]NewPlayer{
	{"alice_smith", "US", 7}, {"bob_jones", "GB", 6}, {"carlos_mendez", "MX", 3},
	{"diana_lee", "KR", 9}, {"ethan_wong", "SG", 5}, {"fatima_khan", "PK", 8},
	{"george_nash", "AU", 4}, {"hana_tanaka", "JP", 6}, {"ian_clark", "CA", 7},
	{"julia_roberts", "US", 3}, {"kevin_singh", "IN", 5}, {"laura_gomez", "ES", 2},
	{"mohamed_fahmy", "EG", 8}, {"natalie_cho", "KR", 7}, {"oliver_klein", "DE", 6},
	{"paula_martins", "BR", 4}, {"quentin_diaz", "AR", 3}, {"rashid_ali", "AE", 7},
	{"sofia_perez", "CO", 5}, {"tommy_nilsen", "NO", 6}, {"ursula_meier", "CH", 4},
	{"viktor_ivanov", "RU", 9}, {"wanda_nowak", "PL", 2}, {"xinyi_zhang", "CN", 7},
	{"youssef_khaled", "MA", 5}, {"zara_kapoor", "IN", 6}, {"anton_karlsen", "SE", 3},
	{"bella_larsson", "SE", 4}, {"cedric_dubois", "FR", 5}, {"daria_smirnova", "RU", 7},
	{"emil_petrov", "BG", 4}, {"fiona_kennedy", "IE", 6}, {"gustavo_silva", "BR", 7},
	{"harper_clarke", "CA", 3}, {"isabel_nunez", "ES", 5}, {"jack_yamamoto", "JP", 6},
	{"karim_rahman", "BD", 2}, {"lucia_ferrari", "IT", 6}, {"matteo_gallo", "IT", 7},
	{"nadine_bakker", "NL", 8}, {"oscar_svensson", "SE", 5}, {"penelope_dunne", "NZ", 3},
	{"qiang_liu", "CN", 7}, {"raj_patel", "IN", 4}, {"samantha_brown", "US", 6},
	{"thomas_muller", "DE", 5}, {"ulrik_berg", "NO", 6}, {"valeria_ramos", "PE", 7},
	{"william_owen", "GB", 3}, {"xiaoli_chen", "CN", 5}, {"yuki_sato", "JP", 8},
	{"zeynep_aydin", "TR", 9}, {"aaron_clark", "US", 4}, {"bianca_santos", "BR", 6},
	{"cameron_wilson", "GB", 5}, {"daniela_rodriguez", "CL", 4}, {"emre_ozdemir", "TR", 3},
	{"farah_abdul", "SA", 7}, {"goran_jovanovic", "RS", 4}, {"hana_elmadi", "DZ", 6},
	{"ibrahim_sow", "SN", 8}, {"johan_nilsson", "SE", 5}, {"karla_ortega", "MX", 4},
	{"liam_brennan", "IE", 3}, {"maya_cohen", "IL", 6}, {"neha_sharma", "IN", 7},
	{"otto_meier", "DE", 5}, {"pia_koskinen", "FI", 6}, {"quentin_martel", "FR", 7},
	{"rohan_desai", "IN", 3}, {"sara_haddad", "MA", 8}, {"timo_ahlberg", "FI", 5},
	{"uliana_kuznetsova", "RU", 7}, {"victoria_lopes", "BR", 4}, {"wassim_hamdi", "TN", 6},
	{"ximena_vega", "PE", 2}, {"yasir_khan", "PK", 9}, {"zsolt_szabo", "HU", 5},
	{"adriana_dumitru", "RO", 6}, {"benjamin_hoffman", "DE", 4}, {"chloe_marchand", "FR", 3},
	{"diego_fuentes", "UY", 5}, {"eva_meyer", "CH", 6}, {"felix_morales", "EC", 7},
	{"grace_yip", "SG", 4}, {"hugo_leblanc", "CA", 5}, {"ines_fernandez", "ES", 6},
	{"james_ryan", "US", 8}, {"khalid_nasser", "JO", 7}, {"leila_saidi", "DZ", 4},
	{"marko_todorovic", "RS", 5}, {"nadia_samir", "EG", 6}, {"omar_hassan", "SD", 3},
	{"patricia_silva", "PT", 5}, {"ricardo_garcia", "MX", 6}, {"salma_abdel", "EG", 7},
	{"tarek_mansour", "LB", 4}, {"ulric_weber", "DE", 6}, {"valentina_mora", "CL", 5},
	{"wenjie_li", "CN", 7}, {"xander_blake", "US", 3}, {"yasmin_almawi", "SA", 6},
	{"zachary_reed", "US", 4},
}
