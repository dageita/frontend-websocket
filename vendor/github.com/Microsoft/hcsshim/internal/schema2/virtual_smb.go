/*
 * HCS API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 2.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package hcsschema

type VirtualSmb struct {
	Shares []VirtualSmbShare `json:"Shares,omitempty"`

	DirectFileMappingInMB int64 `json:"DirectFileMappingInMB,omitempty"`
}
