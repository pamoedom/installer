package slb

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// TCPListenerConfig is a nested struct in slb response
type TCPListenerConfig struct {
	ConnectionDrain           string      `json:"ConnectionDrain" xml:"ConnectionDrain"`
	ConnectionDrainTimeout    int         `json:"ConnectionDrainTimeout" xml:"ConnectionDrainTimeout"`
	EstablishedTimeout        int         `json:"EstablishedTimeout" xml:"EstablishedTimeout"`
	HealthCheck               string      `json:"HealthCheck" xml:"HealthCheck"`
	HealthCheckConnectPort    int         `json:"HealthCheckConnectPort" xml:"HealthCheckConnectPort"`
	HealthCheckConnectTimeout int         `json:"HealthCheckConnectTimeout" xml:"HealthCheckConnectTimeout"`
	HealthCheckDomain         string      `json:"HealthCheckDomain" xml:"HealthCheckDomain"`
	HealthCheckHttpCode       string      `json:"HealthCheckHttpCode" xml:"HealthCheckHttpCode"`
	HealthCheckInterval       int         `json:"HealthCheckInterval" xml:"HealthCheckInterval"`
	HealthCheckMethod         string      `json:"HealthCheckMethod" xml:"HealthCheckMethod"`
	HealthCheckType           string      `json:"HealthCheckType" xml:"HealthCheckType"`
	HealthCheckURI            string      `json:"HealthCheckURI" xml:"HealthCheckURI"`
	HealthyThreshold          int         `json:"HealthyThreshold" xml:"HealthyThreshold"`
	MasterSlaveServerGroupId  string      `json:"MasterSlaveServerGroupId" xml:"MasterSlaveServerGroupId"`
	PersistenceTimeout        int         `json:"PersistenceTimeout" xml:"PersistenceTimeout"`
	UnhealthyThreshold        int         `json:"UnhealthyThreshold" xml:"UnhealthyThreshold"`
	HealthCheckSwitch         string      `json:"HealthCheckSwitch" xml:"HealthCheckSwitch"`
	PortRanges                []PortRange `json:"PortRanges" xml:"PortRanges"`
}