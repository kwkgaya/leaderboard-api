package storage

import (
	"encoding/json"
	"leaderboard/internal/model"
)

// Holds all Players and Competitions in memory
var (
	Players      = map[string]*model.Player{}
	Competitions = map[string]*model.Competition{}
)

// TODO: Define an interface
func LoadDummyPlayers() {
	var dummyPlayers []NewPlayer
	err := json.Unmarshal([]byte(dummyPlayersJson), &dummyPlayers)
	if err != nil {
		panic("Failed to load dummy players: " + err.Error())
	}
	AddPlayers(dummyPlayers)
}

func AddPlayers(players []NewPlayer) {
	for _, dummy := range players {
		player := model.NewPlayer(dummy.Id, dummy.Level, dummy.CountryCode)
		Players[player.Id()] = player
	}
}

type NewPlayer struct {
	Id          string `json:"id"`
	CountryCode string `json:"country_code"`
	Level       int    `json:"level"`
}

var dummyPlayersJson = `[
	{"id":"alice_smith","country_code":"US","level":7},
	{"id":"bob_jones","country_code":"GB","level":6},
	{"id":"carlos_mendez","country_code":"MX","level":3},
	{"id":"diana_lee","country_code":"KR","level":9},
	{"id":"ethan_wong","country_code":"SG","level":5},
	{"id":"fatima_khan","country_code":"PK","level":8},
	{"id":"george_nash","country_code":"AU","level":4},
	{"id":"hana_tanaka","country_code":"JP","level":6},
	{"id":"ian_clark","country_code":"CA","level":7},
	{"id":"julia_roberts","country_code":"US","level":3},
	{"id":"kevin_singh","country_code":"IN","level":5},
	{"id":"laura_gomez","country_code":"ES","level":2},
	{"id":"mohamed_fahmy","country_code":"EG","level":8},
	{"id":"natalie_cho","country_code":"KR","level":7},
	{"id":"oliver_klein","country_code":"DE","level":6},
	{"id":"paula_martins","country_code":"BR","level":4},
	{"id":"quentin_diaz","country_code":"AR","level":3},
	{"id":"rashid_ali","country_code":"AE","level":7},
	{"id":"sofia_perez","country_code":"CO","level":5},
	{"id":"tommy_nilsen","country_code":"NO","level":6},
	{"id":"ursula_meier","country_code":"CH","level":4},
	{"id":"viktor_ivanov","country_code":"RU","level":9},
	{"id":"wanda_nowak","country_code":"PL","level":2},
	{"id":"xinyi_zhang","country_code":"CN","level":7},
	{"id":"youssef_khaled","country_code":"MA","level":5},
	{"id":"zara_kapoor","country_code":"IN","level":6},
	{"id":"anton_karlsen","country_code":"SE","level":3},
	{"id":"bella_larsson","country_code":"SE","level":4},
	{"id":"cedric_dubois","country_code":"FR","level":5},
	{"id":"daria_smirnova","country_code":"RU","level":7},
	{"id":"emil_petrov","country_code":"BG","level":4},
	{"id":"fiona_kennedy","country_code":"IE","level":6},
	{"id":"gustavo_silva","country_code":"BR","level":7},
	{"id":"harper_clarke","country_code":"CA","level":3},
	{"id":"isabel_nunez","country_code":"ES","level":5},
	{"id":"jack_yamamoto","country_code":"JP","level":6},
	{"id":"karim_rahman","country_code":"BD","level":2},
	{"id":"lucia_ferrari","country_code":"IT","level":6},
	{"id":"matteo_gallo","country_code":"IT","level":7},
	{"id":"nadine_bakker","country_code":"NL","level":8},
	{"id":"oscar_svensson","country_code":"SE","level":5},
	{"id":"penelope_dunne","country_code":"NZ","level":3},
	{"id":"qiang_liu","country_code":"CN","level":7},
	{"id":"raj_patel","country_code":"IN","level":4},
	{"id":"samantha_brown","country_code":"US","level":6},
	{"id":"thomas_muller","country_code":"DE","level":5},
	{"id":"ulrik_berg","country_code":"NO","level":6},
	{"id":"valeria_ramos","country_code":"PE","level":7},
	{"id":"william_owen","country_code":"GB","level":3},
	{"id":"xiaoli_chen","country_code":"CN","level":5},
	{"id":"yuki_sato","country_code":"JP","level":8},
	{"id":"zeynep_aydin","country_code":"TR","level":9},
	{"id":"aaron_clark","country_code":"US","level":4},
	{"id":"bianca_santos","country_code":"BR","level":6},
	{"id":"cameron_wilson","country_code":"GB","level":5},
	{"id":"daniela_rodriguez","country_code":"CL","level":4},
	{"id":"emre_ozdemir","country_code":"TR","level":3},
	{"id":"farah_abdul","country_code":"SA","level":7},
	{"id":"goran_jovanovic","country_code":"RS","level":4},
	{"id":"hana_elmadi","country_code":"DZ","level":6},
	{"id":"ibrahim_sow","country_code":"SN","level":8},
	{"id":"johan_nilsson","country_code":"SE","level":5},
	{"id":"karla_ortega","country_code":"MX","level":4},
	{"id":"liam_brennan","country_code":"IE","level":3},
	{"id":"maya_cohen","country_code":"IL","level":6},
	{"id":"neha_sharma","country_code":"IN","level":7},
	{"id":"otto_meier","country_code":"DE","level":5},
	{"id":"pia_koskinen","country_code":"FI","level":6},
	{"id":"quentin_martel","country_code":"FR","level":7},
	{"id":"rohan_desai","country_code":"IN","level":3},
	{"id":"sara_haddad","country_code":"MA","level":8},
	{"id":"timo_ahlberg","country_code":"FI","level":5},
	{"id":"uliana_kuznetsova","country_code":"RU","level":7},
	{"id":"victoria_lopes","country_code":"BR","level":4},
	{"id":"wassim_hamdi","country_code":"TN","level":6},
	{"id":"ximena_vega","country_code":"PE","level":2},
	{"id":"yasir_khan","country_code":"PK","level":9},
	{"id":"zsolt_szabo","country_code":"HU","level":5},
	{"id":"adriana_dumitru","country_code":"RO","level":6},
	{"id":"benjamin_hoffman","country_code":"DE","level":4},
	{"id":"chloe_marchand","country_code":"FR","level":3},
	{"id":"diego_fuentes","country_code":"UY","level":5},
	{"id":"eva_meyer","country_code":"CH","level":6},
	{"id":"felix_morales","country_code":"EC","level":7},
	{"id":"grace_yip","country_code":"SG","level":4},
	{"id":"hugo_leblanc","country_code":"CA","level":5},
	{"id":"ines_fernandez","country_code":"ES","level":6},
	{"id":"james_ryan","country_code":"US","level":8},
	{"id":"khalid_nasser","country_code":"JO","level":7},
	{"id":"leila_saidi","country_code":"DZ","level":4},
	{"id":"marko_todorovic","country_code":"RS","level":5},
	{"id":"nadia_samir","country_code":"EG","level":6},
	{"id":"omar_hassan","country_code":"SD","level":3},
	{"id":"patricia_silva","country_code":"PT","level":5},
	{"id":"ricardo_garcia","country_code":"MX","level":6},
	{"id":"salma_abdel","country_code":"EG","level":7},
	{"id":"tarek_mansour","country_code":"LB","level":4},
	{"id":"ulric_weber","country_code":"DE","level":6},
	{"id":"valentina_mora","country_code":"CL","level":5},
	{"id":"wenjie_li","country_code":"CN","level":7},
	{"id":"xander_blake","country_code":"US","level":3},
	{"id":"yasmin_almawi","country_code":"SA","level":6},
	{"id":"zachary_reed","country_code":"US","level":4}
]`
