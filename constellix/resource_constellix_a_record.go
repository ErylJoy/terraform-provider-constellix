package constellix

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/Constellix/constellix-go-client/client"
	"github.com/Constellix/constellix-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceConstellixARecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceConstellixARecordCreate,
		Update: resourceConstellixARecordUpdate,
		Read:   resourceConstellixARecordRead,
		Delete: resourceConstellixARecordDelete,

		Importer: &schema.ResourceImporter{
			State: resourceConstellixARecordImport,
		},

		Schema: map[string]*schema.Schema{
			"domain_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ttl": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			"source_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"geo_location": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"geo_ip_user_region": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"drop": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"geo_ip_proximity": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"geo_ip_failover": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"record_option": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"noanswer": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"note": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"gtd_region": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"contact_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},

			"roundrobin": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"disable_flag": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Optional: true,
				Computed: true,
			},
			"roundrobin_failover": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"disable_flag": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"sort_order": {
							Type:     schema.TypeString,
							Required: true,
						},
						"check_id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
				Optional: true,
				Computed: true,
			},

			"record_failover_values": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"check_id": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"sort_order": {
							Type:     schema.TypeString,
							Required: true,
						},
						"disable_flag": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Optional: true,
			},

			"record_failover_failover_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"record_failover_disable_flag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pools": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func resourceConstellixARecordImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	constellixClient := m.(*client.Client)
	params := strings.Split(d.Id(), ":")
	resp, err := constellixClient.GetbyId("v1/" + params[0] + "/" + params[1] + "/records/a/" + params[2])
	if err != nil {
		if resp.StatusCode == 404 {
			d.SetId("")
			return nil, err
		}
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bodyString := string(bodyBytes)
	var data map[string]interface{}
	err = json.Unmarshal([]byte(bodyString), &data)
	if err != nil {
		return nil, err
	}

	geoloc1 := data["geolocation"]
	log.Println("GEOLOC VALUE INSIDE READ :", geoloc1)

	geoLocMap := make(map[string]interface{})
	if geoloc1 != nil {
		geoloc := geoloc1.(map[string]interface{})
		if geoloc["geoipFilter"] != nil {
			geoLocMap["geo_ip_user_region"] = fmt.Sprintf("%v", geoloc["geoipFilter"])
		}
		if geoloc["drop"] != nil {
			geoLocMap["drop"] = fmt.Sprintf("%v", geoloc["drop"])
		}
		if geoloc["geoipFailover"] != nil {
			geoLocMap["geo_ip_failover"] = fmt.Sprintf("%v", geoloc["geoipFailover"])
		}
		if geoloc["geoipProximity"] != nil {
			geoLocMap["geo_ip_proximity"] = fmt.Sprintf("%v", geoloc["geoipProximity"])
		}
		d.Set("geo_location", geoLocMap)
	} else {
		d.Set("geo_location", geoLocMap)
	}

	arecroundrobin := data["roundRobin"].([]interface{})
	rrlist := make([]interface{}, 0, 1)
	for _, valrrf := range arecroundrobin {
		map1 := make(map[string]interface{})
		val1 := valrrf.(map[string]interface{})
		map1["value"] = fmt.Sprintf("%v", val1["value"])
		map1["disable_flag"] = fmt.Sprintf("%v", val1["disableFlag"])

		rrlist = append(rrlist, map1)
	}
	log.Printf("uniiiiii %v", rrlist)

	arecroundrobinfailover := data["roundRobinFailover"].([]interface{})

	rrflist := make([]interface{}, 0, 1)
	for _, valrrf := range arecroundrobinfailover {
		map1 := make(map[string]interface{})
		val1 := valrrf.(map[string]interface{})
		map1["value"] = fmt.Sprintf("%v", val1["value"])
		map1["sort_order"] = fmt.Sprintf("%v", val1["sortOrder"])
		map1["disable_flag"] = fmt.Sprintf("%v", val1["disableFlag"])

		rrflist = append(rrflist, map1)
	}

	rcdf := data["recordFailover"]
	rcdflist := make([]interface{}, 0, 1)
	if rcdf != nil {
		rcdf1 := rcdf.(map[string]interface{})
		d.Set("record_failover_failover_type", fmt.Sprintf("%v", rcdf1["failoverType"]))
		d.Set("record_failover_disable_flag", fmt.Sprintf("%v", rcdf1["disabled"]))

		rcdfvalues := rcdf1["values"].([]interface{})

		for _, valrcdf := range rcdfvalues {
			map1 := make(map[string]interface{})
			val1 := valrcdf.(map[string]interface{})
			map1["value"] = fmt.Sprintf("%v", val1["value"])
			map1["sort_order"] = fmt.Sprintf("%v", val1["sortOrder"])
			map1["disable_flag"] = fmt.Sprintf("%v", val1["disableFlag"])
			map1["check_id"] = fmt.Sprintf("%v", val1["checkId"])
			rcdflist = append(rcdflist, map1)
		}
	}

	d.SetId(fmt.Sprintf("%.0f", data["id"]))
	d.Set("name", data["name"])
	d.Set("domain_id", params[1])
	d.Set("source_type", params[0])
	d.Set("ttl", data["ttl"])
	d.Set("record_option", data["recordOption"])
	d.Set("noanswer", data["noAnswer"])
	d.Set("note", data["note"])
	d.Set("gtd_region", data["gtdRegion"])
	d.Set("type", data["type"])
	d.Set("pools", data["pools"])
	d.Set("contact_ids", data["contactIds"])
	d.Set("roundrobin", rrlist)
	d.Set("roundrobin_failover", rrflist)
	d.Set("record_failover_values", rcdflist)

	log.Printf("[DEBUG] %s finished import", d.Id())
	return []*schema.ResourceData{d}, nil

}

