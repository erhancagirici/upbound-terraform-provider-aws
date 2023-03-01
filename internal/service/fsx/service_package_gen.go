// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package fsx

import (
	"context"

	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  DataSourceOpenzfsSnapshot,
			TypeName: "aws_fsx_openzfs_snapshot",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  ResourceBackup,
			TypeName: "aws_fsx_backup",
		},
		{
			Factory:  ResourceDataRepositoryAssociation,
			TypeName: "aws_fsx_data_repository_association",
		},
		{
			Factory:  ResourceFileCache,
			TypeName: "aws_fsx_file_cache",
		},
		{
			Factory:  ResourceLustreFileSystem,
			TypeName: "aws_fsx_lustre_file_system",
		},
		{
			Factory:  ResourceOntapFileSystem,
			TypeName: "aws_fsx_ontap_file_system",
		},
		{
			Factory:  ResourceOntapStorageVirtualMachine,
			TypeName: "aws_fsx_ontap_storage_virtual_machine",
		},
		{
			Factory:  ResourceOntapVolume,
			TypeName: "aws_fsx_ontap_volume",
		},
		{
			Factory:  ResourceOpenzfsFileSystem,
			TypeName: "aws_fsx_openzfs_file_system",
		},
		{
			Factory:  ResourceOpenzfsSnapshot,
			TypeName: "aws_fsx_openzfs_snapshot",
		},
		{
			Factory:  ResourceOpenzfsVolume,
			TypeName: "aws_fsx_openzfs_volume",
		},
		{
			Factory:  ResourceWindowsFileSystem,
			TypeName: "aws_fsx_windows_file_system",
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.FSx
}

var ServicePackage = &servicePackage{}
