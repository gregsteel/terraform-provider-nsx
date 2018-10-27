package main

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/virtualwire"
)

func dataSourceLogicalSwitch() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLogicalSwitchRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the logical switch",
			},
			"scopeid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The transport zone ID.",
			},
		},
	}
}

func dataSourceLogicalSwitchRead(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*gonsx.NSXClient)
	name := d.Get("name").(string)
	scopeid := d.Get("scopeid").(string)

	getAllAPI := virtualwire.NewGetAll(scopeid)
	err := nsxClient.Do(getAllAPI)
	if err != nil {
		return fmt.Errorf("Error while reading logical switchs. Error: %v", err)
	}

	logicalswitch := getAllAPI.GetResponse().FilterByName(name)
	d.SetId(logicalswitch.ObjectID)

	return nil
}
