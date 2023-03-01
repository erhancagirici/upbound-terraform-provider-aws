// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package ecs

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
			Factory:  DataSourceCluster,
			TypeName: "aws_ecs_cluster",
		},
		{
			Factory:  DataSourceContainerDefinition,
			TypeName: "aws_ecs_container_definition",
		},
		{
			Factory:  DataSourceService,
			TypeName: "aws_ecs_service",
		},
		{
			Factory:  DataSourceTaskDefinition,
			TypeName: "aws_ecs_task_definition",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  ResourceAccountSettingDefault,
			TypeName: "aws_ecs_account_setting_default",
		},
		{
			Factory:  ResourceCapacityProvider,
			TypeName: "aws_ecs_capacity_provider",
		},
		{
			Factory:  ResourceCluster,
			TypeName: "aws_ecs_cluster",
		},
		{
			Factory:  ResourceClusterCapacityProviders,
			TypeName: "aws_ecs_cluster_capacity_providers",
		},
		{
			Factory:  ResourceService,
			TypeName: "aws_ecs_service",
		},
		{
			Factory:  ResourceTag,
			TypeName: "aws_ecs_tag",
		},
		{
			Factory:  ResourceTaskDefinition,
			TypeName: "aws_ecs_task_definition",
		},
		{
			Factory:  ResourceTaskSet,
			TypeName: "aws_ecs_task_set",
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.ECS
}

var ServicePackage = &servicePackage{}