func resourceConstellixARecordCreate(d *schema.ResourceData, m interface{}) error {

	constellixConnect := m.(*client.Client)

	aAttr := models.ARecordAttributes{}

	if name, ok := d.GetOk("name"); ok {
		aAttr.Name = name.(string)
	}
	if TTL, ok := d.GetOk("ttl"); ok {
		aAttr.TTL = TTL.(int)
	}
	if RecordOption, ok := d.GetOk("record_option"); ok {
		aAttr.RecordOption = RecordOption.(string)
	}
	if NoAnswer, ok := d.GetOk("noanswer"); ok {
		aAttr.NoAnswer = NoAnswer.(bool)
	}
	if Note, ok := d.GetOk("note"); ok {
		aAttr.Note = Note.(string)
	}
	if GtdRegion, ok := d.GetOk("gtd_region"); ok {
		aAttr.GtdRegion = GtdRegion.(int)
	}
	if Type, ok := d.GetOk("type"); ok {
		aAttr.Type = Type.(string)
	}
	if contactid, ok := d.GetOk("contact_ids"); ok {
		aAttr.ContactId = toListOfInt(contactid)
	}
	if pools, ok := d.GetOk("pools"); ok {
		aAttr.Pools = toListOfInt(pools)
	}

	geoloc := &models.GeolocationArecord{}
	if geoipuserregion, ok := d.GetOk("geo_location"); ok {
		geouserlist := make([]int, 0, 1)
		tp := geoipuserregion.(map[string]interface{})
		if tp["geo_ip_user_region"] != nil {
			var1, _ := strconv.Atoi(fmt.Sprintf("%v", tp["geo_ip_user_region"]))
			geouserlist = append(geouserlist, var1)
			geoloc.GeoIpUserRegion = geouserlist
		}
		if tp["drop"] != nil {
			geoloc.Drop, _ = strconv.ParseBool(fmt.Sprintf("%v", tp["drop"]))
		}
		if tp["geo_ip_failover"] != nil {
			geoloc.GeoIpFailOver, _ = strconv.ParseBool(fmt.Sprintf("%v", tp["geo_ip_failover"]))
		}
		if tp["geo_ip_proximity"] != nil {
			geoloc.GeoIpProximity, _ = strconv.Atoi(fmt.Sprintf("%v", tp["geo_ip_proximity"]))
		}
	}
	aAttr.GeoLocation = geoloc

	maplistrr := make([]interface{}, 0, 1)
	if val, ok := d.GetOk("roundrobin"); ok {
		tp := val.(*schema.Set).List()
		for _, val := range tp {
			map1 := make(map[string]interface{})
			inner := val.(map[string]interface{})
			map1["value"] = fmt.Sprintf("%v", inner["value"])
			map1["disableFlag"], _ = strconv.ParseBool(fmt.Sprintf("%v", inner["disable_flag"]))
			maplistrr = append(maplistrr, map1)
		}
		aAttr.RoundRobin = maplistrr
	}

	maplist := make([]interface{}, 0, 1)
	if value, ok := d.GetOk("roundrobin_failover"); ok {
		tp := value.(*schema.Set).List()

		for _, val := range tp {
			map1 := make(map[string]interface{})
			inner := val.(map[string]interface{})
			map1["value"] = fmt.Sprintf("%v", inner["value"])
			map1["disableFlag"], _ = strconv.ParseBool(fmt.Sprintf("%v", inner["disable_flag"]))
			map1["sortOrder"], _ = strconv.Atoi(fmt.Sprintf("%v", inner["sort_order"]))
			map1["checkId"], _ = strconv.Atoi(fmt.Sprintf("%v", inner["check_id"]))
			maplist = append(maplist, map1)
		}
		aAttr.RoundRobinFailoverA = sortAccordingToSortOrder(maplist)
	}

	valueslist := make([]interface{}, 0, 1)
	if value, ok := d.GetOk("record_failover_values"); ok {
		rcdfa := &models.RCDFAARecord{} //added
		tp := value.(*schema.Set).List()
		for _, val := range tp {
			map1 := make(map[string]interface{})
			inner := val.(map[string]interface{})
			map1["checkId"] = fmt.Sprintf("%v", inner["check_id"])
			map1["value"] = fmt.Sprintf("%v", inner["value"])
			map1["sortOrder"], _ = strconv.Atoi(fmt.Sprintf("%v", inner["sort_order"]))
			map1["disableFlag"], _ = strconv.ParseBool(fmt.Sprintf("%v", inner["disable_flag"]))
			valueslist = append(valueslist, map1)
		}

		if failovertype, ok := d.GetOk("record_failover_failover_type"); ok {
			rcdfa.FailoverTypeRCDFA, _ = strconv.Atoi(fmt.Sprintf("%v", failovertype)) //added
		}

		if disableflag, ok := d.GetOk("record_failover_disable_flag"); ok {
			rcdfa.DisableFlagRCDFA, _ = strconv.ParseBool(fmt.Sprintf("%v", disableflag)) //added
		}

		rcdfa.Values = sortAccordingToSortOrder(valueslist) //added
		aAttr.RecordFailoverA = rcdfa                       //added
	}

	resp, err := constellixConnect.Save(aAttr, "v1/"+d.Get("source_type").(string)+"/"+d.Get("domain_id").(string)+"/records/a")

	if err != nil {
		return err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	bodyString := string(bodyBytes)
	var data map[string]interface{}
	json.Unmarshal([]byte(bodyString[1:len(bodyString)-1]), &data)

	d.SetId(fmt.Sprintf("%.0f", data["id"]))

	return resourceConstellixARecordRead(d, m)
}

func resourceConstellixARecordRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Begining Read %s", d.Id())
	constellixClient := m.(*client.Client)
	arecordid := d.Id()

	resp, err := constellixClient.GetbyId("v1/" + d.Get("source_type").(string) + "/" + d.Get("domain_id").(string) + "/records/a/" + arecordid)
	if err != nil {
		if resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	bodyString := string(bodyBytes)
	var data map[string]interface{}
	err = json.Unmarshal([]byte(bodyString), &data)
	if err != nil {
		return err
	}

	geoloc1 := data["geolocation"]
	log.Println("GEOLOC VALUE INSIDE READ :", geoloc1)

	geoLocMap := make(map[string]interface{})
	if geoloc1 != nil {
		geoloc := geoloc1.(map[string]interface{})
		if geoloc["geoipFilter"] != nil {
			geoLocMap["geo_ip_user_region"] = fmt.Sprintf("%v", geoloc["geoipFilter"])
		}
		if geoloc["drop"] != nil {
			geoLocMap["drop"] = fmt.Sprintf("%v", geoloc["drop"])
		}
		if geoloc["geoipFailover"] != nil {
			geoLocMap["geo_ip_failover"] = fmt.Sprintf("%v", geoloc["geoipFailover"])
		}
		if geoloc["geoipProximity"] != nil {
			geoLocMap["geo_ip_proximity"] = fmt.Sprintf("%v", geoloc["geoipProximity"])
		}
		d.Set("geo_location", geoLocMap)
	} else {
		d.Set("geo_location", geoLocMap)
	}

	arecroundrobin := data["roundRobin"].([]interface{})
	rrlist := make([]interface{}, 0, 1)
	for _, valrrf := range arecroundrobin {
		map1 := make(map[string]interface{})
		val1 := valrrf.(map[string]interface{})
		map1["value"] = fmt.Sprintf("%v", val1["value"])
		map1["disable_flag"] = fmt.Sprintf("%v", val1["disableFlag"])

		rrlist = append(rrlist, map1)
	}

	arecroundrobinfailover := data["roundRobinFailover"].([]interface{})

	rrflist := make([]interface{}, 0, 1)
	for _, valrrf := range arecroundrobinfailover {
		map1 := make(map[string]interface{})
		val1 := valrrf.(map[string]interface{})
		map1["value"] = fmt.Sprintf("%v", val1["value"])
		map1["sort_order"] = fmt.Sprintf("%v", val1["sortOrder"])
		map1["disable_flag"] = fmt.Sprintf("%v", val1["disableFlag"])
		map1["check_id"] = fmt.Sprintf("%v", val1["checkId"])

		rrflist = append(rrflist, map1)
	}

	rcdf := data["recordFailover"]
	rcdflist := make([]interface{}, 0, 1)
	if rcdf != nil {
		rcdf1 := rcdf.(map[string]interface{})
		d.Set("record_failover_failover_type", fmt.Sprintf("%v", rcdf1["failoverType"]))
		d.Set("record_failover_disable_flag", fmt.Sprintf("%v", rcdf1["disabled"]))

		rcdfvalues := rcdf1["values"].([]interface{})

		for _, valrcdf := range rcdfvalues {
			map1 := make(map[string]interface{})
			val1 := valrcdf.(map[string]interface{})
			map1["value"] = fmt.Sprintf("%v", val1["value"])
			map1["sort_order"] = fmt.Sprintf("%v", val1["sortOrder"])
			map1["disable_flag"] = fmt.Sprintf("%v", val1["disableFlag"])
			map1["check_id"] = fmt.Sprintf("%v", val1["checkId"])
			rcdflist = append(rcdflist, map1)
		}
	}

	d.SetId(fmt.Sprintf("%.0f", data["id"]))
	d.Set("name", data["name"])
	d.Set("ttl", data["ttl"])
	d.Set("record_option", data["recordOption"])
	d.Set("noanswer", data["noAnswer"])
	d.Set("note", data["note"])
	d.Set("gtd_region", data["gtdRegion"])
	d.Set("type", data["type"])
	d.Set("pools", data["pools"])
	d.Set("contact_ids", data["contactIds"])
	d.Set("roundrobin", rrlist)
	d.Set("roundrobin_failover", rrflist)
	d.Set("record_failover_values", rcdflist)

	return nil
}

func resourceConstellixARecordUpdate(d *schema.ResourceData, m interface{}) error {
	constellixClient := m.(*client.Client)
	aAttr := models.ARecordAttributes{}

	if ttl, ok := d.GetOk("ttl"); ok {
		aAttr.TTL = ttl.(int)
	}
	if name, ok := d.GetOk("name"); ok {
		aAttr.Name = name.(string)
	}

	if recordoption, ok := d.GetOk("record_option"); ok {
		aAttr.RecordOption = recordoption.(string)
	}

	if _, ok := d.GetOk("noanswer"); ok {
		aAttr.NoAnswer = d.Get("noanswer").(bool)
	}

	if note, ok := d.GetOk("note"); ok {
		aAttr.Note = note.(string)
	}

	if _, ok := d.GetOk("gtd_region"); ok {
		aAttr.GtdRegion = d.Get("gtd_region").(int)
	}

	if _, ok := d.GetOk("type"); ok {
		aAttr.Type = d.Get("type").(string)
	}
	if contactid, ok := d.GetOk("contact_ids"); ok {
		aAttr.ContactId = toListOfInt(contactid)
	}
	if pools, ok := d.GetOk("pools"); ok {
		aAttr.Pools = toListOfInt(pools)
	}
	geoloc := &models.GeolocationArecord{}
	if geoipuserregion, ok := d.GetOk("geo_location"); ok {
		geouserlist := make([]int, 0, 1)
		tp := geoipuserregion.(map[string]interface{})
		if tp["geo_ip_user_region"] != nil {
			var1, _ := strconv.Atoi(fmt.Sprintf("%v", tp["geo_ip_user_region"]))
			geouserlist = append(geouserlist, var1)
			geoloc.GeoIpUserRegion = geouserlist
		}
		if tp["drop"] != nil {
			geoloc.Drop, _ = strconv.ParseBool(fmt.Sprintf("%v", tp["drop"]))
		}
		if tp["geo_ip_failover"] != nil {
			geoloc.GeoIpFailOver, _ = strconv.ParseBool(fmt.Sprintf("%v", tp["geo_ip_failover"]))
		}
		if tp["geo_ip_proximity"] != nil {
			geoloc.GeoIpProximity, _ = strconv.Atoi(fmt.Sprintf("%v", tp["geo_ip_proximity"]))
		}
	}
	aAttr.GeoLocation = geoloc

	maplistrr := make([]interface{}, 0, 1)
	if val, ok := d.GetOk("roundrobin"); ok {
		tp := val.(*schema.Set).List()
		for _, val := range tp {
			map1 := make(map[string]interface{})
			inner := val.(map[string]interface{})
			map1["value"] = fmt.Sprintf("%v", inner["value"])
			map1["disableFlag"], _ = strconv.ParseBool(fmt.Sprintf("%v", inner["disable_flag"]))
			maplistrr = append(maplistrr, map1)
		}
		aAttr.RoundRobin = maplistrr
	}

	maplist := make([]interface{}, 0, 1)
	if value, ok := d.GetOk("roundrobin_failover"); ok {
		tp := value.(*schema.Set).List()

		for _, val := range tp {
			map1 := make(map[string]interface{})
			inner := val.(map[string]interface{})
			map1["value"] = fmt.Sprintf("%v", inner["value"])
			map1["disableFlag"], _ = strconv.ParseBool(fmt.Sprintf("%v", inner["disable_flag"]))
			map1["sortOrder"], _ = strconv.Atoi(fmt.Sprintf("%v", inner["sort_order"]))
			map1["checkId"], _ = strconv.Atoi(fmt.Sprintf("%v", inner["check_id"]))
			maplist = append(maplist, map1)
		}
		aAttr.RoundRobinFailoverA = sortAccordingToSortOrder(maplist)
	}

	valueslist := make([]interface{}, 0, 1)
	if value, ok := d.GetOk("record_failover_values"); ok {
		rcdfa := &models.RCDFAARecord{} //added
		tp := value.(*schema.Set).List()
		for _, val := range tp {
			map1 := make(map[string]interface{})
			inner := val.(map[string]interface{})
			map1["checkId"] = fmt.Sprintf("%v", inner["check_id"])
			map1["value"] = fmt.Sprintf("%v", inner["value"])
			map1["sortOrder"], _ = strconv.Atoi(fmt.Sprintf("%v", inner["sort_order"]))
			map1["disableFlag"], _ = strconv.ParseBool(fmt.Sprintf("%v", inner["disable_flag"]))
			valueslist = append(valueslist, map1)
		}

		if failovertype, ok := d.GetOk("record_failover_failover_type"); ok {
			rcdfa.FailoverTypeRCDFA, _ = strconv.Atoi(fmt.Sprintf("%v", failovertype))
		}

		if disableflag, ok := d.GetOk("record_failover_disable_flag"); ok {
			rcdfa.DisableFlagRCDFA, _ = strconv.ParseBool(fmt.Sprintf("%v", disableflag))
		}

		rcdfa.Values = sortAccordingToSortOrder(valueslist)
		aAttr.RecordFailoverA = rcdfa
	}

	arecordid := d.Id()

	_, err := constellixClient.UpdatebyID(aAttr, "v1/"+d.Get("source_type").(string)+"/"+d.Get("domain_id").(string)+"/records/a/"+arecordid)
	if err != nil {
		return err
	}
	return resourceConstellixARecordRead(d, m)

}

func resourceConstellixARecordDelete(d *schema.ResourceData, m interface{}) error {
	constellixClient := m.(*client.Client)
	arecordid := d.Id()

	err := constellixClient.DeletebyId("v1/" + d.Get("source_type").(string) + "/" + d.Get("domain_id").(string) + "/records/a/" + arecordid)
	if err != nil {
		return err
	}
	d.SetId("")
	return err
}

func sortAccordingToSortOrder(mapList []interface{}) []interface{} {
	sortOrders := make([]int, 0)
	sortOrdersMap := make(map[int]bool, 0)
	for _, dict := range mapList {
		sortOrdersMap[dict.(map[string]interface{})["sortOrder"].(int)] = true
	}

	for k, _ := range sortOrdersMap {
		sortOrders = append(sortOrders, k)
	}

	sort.Ints(sortOrders)

	sortedMapList := make([]interface{}, 0)
	for _, order := range sortOrders {
		for i := 0; i < len(mapList); i++ {
			if mapList[i].(map[string]interface{})["sortOrder"].(int) == order {
				sortedMapList = append(sortedMapList, mapList[i])
			}
		}
	}
	return sortedMapList
}
