package fwprovider

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	fwtypes "github.com/hashicorp/terraform-provider-aws/internal/framework/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// New returns a new, initialized Terraform Plugin Framework-style provider instance.
// The provider instance is fully configured once the `Configure` method has been called.
func New(primary interface{ Meta() interface{} }) provider.Provider {
	return &fwprovider{
		Primary: primary,
	}
}

type fwprovider struct {
	Primary interface{ Meta() interface{} }
}

// GetSchema returns the schema for this provider's configuration.
func (p *fwprovider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	var diags diag.Diagnostics

	// This schema must match exactly the Terraform Protocol v5 (Terraform Plugin SDK v2) provider's schema.
	schema := tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"access_key": {
				Type:        types.StringType,
				Optional:    true,
				Description: "The access key for API operations. You can retrieve this\nfrom the 'Security & Credentials' section of the AWS console.",
			},
			"allowed_account_ids": {
				Type:     types.SetType{ElemType: types.StringType},
				Optional: true,
			},
			"custom_ca_bundle": {
				Type:        types.StringType,
				Optional:    true,
				Description: "File containing custom root and intermediate certificates. Can also be configured using the `AWS_CA_BUNDLE` environment variable. (Setting `ca_bundle` in the shared config file is not supported.)",
			},
			"ec2_metadata_service_endpoint": {
				Type:        types.StringType,
				Optional:    true,
				Description: "Address of the EC2 metadata service endpoint to use. Can also be configured using the `AWS_EC2_METADATA_SERVICE_ENDPOINT` environment variable.",
			},
			"ec2_metadata_service_endpoint_mode": {
				Type:        types.StringType,
				Optional:    true,
				Description: "Protocol to use with EC2 metadata service endpoint.Valid values are `IPv4` and `IPv6`. Can also be configured using the `AWS_EC2_METADATA_SERVICE_ENDPOINT_MODE` environment variable.",
			},
			"forbidden_account_ids": {
				Type:     types.SetType{ElemType: types.StringType},
				Optional: true,
			},
			"http_proxy": {
				Type:        types.StringType,
				Optional:    true,
				Description: "The address of an HTTP proxy to use when accessing the AWS API. Can also be configured using the `HTTP_PROXY` or `HTTPS_PROXY` environment variables.",
			},
			"insecure": {
				Type:        types.BoolType,
				Optional:    true,
				Description: "Explicitly allow the provider to perform \"insecure\" SSL requests. If omitted, default value is `false`",
			},
			"max_retries": {
				Type:        types.Int64Type,
				Optional:    true,
				Description: "The maximum number of times an AWS API request is\nbeing executed. If the API request still fails, an error is\nthrown.",
			},
			"profile": {
				Type:        types.StringType,
				Optional:    true,
				Description: "The profile for API operations. If not set, the default profile\ncreated with `aws configure` will be used.",
			},
			"region": {
				Type:        types.StringType,
				Optional:    true,
				Description: "The region where AWS operations will take place. Examples\nare us-east-1, us-west-2, etc.", // lintignore:AWSAT003
			},
			"s3_force_path_style": {
				Type:               types.BoolType,
				Optional:           true,
				Description:        "Set this to true to enable the request to use path-style addressing,\ni.e., https://s3.amazonaws.com/BUCKET/KEY. By default, the S3 client will\nuse virtual hosted bucket addressing when possible\n(https://BUCKET.s3.amazonaws.com/KEY). Specific to the Amazon S3 service.",
				DeprecationMessage: "Use s3_use_path_style instead.",
			},
			"s3_use_path_style": {
				Type:        types.BoolType,
				Optional:    true,
				Description: "Set this to true to enable the request to use path-style addressing,\ni.e., https://s3.amazonaws.com/BUCKET/KEY. By default, the S3 client will\nuse virtual hosted bucket addressing when possible\n(https://BUCKET.s3.amazonaws.com/KEY). Specific to the Amazon S3 service.",
			},
			"secret_key": {
				Type:        types.StringType,
				Optional:    true,
				Description: "The secret key for API operations. You can retrieve this\nfrom the 'Security & Credentials' section of the AWS console.",
			},
			"shared_config_files": {
				Type:        types.ListType{ElemType: types.StringType},
				Optional:    true,
				Description: "List of paths to shared config files. If not set, defaults to [~/.aws/config].",
			},
			"shared_credentials_file": {
				Type:               types.StringType,
				Optional:           true,
				Description:        "The path to the shared credentials file. If not set, defaults to ~/.aws/credentials.",
				DeprecationMessage: "Use shared_credentials_files instead.",
			},
			"shared_credentials_files": {
				Type:        types.ListType{ElemType: types.StringType},
				Optional:    true,
				Description: "List of paths to shared credentials files. If not set, defaults to [~/.aws/credentials].",
			},
			"skip_credentials_validation": {
				Type:        types.BoolType,
				Optional:    true,
				Description: "Skip the credentials validation via STS API. Used for AWS API implementations that do not have STS available/implemented.",
			},
			"skip_get_ec2_platforms": {
				Type:               types.BoolType,
				Optional:           true,
				Description:        "Skip getting the supported EC2 platforms. Used by users that don't have ec2:DescribeAccountAttributes permissions.",
				DeprecationMessage: `With the retirement of EC2-Classic the skip_get_ec2_platforms attribute has been deprecated and will be removed in a future version.`,
			},
			"skip_metadata_api_check": {
				Type:        types.StringType,
				Optional:    true,
				Description: "Skip the AWS Metadata API check. Used for AWS API implementations that do not have a metadata api endpoint.",
			},
			"skip_region_validation": {
				Type:        types.BoolType,
				Optional:    true,
				Description: "Skip static validation of region name. Used by users of alternative AWS-like APIs or users w/ access to regions that are not public (yet).",
			},
			"skip_requesting_account_id": {
				Type:        types.BoolType,
				Optional:    true,
				Description: "Skip requesting the account ID. Used for AWS API implementations that do not have IAM/STS API and/or metadata API.",
			},
			"sts_region": {
				Type:        types.StringType,
				Optional:    true,
				Description: "The region where AWS STS operations will take place. Examples\nare us-east-1 and us-west-2.", // lintignore:AWSAT003
			},
			"token": {
				Type:        types.StringType,
				Optional:    true,
				Description: "session token. A session token is only required if you are\nusing temporary security credentials.",
			},
			"use_dualstack_endpoint": {
				Type:        types.BoolType,
				Optional:    true,
				Description: "Resolve an endpoint with DualStack capability",
			},
			"use_fips_endpoint": {
				Type:        types.BoolType,
				Optional:    true,
				Description: "Resolve an endpoint with FIPS capability",
			},
		},
		Blocks: map[string]tfsdk.Block{
			"assume_role": {
				Attributes: map[string]tfsdk.Attribute{
					"duration": {
						Type:        fwtypes.DurationType,
						Optional:    true,
						Description: "The duration, between 15 minutes and 12 hours, of the role session. Valid time units are ns, us (or µs), ms, s, h, or m.",
					},
					"duration_seconds": {
						Type:               types.Int64Type,
						Optional:           true,
						Description:        "The duration, in seconds, of the role session.",
						DeprecationMessage: "Use assume_role.duration instead",
					},
					"external_id": {
						Type:        types.StringType,
						Optional:    true,
						Description: "A unique identifier that might be required when you assume a role in another account.",
					},
					"policy": {
						Type:        types.StringType,
						Optional:    true,
						Description: "IAM Policy JSON describing further restricting permissions for the IAM Role being assumed.",
					},
					"policy_arns": {
						Type:        types.SetType{ElemType: types.StringType},
						Optional:    true,
						Description: "Amazon Resource Names (ARNs) of IAM Policies describing further restricting permissions for the IAM Role being assumed.",
					},
					"role_arn": {
						Type:        types.StringType,
						Optional:    true,
						Description: "Amazon Resource Name (ARN) of an IAM Role to assume prior to making API calls.",
					},
					"session_name": {
						Type:        types.StringType,
						Optional:    true,
						Description: "An identifier for the assumed role session.",
					},
					"source_identity": {
						Type:        types.StringType,
						Optional:    true,
						Description: "Source identity specified by the principal assuming the role.",
					},
					"tags": {
						Type:        types.MapType{ElemType: types.StringType},
						Optional:    true,
						Description: "Assume role session tags.",
					},
					"transitive_tag_keys": {
						Type:        types.SetType{ElemType: types.StringType},
						Optional:    true,
						Description: "Assume role session tag keys to pass to any subsequent sessions.",
					},
				},
				NestingMode: tfsdk.BlockNestingModeList,
				MaxItems:    1,
			},
			"assume_role_with_web_identity": {
				Attributes: map[string]tfsdk.Attribute{
					"duration": {
						Type:        fwtypes.DurationType,
						Optional:    true,
						Description: "The duration, between 15 minutes and 12 hours, of the role session. Valid time units are ns, us (or µs), ms, s, h, or m.",
					},
					"policy": {
						Type:        types.StringType,
						Optional:    true,
						Description: "IAM Policy JSON describing further restricting permissions for the IAM Role being assumed.",
					},
					"policy_arns": {
						Type:        types.SetType{ElemType: types.StringType},
						Optional:    true,
						Description: "Amazon Resource Names (ARNs) of IAM Policies describing further restricting permissions for the IAM Role being assumed.",
					},
					"role_arn": {
						Type:        types.StringType,
						Optional:    true,
						Description: "Amazon Resource Name (ARN) of an IAM Role to assume prior to making API calls.",
					},
					"session_name": {
						Type:        types.StringType,
						Optional:    true,
						Description: "An identifier for the assumed role session.",
					},
					"web_identity_token": {
						Type:     types.StringType,
						Optional: true,
					},
					"web_identity_token_file": {
						Type:     types.StringType,
						Optional: true,
					},
				},
				NestingMode: tfsdk.BlockNestingModeList,
				MaxItems:    1,
			},
			"default_tags": {
				Attributes: map[string]tfsdk.Attribute{
					"tags": {
						Type:        types.MapType{ElemType: types.StringType},
						Optional:    true,
						Description: "Resource tags to default across all resources",
					},
				},
				NestingMode: tfsdk.BlockNestingModeList,
				MaxItems:    1,
				Description: "Configuration block with settings to default resource tags across all resources.",
			},
			"endpoints": endpointsBlock(),
			"ignore_tags": {
				Attributes: map[string]tfsdk.Attribute{
					"key_prefixes": {
						Type:        types.SetType{ElemType: types.StringType},
						Optional:    true,
						Description: "Resource tag key prefixes to ignore across all resources.",
					},
					"keys": {
						Type:        types.SetType{ElemType: types.StringType},
						Optional:    true,
						Description: "Resource tag keys to ignore across all resources.",
					},
				},
				NestingMode: tfsdk.BlockNestingModeList,
				MaxItems:    1,
				Description: "Configuration block with settings to ignore resource tags across all resources.",
			},
		},
	}

	return schema, diags
}

