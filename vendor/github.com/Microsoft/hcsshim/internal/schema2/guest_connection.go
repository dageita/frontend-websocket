/*
 * HCS API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 2.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package hcsschema

type GuestConnection struct {

	//  Use Vsock rather than Hyper-V sockets to communicate with the guest service.
	UseVsock bool `json:"UseVsock,omitempty"`

	//  Don't disconnect the guest connection when pausing the virtual machine.
	UseConnectedSuspend bool `json:"UseConnectedSuspend,omitempty"`
}
