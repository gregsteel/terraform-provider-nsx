package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/edgegateway"
	//"strconv"
)

func resourceEdgeGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceEdgeGatewayCreate,
		Read:   resourceEdgeGatewayRead,
		Delete: resourceEdgeGatewayDelete,
		Update: resourceEdgeGatewayUpdate,

		Schema: map[string]*schema.Schema{
			"datacenter_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				Default:  "default",
			},
			"fqdn": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"vse_log_level": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				Default:  "info",
			},
			"enable_aesni": &schema.Schema{
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  true,
			},
			"enable_fips": &schema.Schema{
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  true,
			},
			"enable_core_dump": &schema.Schema{
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  false,
			},
			"appliance_size": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"appliances": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_pool_id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"datastore_id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"vmfolder_id": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Optional: true,
						},
					},
				},
			},
			"vnics": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Optional: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Optional: true,
						},
						"mtu": &schema.Schema{
							Type:     schema.TypeInt,
							Required: false,
							Optional: true,
						},
						"portgroup_id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"isconnected": &schema.Schema{
							Type:     schema.TypeBool,
							Required: false,
							Optional: true,
						},
						"index": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"addressgroups": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"primaryaddress": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"subnetmask": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildAppliances(appliances []interface{}) []edgegateway.EdgeGatewayAppliance {
	var applianceList []edgegateway.EdgeGatewayAppliance
	for _, appliance := range appliances {
		data := appliance.(map[string]interface{})
		appl := edgegateway.EdgeGatewayAppliance{
			ResourcePoolId: data["resource_pool_id"].(string),
			DatastoreId:    data["datastore_id"].(string),
			VmFolderId:     data["vmfolder_id"].(string),
		}
		applianceList = append(applianceList, appl)
	}
	return applianceList
}

func buildVnicAddressGroups(addressGroups []interface{}) []edgegateway.AddressGroup {
	var addressGroupList []edgegateway.AddressGroup
	for _, address := range addressGroups {
		data := address.(map[string]interface{})
		addr := edgegateway.AddressGroup{
			PrimaryAddress: data["primaryaddress"].(string),
			SubnetMask:     data["subnetmask"].(string),
		}
		addressGroupList = append(addressGroupList, addr)
	}
	return addressGroupList
}

func buildVnics(vnics []interface{}) []edgegateway.EdgeVnic {
	var vnicList []edgegateway.EdgeVnic
	for _, vnic := range vnics {
		data := vnic.(map[string]interface{})
		addrGrps := edgegateway.AddressGroups{
			AddressGroups: buildVnicAddressGroups(data["addressgroups"].([]interface{})),
		}
		vn := edgegateway.EdgeVnic{
			Name:          data["name"].(string),
			Type:          data["type"].(string),
			Mtu:           data["mtu"].(int),
			PortgroupId:   data["portgroup_id"].(string),
			IsConnected:   data["isconnected"].(bool),
			Index:         data["index"].(int),
			AddressGroups: addrGrps,
		}
		vnicList = append(vnicList, vn)
	}
	return vnicList
}

func resourceEdgeGatewayCreate(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)

	var edge edgegateway.EdgeGateway

	edge.DatacenterMoid = d.Get("datacenter_id").(string)
	edge.Name = d.Get("name").(string)
	edge.Type = "gatewayServices" //hardcoded for Edge Gateway
	edge.Tenant = d.Get("tenant").(string)
	edge.Fqdn = d.Get("fqdn").(string)
	edge.VseLogLevel = d.Get("vse_log_level").(string)

	if v, ok := d.GetOk("enable_aesni"); ok {
		edge.EnableAesni = v.(bool)
	} else {
		edge.EnableAesni = false
	}
	if v, ok := d.GetOk("enable_fips"); ok {
		edge.EnableFips = v.(bool)
	} else {
		edge.EnableFips = false
	}

	var appliances edgegateway.EdgeGatewayAppliances
	edge.Appliances = appliances
	if v, ok := d.GetOk("enable_core_dump"); ok {
		edge.Appliances.EnableCoreDump = v.(bool)
	} else {
		edge.Appliances.EnableCoreDump = false
	}
	edge.Appliances.ApplianceSize = d.Get("appliance_size").(string)

	if v, ok := d.GetOk("appliances"); ok {
		edge.Appliances.ApplianceList = buildAppliances(v.([]interface{}))
	}
	if v, ok := d.GetOk("vnics"); ok {
		edge.Vnics.VnicList = buildVnics(v.([]interface{}))
	}

	createAPI := edgegateway.NewCreate(&edge)
	err := nsxclient.Do(createAPI)
	if err != nil {
		return err
	}

	createResponseCode := createAPI.StatusCode()
	createResponse := createAPI.GetResponse()
	if createResponseCode != http.StatusCreated {
		return fmt.Errorf("Error while creating Edge Gateway %s. Invalid HTTP response code %d received. Response: %v", edge.Name, createResponseCode, createResponse)
	}

	//create edge doesn't return the ID, need to get it from the Location
	//resp := createAPI.RawResponse()
	//log.Println("raw response for edgegateway: %v", resp)

	//elements := strings.Split(createAPI.Endpoint(), "/")
	d.SetId(createResponse)
	return resourceEdgeGatewayRead(d, m)
}

