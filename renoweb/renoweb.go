package renoweb

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/anderskvist/DVIEnergiSmartControl/log"
)

type addressSearchSearch struct {
	Searchterm          string `json:"searchterm"`
	Addresswithmateriel int    `json:"addresswithmateriel"`
}

// D is a temporary struct to keep data from renoweb
type tempData struct {
	D string `json:"d"`
}

type AddressSearch struct {
	List   []AddressSearchList
	Status AddressSearchStatus
}

type AddressSearchList struct {
	Value int    `json:"value,string"`
	Label string `json:"label"`
}

type AddressSearchStatus struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

type pickupPlanSearch struct {
	Adrid  int  `json:"adrid"`
	Common bool `json:"common"`
}

type PickupPlan struct {
	List []PickupPlanList
}

type PickupPlanList struct {
	ID            int    `json:"id"`
	MaterielNavn  string `json:"materielnavn"`
	OrdningNavn   string `json:"ordningnavn"`
	ToemningsDage string `json:"toemningsdage"`
	ToemningsDato string `json:"toemningsdato"`
}

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}

func GetRenoWebAddressID(needle string) int {
	search := addressSearchSearch{
		Searchterm:          needle,
		Addresswithmateriel: 3}

	jsondata, _ := json.Marshal(search)
	log.Debugf("%s\n", jsonPrettyPrint(string(jsondata)))

	response, _ := http.Post("https://rebild-sb.renoweb.dk/Legacy/JService.asmx/Adresse_SearchByString", "application/json", bytes.NewBuffer(jsondata))
	data, _ := ioutil.ReadAll(response.Body)
	log.Debugf("%s\n", jsonPrettyPrint(string(data)))

	var d tempData
	err := json.Unmarshal(data, &d)
	if err != nil {
		panic(err)
	}
	log.Debugf("%s\n", d.D)

	var addressSearch AddressSearch
	err = json.Unmarshal([]byte(d.D), &addressSearch)
	if err != nil {
		panic(err)
	}
	log.Infof("%#v\n", addressSearch)
	return addressSearch.List[0].Value
}

func GetRenoWebPickupPlan(id int) PickupPlan {
	search := pickupPlanSearch{
		Adrid:  id,
		Common: false}

	jsondata, _ := json.Marshal(search)
	log.Debugf("%s\n", jsonPrettyPrint(string(jsondata)))

	response, _ := http.Post("https://rebild-sb.renoweb.dk/Legacy/JService.asmx/GetAffaldsplanMateriel_mitAffald", "application/json", bytes.NewBuffer(jsondata))
	data, _ := ioutil.ReadAll(response.Body)
	log.Debugf("%s\n", jsonPrettyPrint(string(data)))

	var d tempData
	err := json.Unmarshal(data, &d)
	if err != nil {
		panic(err)
	}
	log.Debugf("%s\n", d.D)

	var pickupPlan PickupPlan
	err = json.Unmarshal([]byte(d.D), &pickupPlan)
	if err != nil {
		panic(err)
	}
	log.Infof("%#v\n", pickupPlan)
	return pickupPlan
}
