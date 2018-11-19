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
			"firewall_default_policy_action": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				Default:  "accept",
			},
			"firewall_default_policy_logging_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  false,
			},
			"firewall_tcp_timeout_established": &schema.Schema{
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
				Default:  300000,
			},
			"router_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"router_ecmp": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"static_routing_default_vnic": &schema.Schema{
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
			},
			"static_routing_default_gateway": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"static_routing_default_mtu": &schema.Schema{
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
			},
			"routing_ospf_areas": &schema.Schema{
				Type:     schema.TypeList,
				Required: false,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"area_id": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
							Optional: false,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Optional: true,
							Default:  "normal",
						},
					},
				},
			},
			"routing_ospf_interfaces": &schema.Schema{
				Type:     schema.TypeList,
				Required: false,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"area_id": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
							Optional: false,
						},
						"vnic_id": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
							Optional: false,
						},
					},
				},
			},
			"routing_ospf_redistribution_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  false,
			},
			"routing_ospf_graceful_restart": &schema.Schema{
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  true,
			},
			"routing_ospf_default_originate": &schema.Schema{
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  false,
			},
			"dhcp_static_bindings": &schema.Schema{
				Type:     schema.TypeList,
				Required: false,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"macaddress": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Optional: true,
						},
						"vm_id": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Optional: true,
						},
						"vnic_id": &schema.Schema{
							Type:     schema.TypeInt,
							Required: false,
							Optional: true,
						},
						"hostname": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Optional: true,
						},
						"ipaddress": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Optional: true,
						},
						"subnetmask": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Optional: true,
						},
						"defaultgateway": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Optional: true,
						},
						"domainname": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Optional: true,
						},
						"primarydns": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Optional: true,
						},
						"secondarydns": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Optional: true,
						},
						"lease": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Optional: true,
						},
						"autoconfigure_dns": &schema.Schema{
							Type:     schema.TypeBool,
							Required: false,
							Optional: true,
						},
					},
				},
			},
			"dhcp_ip_pools": &schema.Schema{
				Type:     schema.TypeList,
				Required: false,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"iprange": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Optional: true,
						},
						"defaultgateway": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Optional: true,
						},
						"subnetmask": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Optional: true,
						},
						"domainname": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Optional: true,
						},
						"primarydns": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Optional: true,
						},
						"secondarydns": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Optional: true,
						},
						"lease": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Optional: true,
						},
						"autoconfigure_dns": &schema.Schema{
							Type:     schema.TypeBool,
							Required: false,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func buildAppliances(appliances []interface{}) []edgegateway.Appliance {
	var applianceList []edgegateway.Appliance
	for _, appliance := range appliances {
		data := appliance.(map[string]interface{})
		appl := edgegateway.Appliance{
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

func buildOSPFInterfaces(ospfInterfaces []interface{}) []edgegateway.OSPFInterface {
	var ospfInterfaceList []edgegateway.OSPFInterface
	for _, ospfInterface := range ospfInterfaces {
		data := ospfInterface.(map[string]interface{})
		intf := edgegateway.OSPFInterface{
			Vnic:   data["vnic_id"].(int),
			AreaId: data["area_id"].(int),
		}
		ospfInterfaceList = append(ospfInterfaceList, intf)
	}
	return ospfInterfaceList
}

func buildOSPFAreas(ospfAreas []interface{}) []edgegateway.OSPFArea {
	var ospfAreaList []edgegateway.OSPFArea
	for _, ospfArea := range ospfAreas {
		data := ospfArea.(map[string]interface{})
		area := edgegateway.OSPFArea{
			AreaId: data["area_id"].(int),
			Type:   data["type"].(string),
		}
		ospfAreaList = append(ospfAreaList, area)
	}
	return ospfAreaList
}

func buildDhcpStaticBindings(staticBindings []interface{}) []edgegateway.DhcpStaticBinding {
	var staticBindingList []edgegateway.DhcpStaticBinding
	for _, staticBinding := range staticBindings {
		data := staticBinding.(map[string]interface{})
		sb := edgegateway.DhcpStaticBinding{
			MacAddress:          data["macaddress"].(string),
			VmId:                data["vm_id"].(string),
			VnicId:              data["vnic_id"].(int),
			Hostname:            data["hostname"].(string),
			IpAddress:           data["ipaddress"].(string),
			SubnetMask:          data["subnetmask"].(string),
			DefaultGateway:      data["defaultgateway"].(string),
			DomainName:          data["domainname"].(string),
			PrimaryNameServer:   data["primarydns"].(string),
			SecondaryNameServer: data["secondarydns"].(string),
			LeaseTime:           data["lease"].(string),
			AutoConfigureDNS:    data["autoconfigure_dns"].(bool),
		}
		staticBindingList = append(staticBindingList, sb)
	}
	return staticBindingList
}

func buildDhcpIpPools(ipPools []interface{}) []edgegateway.DhcpIpPool {
	var ipPoolList []edgegateway.DhcpIpPool
	for _, ipPool := range ipPools {
		data := ipPool.(map[string]interface{})
		ipp := edgegateway.DhcpIpPool{
			IpRange:             data["iprange"].(string),
			SubnetMask:          data["subnetmask"].(string),
			DefaultGateway:      data["defaultgateway"].(string),
			DomainName:          data["domainname"].(string),
			PrimaryNameServer:   data["primarydns"].(string),
			SecondaryNameServer: data["secondarydns"].(string),
			LeaseTime:           data["lease"].(string),
			AutoConfigureDNS:    data["autoconfigure_dns"].(bool),
		}
		ipPoolList = append(ipPoolList, ipp)
	}
	return ipPoolList
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

	var appliances edgegateway.Appliances
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

	var features edgegateway.Features
	edge.Features = features

	var firewall edgegateway.Firewall
	edge.Features.Firewall = firewall
	var defaultPolicy edgegateway.FirewallDefaultPolicy
	edge.Features.Firewall.DefaultPolicy = defaultPolicy
	if v, ok := d.GetOk("firewall_default_policy_action"); ok {
		edge.Features.Firewall.DefaultPolicy.Action = v.(string)
	} else {
		edge.Features.Firewall.DefaultPolicy.Action = "accept"
	}
	if v, ok := d.GetOk("firewall_default_policy_logging_enabled"); ok {
		edge.Features.Firewall.DefaultPolicy.LoggingEnabled = v.(bool)
	} else {
		edge.Features.Firewall.DefaultPolicy.LoggingEnabled = false
	}

	var firewallGlobalConfig edgegateway.FirewallGlobalConfig
	edge.Features.Firewall.GlobalConfig = firewallGlobalConfig
	if v, ok := d.GetOk("firewall_tcp_timeout_established"); ok {
		edge.Features.Firewall.GlobalConfig.TcpTimeoutEstablished = v.(int)
	} else {
		edge.Features.Firewall.GlobalConfig.TcpTimeoutEstablished = 300000
	}

	//Routing config is optional
	if v, ok := d.GetOk("router_id"); ok {
		var routing edgegateway.Routing
		edge.Features.Routing = routing
		var routingGlobalConfig edgegateway.RoutingGlobalConfig
		edge.Features.Routing.GlobalConfig = routingGlobalConfig
		edge.Features.Routing.GlobalConfig.RouterId = v.(string)
		if v, ok := d.GetOk("router_ecmp"); ok {
			edge.Features.Routing.GlobalConfig.ECMP = v.(bool)
		} else {
			edge.Features.Routing.GlobalConfig.ECMP = false
		}
		var routingStaticRouting edgegateway.StaticRouting
		edge.Features.Routing.StaticRouting = routingStaticRouting
		var routingStaticRoutingDefault edgegateway.StaticRoutingDefault
		edge.Features.Routing.StaticRouting.DefaultRoute = routingStaticRoutingDefault
		if v, ok := d.GetOk("static_routing_default_vnic"); ok {
			edge.Features.Routing.StaticRouting.DefaultRoute.Vnic = v.(int)
		}
		if v, ok := d.GetOk("static_routing_default_gateway"); ok {
			edge.Features.Routing.StaticRouting.DefaultRoute.GatewayAddress = v.(string)
		}
		if v, ok := d.GetOk("static_routing_default_mtu"); ok {
			edge.Features.Routing.StaticRouting.DefaultRoute.Mtu = v.(int)
		} else {
			edge.Features.Routing.StaticRouting.DefaultRoute.Mtu = 1500
		}

		//if OSPF interfaces are defined
		if v, ok := d.GetOk("routing_ospf_interfaces"); ok {
			var ospfRouting edgegateway.OSPFRouting
			edge.Features.Routing.OSPFRouting = ospfRouting
			var ospfRoutingInterfaces edgegateway.OSPFInterfaces
			var ospfRoutingAreas edgegateway.OSPFAreas
			edge.Features.Routing.OSPFRouting.Enabled = true
			edge.Features.Routing.OSPFRouting.OSPFInterfaces = ospfRoutingInterfaces
			edge.Features.Routing.OSPFRouting.OSPFAreas = ospfRoutingAreas
			edge.Features.Routing.OSPFRouting.OSPFInterfaces.OSPFInterfaceList = buildOSPFInterfaces(v.([]interface{}))
			if v, ok := d.GetOk("routing_ospf_areas"); ok {
				edge.Features.Routing.OSPFRouting.OSPFAreas.OSPFAreaList = buildOSPFAreas(v.([]interface{}))
			}
			if v, ok := d.GetOk("routing_ospf_redistribution_enabled"); ok {
				var ospfRedistribution edgegateway.OSPFRedistribution
				edge.Features.Routing.OSPFRouting.Redistribution = ospfRedistribution
				edge.Features.Routing.OSPFRouting.Redistribution.Enabled = v.(bool)
				edge.Features.Routing.OSPFRouting.Redistribution.Rules = []edgegateway.OSPFRule{}
			}
			if v, ok := d.GetOk("routing_ospf_graceful_restart"); ok {
				edge.Features.Routing.OSPFRouting.GracefulRestart = v.(bool)
			} else {
				edge.Features.Routing.OSPFRouting.GracefulRestart = true
			}
			if v, ok := d.GetOk("routing_ospf_default_originate"); ok {
				edge.Features.Routing.OSPFRouting.DefaultOriginate = v.(bool)
			} else {
				edge.Features.Routing.OSPFRouting.DefaultOriginate = false
			}
		}
	}

	var dhcp edgegateway.Dhcp
	edge.Features.Dhcp = dhcp
	if v, ok := d.GetOk("dhcp_static_bindings"); ok {
		edge.Features.Dhcp.Enabled = true
		var dhcpStaticBindings edgegateway.DhcpStaticBindings
		edge.Features.Dhcp.StaticBindings = dhcpStaticBindings
		edge.Features.Dhcp.StaticBindings.StaticBindingList = buildDhcpStaticBindings(v.([]interface{}))
	}
	if v, ok := d.GetOk("dhcp_ip_pools"); ok {
		edge.Features.Dhcp.Enabled = true
		var dhcpIpPools edgegateway.DhcpIpPools
		edge.Features.Dhcp.IpPools = dhcpIpPools
		edge.Features.Dhcp.IpPools.IpPoolList = buildDhcpIpPools(v.([]interface{}))
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
	d.Set("firewall_default_policy_action", edgegateway.Features.Firewall.DefaultPolicy.Action)
	d.Set("firewall_default_policy_logging_enabled", edgegateway.Features.Firewall.DefaultPolicy.LoggingEnabled)
	d.Set("firewall_tcp_timeout_established", edgegateway.Features.Firewall.GlobalConfig.TcpTimeoutEstablished)

	d.Set("router_id", edgegateway.Features.Routing.GlobalConfig.RouterId)
	d.Set("router_ecmp", edgegateway.Features.Routing.GlobalConfig.ECMP)
	d.Set("static_routing_default_vnic", edgegateway.Features.Routing.GlobalConfig.ECMP)
	var ospfAreaList []*schema.ResourceData
	for _, area := range edgegateway.Features.Routing.OSPFRouting.OSPFAreas.OSPFAreaList {
		var ospfAreaResource schema.ResourceData
		ospfAreaResource.Set("area_id", area.AreaId)
		ospfAreaResource.Set("type", area.Type)
		ospfAreaList = append(ospfAreaList, &ospfAreaResource)
	}
	d.Set("routing_ospf_areas", ospfAreaList)
	var ospfInterfaceList []*schema.ResourceData
	for _, intf := range edgegateway.Features.Routing.OSPFRouting.OSPFInterfaces.OSPFInterfaceList {
		var ospfInterfaceResource schema.ResourceData
		ospfInterfaceResource.Set("vnic_id", intf.Vnic)
		ospfInterfaceResource.Set("area_id", intf.AreaId)
		ospfInterfaceList = append(ospfInterfaceList, &ospfInterfaceResource)
	}
	d.Set("routing_ospf_interfaces", ospfInterfaceList)
	d.Set("routing_ospf_redistribution_enabled", edgegateway.Features.Routing.OSPFRouting.Redistribution.Enabled)
	d.Set("routing_ospf_graceful_restart", edgegateway.Features.Routing.OSPFRouting.GracefulRestart)
	d.Set("routing_ospf_default_originate", edgegateway.Features.Routing.OSPFRouting.DefaultOriginate)

	var dhcpStaticBindingList []*schema.ResourceData
	for _, sb := range edgegateway.Features.Dhcp.StaticBindings.StaticBindingList {
		var dhcpStaticBinding schema.ResourceData
		dhcpStaticBinding.Set("macaddress", sb.MacAddress)
		dhcpStaticBinding.Set("vm_id", sb.VmId)
		dhcpStaticBinding.Set("vnic_id", sb.VnicId)
		dhcpStaticBinding.Set("hostname", sb.Hostname)
		dhcpStaticBinding.Set("ipaddress", sb.IpAddress)
		dhcpStaticBinding.Set("subnetmask", sb.SubnetMask)
		dhcpStaticBinding.Set("defaultgateway", sb.DefaultGateway)
		dhcpStaticBinding.Set("domainname", sb.DomainName)
		dhcpStaticBinding.Set("primarydns", sb.PrimaryNameServer)
		dhcpStaticBinding.Set("secondarydns", sb.SecondaryNameServer)
		dhcpStaticBinding.Set("lease", sb.LeaseTime)
		dhcpStaticBinding.Set("autoconfigure_dns", sb.AutoConfigureDNS)
		dhcpStaticBindingList = append(dhcpStaticBindingList, &dhcpStaticBinding)
	}
	d.Set("dhcp_static_bindings", dhcpStaticBindingList)
	var dhcpIpPoolList []*schema.ResourceData
	for _, sb := range edgegateway.Features.Dhcp.IpPools.IpPoolList {
		var dhcpIpPool schema.ResourceData
		dhcpIpPool.Set("iprange", sb.IpRange)
		dhcpIpPool.Set("subnetmask", sb.SubnetMask)
		dhcpIpPool.Set("defaultgateway", sb.DefaultGateway)
		dhcpIpPool.Set("domainname", sb.DomainName)
		dhcpIpPool.Set("primarydns", sb.PrimaryNameServer)
		dhcpIpPool.Set("secondarydns", sb.SecondaryNameServer)
		dhcpIpPool.Set("lease", sb.LeaseTime)
		dhcpIpPool.Set("autoconfigure_dns", sb.AutoConfigureDNS)
		dhcpIpPoolList = append(dhcpIpPoolList, &dhcpIpPool)
	}
	d.Set("dhcp_ip_pools", dhcpIpPoolList)

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
