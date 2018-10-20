package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/tzone"
)

func dataSourceTransportZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTransportZoneRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the transport zone",
			},
			"desc": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A description of the transport zone",
			},
			"controlplanemode": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The control plane mode to use with the transport zone. Typically this will be UNICAST_MODE. Other valid options are HYBRID_MODE and MULTICAST_MODE",
			},
		},
	}
}

func dataSourceTransportZoneRead(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*gonsx.NSXClient)
	name := d.Get("name").(string)

	getAllAPI := tzone.NewGetAll()
	err := nsxClient.Do(getAllAPI)
	if err != nil {
		return fmt.Errorf("Error while reading transport zones. Error: %v", err)
	}

	TransportZone := getAllAPI.GetResponse().FilterByName(name)
	d.SetId(TransportZone.ObjectID)
	d.Set("desc", TransportZone.Description)
	d.Set("controlplanemode", TransportZone.ControlPlaneMode)

	return nil
}