// Configure is called at the beginning of the provider lifecycle, when
// Terraform sends to the provider the values the user specified in the
// provider configuration block.
func (p *fwprovider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	// Provider's parsed configuration (its instance state) is available through the primary provider's Meta() method.
	v := p.Primary.Meta()
	response.DataSourceData = v
	response.ResourceData = v
}

// DataSources returns a slice of functions to instantiate each DataSource
// implementation.
//
// The data source type name is determined by the DataSource implementing
// the Metadata method. All data sources must have unique names.
func (p *fwprovider) DataSources(ctx context.Context) []func() datasource.DataSource {
	var dataSources []func() datasource.DataSource

	for _, sp := range p.Primary.Meta().(*conns.AWSClient).ServicePackages {
		for _, v := range sp.FrameworkDataSources(ctx) {
			v, err := v(ctx)

			if err != nil {
				tflog.Warn(ctx, "creating data source", map[string]interface{}{
					"service_package_name": sp.ServicePackageName(),
					"error":                err.Error(),
				})

				continue
			}

			dataSources = append(dataSources, func() datasource.DataSource {
				return newWrappedDataSource(v)
			})
		}
	}

	return dataSources
}

// Resources returns a slice of functions to instantiate each Resource
// implementation.
//
// The resource type name is determined by the Resource implementing
// the Metadata method. All resources must have unique names.
func (p *fwprovider) Resources(ctx context.Context) []func() resource.Resource {
	var resources []func() resource.Resource

	for _, sp := range p.Primary.Meta().(*conns.AWSClient).ServicePackages {
		for _, v := range sp.FrameworkResources(ctx) {
			v, err := v(ctx)

			if err != nil {
				tflog.Warn(ctx, "creating resource", map[string]interface{}{
					"service_package_name": sp.ServicePackageName(),
					"error":                err.Error(),
				})

				continue
			}

			resources = append(resources, func() resource.Resource {
				return newWrappedResource(v)
			})
		}
	}

	return resources
}

