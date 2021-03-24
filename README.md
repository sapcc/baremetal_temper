# baremetal_temper

This is an ironic out of bound inspection using redfish, netbox and ironic-inspector.

## what it does
It can automatically create ironic nodes based on the redfish data provided by the redfish api of a server.

To accomplish this, all it needs is the ip and name of the server, which should be onboarded in ironic. 
This data can be provided via a file or the netbox api.   
The data structure looks as follows:
```
type NetboxDiscovery struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}
```

The following list shows the steps executed for each new node:

* creates a dns record for the new node
* loads server information via the redfish api
* creates new node via the ironic inspector api
* applys rules which are defined by the user
* validates the new ironic node
* tries to power on the new node via the ironic api
* puts node in status available
* waits for nova propagation to be finished
* deploys a test instance on this new node via the compute api
* deletes new test instance
* prepares node for customer use (conductor group etc)
* sets node status to active

If any of the above steps fail, the node will be flagged as not ready in netbox and the node will be deleted in ironic.

## extra feature
The user can provide a rules json template which should be applied for each new node, such as a specific node name, port infos etc.  
Example:
```
{
  "properties": {
    "node": [
      {
        "op": "replace",
        "path": "/name",
        "value": "{{ .node.Name }}"
      },
      {
        "op": "add",
        "path": "/driver_info/deploy_kernel",
        "value": "{{ imageToID `some_image_name` }}"
      },
      {
        "op": "replace",
        "path": "/resource_class",
        "value": "{{ getMatchingFlavorFor .node }}"
      }
    ],
    "port": [
      {
        "op": "add",
        "path": "/local_link_connection/switch_id"
        "value": "aa:bb:cc:dd:ee:ff"
      }
    ]
  }
}
```