func resourceEdgeGatewayRead(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)

	edgegatewayid := d.Id()

	log.Println("Looking for edgegatewayid: %s", edgegatewayid)

	api := edgegateway.NewGet(edgegatewayid)
	err := nsxclient.Do(api)
	if err != nil {
		d.SetId("")
		return nil
	}

	if api.StatusCode() != http.StatusOK {
		d.SetId("")
		return fmt.Errorf("Error getting edgegateway s: Status code: %d", api.StatusCode())
	}

	edgegateway := api.GetResponse()

	d.Set("datacenter_id", edgegateway.DatacenterMoid)
	d.Set("name", edgegateway.Name)
	d.Set("type", edgegateway.Type)
	d.Set("tenant", edgegateway.Tenant)
	d.Set("fqdn", edgegateway.Fqdn)
	d.Set("enable_aesni", edgegateway.EnableAesni)
	d.Set("enable_fips", edgegateway.EnableFips)
	d.Set("vse_log_level", edgegateway.VseLogLevel)
	d.Set("appliance_size", edgegateway.Appliances.ApplianceSize)
	d.Set("enable_core_dump", edgegateway.Appliances.EnableCoreDump)
	var applianceList []*schema.ResourceData
	for _, appliance := range edgegateway.Appliances.ApplianceList {
		var applianceResource schema.ResourceData
		applianceResource.Set("resource_pool_id", appliance.ResourcePoolId)
		applianceResource.Set("datastore_id", appliance.DatastoreId)
		applianceResource.Set("vmfolder_id", appliance.VmFolderId)
		applianceList = append(applianceList, &applianceResource)
	}
	d.Set("appliances", applianceList)
	var vnicList []*schema.ResourceData
	for _, vnic := range edgegateway.Vnics.VnicList {
		var vnicResource schema.ResourceData
		vnicResource.Set("name", vnic.Name)
		vnicResource.Set("type", vnic.Type)
		vnicResource.Set("mtu", vnic.Mtu)
		vnicResource.Set("portgroup_id", vnic.PortgroupId)
		vnicResource.Set("isconnected", vnic.IsConnected)
		vnicResource.Set("index", vnic.Index)
		var adrGroupList []*schema.ResourceData
		for _, adrGroup := range vnic.AddressGroups.AddressGroups {
			var adrgrpResource schema.ResourceData
			adrgrpResource.Set("primaryaddress", adrGroup.PrimaryAddress)
			adrgrpResource.Set("subnetmask", adrGroup.SubnetMask)
			adrGroupList = append(adrGroupList, &adrgrpResource)
		}
		vnicResource.Set("addressgroups", adrGroupList)
		vnicList = append(vnicList, &vnicResource)
	}
	d.Set("vnics", vnicList)

	return nil
}

func resourceEdgeGatewayDelete(d *schema.ResourceData, m interface{}) error {

	nsxclient := m.(*gonsx.NSXClient)
	edgegatewayid := d.Id()

	deleteAPI := edgegateway.NewDelete(edgegatewayid)
	err := nsxclient.Do(deleteAPI)
	if err != nil {
		return err
	}
	if deleteAPI.StatusCode() == http.StatusNoContent {
		d.SetId("")
		return nil
	}
	return fmt.Errorf("Error deleting  from NSX, Status code: %d", deleteAPI.StatusCode())
}

func resourceEdgeGatewayUpdate(d *schema.ResourceData, m interface{}) error {

	nsxclient := m.(*gonsx.NSXClient)
	hasChanges := false

	var updatedEdgeGateway edgegateway.EdgeGateway

	edgegatewayid := d.Id()

	oldName, newName := d.GetChange("name")
	if d.HasChange("name") {
		hasChanges = true
		updatedEdgeGateway.Name = newName.(string)
	} else {
		updatedEdgeGateway.Name = oldName.(string)
	}

	//work in progress - add more attributes

	if hasChanges {
		updateAPI := edgegateway.NewUpdate(edgegatewayid, updatedEdgeGateway)

		nsxMutexKV.Lock(edgegatewayid)
		defer nsxMutexKV.Unlock(edgegatewayid)

		err := nsxclient.Do(updateAPI)
		if err != nil {
			return err
		}

		if updateAPI.StatusCode() != http.StatusNoContent {
			return fmt.Errorf("Error updating resource: status code: %d", updateAPI.StatusCode())
		}
	}

	return nil
}