func endpointsBlock() tfsdk.Block {
	endpointsAttributes := make(map[string]tfsdk.Attribute)

	for _, serviceKey := range names.Aliases() {
		endpointsAttributes[serviceKey] = tfsdk.Attribute{
			Type:        types.StringType,
			Optional:    true,
			Description: "Use this to override the default service endpoint URL",
		}
	}

	return tfsdk.Block{
		Attributes:  endpointsAttributes,
		NestingMode: tfsdk.BlockNestingModeSet,
	}
}

// wrappedDataSource wraps a data source, adding common functionality.
type wrappedDataSource struct {
	inner    datasource.DataSourceWithConfigure
	typeName string
}

func newWrappedDataSource(inner datasource.DataSourceWithConfigure) datasource.DataSourceWithConfigure {
	return &wrappedDataSource{inner: inner, typeName: strings.TrimPrefix(reflect.TypeOf(inner).String(), "*")}
}

func (w *wrappedDataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	w.inner.Metadata(ctx, request, response)
}

func (w *wrappedDataSource) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	if v, ok := w.inner.(datasource.DataSourceWithSchema); ok {
		v.Schema(ctx, request, response)
		return
	}
	response.Diagnostics.AddError(
		"DataSource Schema Not Implemented",
		"This data source does not support get schema. Please contact the provider developer for additional information.",
	)
}

