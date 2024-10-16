package clusterapi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	capg "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/cmd/openshift-install/command"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/lbconfig"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/gcp"
)

const (
	infrastructureFilepath = "/opt/openshift/manifests/cluster-infrastructure-02-config.yml"

	// replaceable is the string that precedes the encoded data in the ignition data.
	// The data must be replaced before decoding the string, and the string must be
	// prepended to the encoded data.
	replaceable = "data:text/plain;charset=utf-8;base64,"
)

// EditIgnition attempts to edit the contents of the bootstrap ignition when the user has selected
// a custom DNS configuration. Find the public and private load balancer addresses and fill in the
// infrastructure file within the ignition struct.
func EditIgnition(ctx context.Context, in clusterapi.IgnitionInput) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*2)
	defer cancel()

	if in.InstallConfig.Config.GCP.UserProvisionedDNS == gcp.UserProvisionedDNSEnabled {
		gcpCluster := &capg.GCPCluster{}
		key := client.ObjectKey{
			Name:      in.InfraID,
			Namespace: capiutils.Namespace,
		}
		if err := in.Client.Get(ctx, key, gcpCluster); err != nil {
			return nil, fmt.Errorf("failed to get GCP cluster: %w", err)
		}

		svc, err := NewComputeService()
		if err != nil {
			return nil, err
		}

		project := in.InstallConfig.Config.GCP.ProjectID
		if in.InstallConfig.Config.GCP.NetworkProjectID != "" {
			project = in.InstallConfig.Config.GCP.NetworkProjectID
		}

		computeAddress := ""
		if in.InstallConfig.Config.Publish == types.ExternalPublishingStrategy {
			apiIPAddress := *gcpCluster.Status.Network.APIServerAddress
			addressCut := apiIPAddress[strings.LastIndex(apiIPAddress, "/")+1:]
			computeAddressObj, err := svc.GlobalAddresses.Get(project, addressCut).Context(ctx).Do()
			if err != nil {
				return nil, fmt.Errorf("failed to get global compute address: %w", err)
			}
			computeAddress = computeAddressObj.Address
		}

		apiIntIPAddress := *gcpCluster.Status.Network.APIInternalAddress
		addressIntCut := apiIntIPAddress[strings.LastIndex(apiIntIPAddress, "/")+1:]
		computeIntAddress, err := svc.Addresses.Get(project, in.InstallConfig.Config.GCP.Region, addressIntCut).Context(ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("failed to get compute address: %w", err)
		}

		ignData := &igntypes.Config{}
		err = json.Unmarshal(in.BootstrapIgnData, ignData)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal bootstrap ignition: %w", err)
		}

		err = addLoadBalancersToInfra(gcp.Name, ignData, []string{computeAddress}, []string{computeIntAddress.Address})
		if err != nil {
			return nil, fmt.Errorf("failed to add load balancers to ignition config: %w", err)
		}

		lbConfig, err := lbconfig.GenerateLBConfigOverride(computeIntAddress.Address, computeAddress)
		if err != nil {
			return nil, err
		}
		if err := asset.NewDefaultFileWriter(lbConfig).PersistToFile(command.RootOpts.Dir); err != nil {
			return nil, fmt.Errorf("failed to save %s to state file: %w", lbConfig.Name(), err)
		}

		editedIgnBytes, err := json.Marshal(ignData)
		if err != nil {
			return nil, fmt.Errorf("failed to convert ignition data to json: %w", err)
		}

		return editedIgnBytes, nil
	}

	return nil, nil
}

// addLoadBalancersToInfra will load the public and private load balancer information into
// the infrastructure CR. This will occur after the data has already been inserted into the
// ignition file.
func addLoadBalancersToInfra(platform string, config *igntypes.Config, publicLBs []string, privateLBs []string) error {
	for i, fileData := range config.Storage.Files {
		// update the contents of this file
		if fileData.Path == infrastructureFilepath {
			contents := config.Storage.Files[i].Contents.Source
			replaced := strings.Replace(*contents, replaceable, "", 1)

			rawDecodedText, err := base64.StdEncoding.DecodeString(replaced)
			if err != nil {
				return fmt.Errorf("failed to decode contents of ignition file: %w", err)
			}

			infra := &configv1.Infrastructure{}
			if err := yaml.Unmarshal(rawDecodedText, infra); err != nil {
				return fmt.Errorf("failed to unmarshal infrastructure: %w", err)
			}

			// convert the list of strings to a list of IPs
			apiIntLbs := []configv1.IP{}
			for _, ip := range privateLBs {
				apiIntLbs = append(apiIntLbs, configv1.IP(ip))
			}
			apiLbs := []configv1.IP{}
			for _, ip := range publicLBs {
				apiLbs = append(apiLbs, configv1.IP(ip))
			}
			cloudLBInfo := configv1.CloudLoadBalancerIPs{
				APIIntLoadBalancerIPs: apiIntLbs,
				APILoadBalancerIPs:    apiLbs,
			}

			if infra.Status.PlatformStatus.GCP.CloudLoadBalancerConfig.DNSType == configv1.ClusterHostedDNSType {
				infra.Status.PlatformStatus.GCP.CloudLoadBalancerConfig.ClusterHosted = &cloudLBInfo
			}

			// convert the infrastructure back to an encoded string
			infraContents, err := yaml.Marshal(infra)
			if err != nil {
				return fmt.Errorf("failed to marshal infrastructure: %w", err)
			}

			encoded := fmt.Sprintf("%s%s", replaceable, base64.StdEncoding.EncodeToString(infraContents))
			// replace the contents with the edited information
			config.Storage.Files[i].Contents.Source = &encoded

			break
		}
	}

	return nil
}