func (w *wrappedDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	tflog.Debug(ctx, fmt.Sprintf("%s.Read enter", w.typeName))

	w.inner.Read(ctx, request, response)

	tflog.Debug(ctx, fmt.Sprintf("%s.Read exit", w.typeName))
}

func (w *wrappedDataSource) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	w.inner.Configure(ctx, request, response)
}

// wrappedResource wraps a resource, adding common functionality.
type wrappedResource struct {
	inner    resource.ResourceWithConfigure
	typeName string
}

func newWrappedResource(inner resource.ResourceWithConfigure) resource.ResourceWithConfigure {
	return &wrappedResource{inner: inner, typeName: strings.TrimPrefix(reflect.TypeOf(inner).String(), "*")}
}

func (w *wrappedResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	w.inner.Metadata(ctx, request, response)
}

func (w *wrappedResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	if v, ok := w.inner.(resource.ResourceWithSchema); ok {
		v.Schema(ctx, request, response)
		return
	}
	response.Diagnostics.AddError(
		"Resource Schema Not Implemented",
		"This resource does not support get schema. Please contact the provider developer for additional information.",
	)
}

func (w *wrappedResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	tflog.Debug(ctx, fmt.Sprintf("%s.Create enter", w.typeName))

	w.inner.Create(ctx, request, response)

	tflog.Debug(ctx, fmt.Sprintf("%s.Create exit", w.typeName))
}

func (w *wrappedResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	tflog.Debug(ctx, fmt.Sprintf("%s.Read enter", w.typeName))

	w.inner.Read(ctx, request, response)

	tflog.Debug(ctx, fmt.Sprintf("%s.Read exit", w.typeName))
}

func (w *wrappedResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	tflog.Debug(ctx, fmt.Sprintf("%s.Update enter", w.typeName))

	w.inner.Update(ctx, request, response)

	tflog.Debug(ctx, fmt.Sprintf("%s.Update exit", w.typeName))
}

func (w *wrappedResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	tflog.Debug(ctx, fmt.Sprintf("%s.Delete enter", w.typeName))

	w.inner.Delete(ctx, request, response)

	tflog.Debug(ctx, fmt.Sprintf("%s.Delete exit", w.typeName))
}

func (w *wrappedResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	w.inner.Configure(ctx, request, response)
}

func (w *wrappedResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	if v, ok := w.inner.(resource.ResourceWithImportState); ok {
		v.ImportState(ctx, request, response)

		return
	}

	response.Diagnostics.AddError(
		"Resource Import Not Implemented",
		"This resource does not support import. Please contact the provider developer for additional information.",
	)
}

func (w *wrappedResource) ModifyPlan(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse) {
	if v, ok := w.inner.(resource.ResourceWithModifyPlan); ok {
		v.ModifyPlan(ctx, request, response)

		return
	}
}

func (w *wrappedResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	if v, ok := w.inner.(resource.ResourceWithConfigValidators); ok {
		return v.ConfigValidators(ctx)
	}

	return nil
}

func (w *wrappedResource) ValidateConfig(ctx context.Context, request resource.ValidateConfigRequest, response *resource.ValidateConfigResponse) {
	if v, ok := w.inner.(resource.ResourceWithValidateConfig); ok {
		v.ValidateConfig(ctx, request, response)
	}
}
