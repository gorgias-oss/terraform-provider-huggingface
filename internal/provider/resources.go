package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorgias-oss/huggingface-endpoints-client-go"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &endpointResource{}
	_ resource.ResourceWithConfigure = &endpointResource{}
)

func NewEndpointResource() resource.Resource {
	return &endpointResource{}
}

type endpointResource struct {
	client *huggingface.Client
}

type endpointResourceModel struct {
	AccountId types.String `tfsdk:"account_id"`
	Compute   Compute      `tfsdk:"compute"`
	Model     Model        `tfsdk:"model"`
	Name      types.String `tfsdk:"name"`
	Cloud     Cloud        `tfsdk:"cloud"`
	Type      types.String `tfsdk:"type"`
}

func (r *endpointResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*huggingface.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"unexpected data source configure type",
			fmt.Sprintf("expected *huggingface.Client, got: %T.", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *endpointResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_endpoint"
}

func (r *endpointResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Hugging Face Inference Endpoint.",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Optional: true,
			},
			"compute": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"accelerator": schema.StringAttribute{
						Required: true,
					},
					"instance_size": schema.StringAttribute{
						Required: true,
					},
					"instance_type": schema.StringAttribute{
						Required: true,
					},
					"scaling": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"max_replica": schema.Int64Attribute{
								Required: true,
							},
							"min_replica": schema.Int64Attribute{
								Required: true,
							},
							"scale_to_zero_timeout": schema.Int64Attribute{
								Optional: true,
								Computed: true,
							},
							"measure": schema.SingleNestedAttribute{
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"hardware_usage": schema.Float64Attribute{
										Optional: true,
									},
									"pending_requests": schema.Float64Attribute{
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"model": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"framework": schema.StringAttribute{
						Required: true,
					},
					"env": schema.MapAttribute{
						Optional:    true,
						ElementType: types.StringType,
					},
					"image": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"huggingface": schema.SingleNestedAttribute{
								Optional: true,
							},
							"custom": schema.SingleNestedAttribute{
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"credentials": schema.SingleNestedAttribute{
										Optional: true,
										Attributes: map[string]schema.Attribute{
											"username": schema.StringAttribute{
												Required: true,
											},
											"password": schema.StringAttribute{
												Required: true,
											},
										},
									},
									"health_route": schema.StringAttribute{
										Optional: true,
									},
									"port": schema.Int64Attribute{
										Optional: true,
										Computed: true,
									},
									"url": schema.StringAttribute{
										Required: true,
									},
								},
							},
							"tei": schema.SingleNestedAttribute{
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"health_route": schema.StringAttribute{
										Optional: true,
									},
									"port": schema.Int64Attribute{
										Optional: true,
										Computed: true,
									},
									"url": schema.StringAttribute{
										Required: true,
									},
									"max_batch_tokens": schema.Int64Attribute{
										Optional: true,
									},
									"max_concurrent_requests": schema.Int64Attribute{
										Optional: true,
									},
									"pooling": schema.StringAttribute{
										Optional: true,
									},
								},
							},
							"tgi": schema.SingleNestedAttribute{
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"health_route": schema.StringAttribute{
										Optional: true,
									},
									"port": schema.Int64Attribute{
										Optional: true,
										Computed: true,
									},
									"url": schema.StringAttribute{
										Required: true,
									},
									"max_batch_prefill_tokens": schema.Int64Attribute{
										Optional: true,
									},
									"max_batch_total_tokens": schema.Int64Attribute{
										Optional: true,
									},
									"max_input_length": schema.Int64Attribute{
										Optional: true,
									},
									"max_total_tokens": schema.Int64Attribute{
										Optional: true,
									},
									"diable_custom_kernels": schema.BoolAttribute{
										Optional: true,
									},
									"quantize": schema.StringAttribute{
										Optional: true,
									},
								},
							},
							"tgi_neuron": schema.SingleNestedAttribute{
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"health_route": schema.StringAttribute{
										Optional: true,
									},
									"port": schema.Int64Attribute{
										Optional: true,
										Computed: true,
									},
									"url": schema.StringAttribute{
										Required: true,
									},
									"max_batch_prefill_tokens": schema.Int64Attribute{
										Optional: true,
									},
									"max_batch_total_tokens": schema.Int64Attribute{
										Optional: true,
									},
									"max_input_length": schema.Int64Attribute{
										Optional: true,
									},
									"max_total_tokens": schema.Int64Attribute{
										Optional: true,
									},
									"hf_auto_cast_type": schema.StringAttribute{
										Optional: true,
									},
									"hf_num_cores": schema.Int64Attribute{
										Optional: true,
									},
								},
							},
							"llamacpp": schema.SingleNestedAttribute{
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"health_route": schema.StringAttribute{
										Optional: true,
									},
									"port": schema.Int64Attribute{
										Optional: true,
										Computed: true,
									},
									"url": schema.StringAttribute{
										Required: true,
									},
									"ctx_size": schema.Int64Attribute{
										Optional: true,
									},
									"embeddings": schema.BoolAttribute{
										Optional: true,
									},
									"model_path": schema.StringAttribute{
										Required: true,
									},
									"n_parallel": schema.Int64Attribute{
										Optional: true,
									},
									"threads_http": schema.Int64Attribute{
										Optional: true,
									},
								},
							},
							"vllm": schema.SingleNestedAttribute{
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"health_route": schema.StringAttribute{
										Optional: true,
									},
									"port": schema.Int64Attribute{
										Optional: true,
										Computed: true,
									},
									"url": schema.StringAttribute{
										Required: true,
									},
									"kv_cache_dtype": schema.StringAttribute{
										Optional: true,
									},
									"max_num_batched_tokens": schema.Int64Attribute{
										Optional: true,
									},
									"max_num_seqs": schema.Int64Attribute{
										Optional: true,
									},
									"tensor_parallel_size": schema.Int64Attribute{
										Optional: true,
									},
									"server_args": schema.ListAttribute{
										Optional:    true,
										ElementType: types.StringType,
									},
								},
							},
						},
					},
					"repository": schema.StringAttribute{
						Required: true,
					},
					"revision": schema.StringAttribute{
						Computed: true,
						Optional: true,
					},
					"task": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"cloud": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"region": schema.StringAttribute{
						Required: true,
					},
					"vendor": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"type": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func clientEndpointToProviderEndpoint(endpoint huggingface.EndpointDetails) endpointResourceModel {
	var image Image
	if endpoint.Model.Image.Huggingface != nil {
		image = Image{
			Huggingface: &Huggingface{},
		}
	} else if endpoint.Model.Image.Custom != nil {
		var port64 *int64
		port := endpoint.Model.Image.Custom.Port
		if port == nil {
			port64 = nil
		} else {
			portInt64 := int64(*port)
			port64 = &portInt64
		}
		image = Image{
			Custom: &Custom{
				HealthRoute: endpoint.Model.Image.Custom.HealthRoute,
				Port:        types.Int64PointerValue(port64),
				URL:         endpoint.Model.Image.Custom.URL,
			},
		}
		if endpoint.Model.Image.Custom.Credentials != nil {
			image.Custom.Credentials = &Credentials{
				Username: endpoint.Model.Image.Custom.Credentials.Username,
				Password: endpoint.Model.Image.Custom.Credentials.Password,
			}
		}
	} else if endpoint.Model.Image.Tgi != nil {
		var port64 *int64
		port := endpoint.Model.Image.Tgi.Port
		if port == nil {
			port64 = nil
		} else {
			portInt64 := int64(*port)
			port64 = &portInt64
		}
		image = Image{
			Tgi: &Tgi{
				HealthRoute:           endpoint.Model.Image.Tgi.HealthRoute,
				Port:                  types.Int64PointerValue(port64),
				URL:                   endpoint.Model.Image.Tgi.URL,
				MaxBatchPrefillTokens: endpoint.Model.Image.Tgi.MaxBatchPrefillTokens,
				MaxBatchTotalTokens:   endpoint.Model.Image.Tgi.MaxBatchTotalTokens,
				MaxInputLength:        endpoint.Model.Image.Tgi.MaxInputLength,
				MaxTotalTokens:        endpoint.Model.Image.Tgi.MaxTotalTokens,
				DisableCustomKernels:  endpoint.Model.Image.Tgi.DisableCustomKernels,
				Quantize:              endpoint.Model.Image.Tgi.Quantize,
			},
		}
	} else if endpoint.Model.Image.TgiNeuron != nil {
		var port64 *int64
		port := endpoint.Model.Image.TgiNeuron.Port
		if port == nil {
			port64 = nil
		} else {
			portInt64 := int64(*port)
			port64 = &portInt64
		}
		image = Image{
			TgiNeuron: &TgiNeuron{
				HealthRoute:           endpoint.Model.Image.TgiNeuron.HealthRoute,
				Port:                  types.Int64PointerValue(port64),
				URL:                   endpoint.Model.Image.TgiNeuron.URL,
				MaxBatchPrefillTokens: endpoint.Model.Image.TgiNeuron.MaxBatchPrefillTokens,
				MaxBatchTotalTokens:   endpoint.Model.Image.TgiNeuron.MaxBatchTotalTokens,
				MaxInputLength:        endpoint.Model.Image.TgiNeuron.MaxInputLength,
				MaxTotalTokens:        endpoint.Model.Image.TgiNeuron.MaxTotalTokens,
				HfAutoCastType:        endpoint.Model.Image.TgiNeuron.HfAutoCastType,
				HfNumCores:            endpoint.Model.Image.TgiNeuron.HfNumCores,
			},
		}
	} else if endpoint.Model.Image.Tei != nil {
		var port64 *int64
		port := endpoint.Model.Image.Tei.Port
		if port == nil {
			port64 = nil
		} else {
			portInt64 := int64(*port)
			port64 = &portInt64
		}
		image = Image{
			Tei: &Tei{
				HealthRoute:           endpoint.Model.Image.Tei.HealthRoute,
				Port:                  types.Int64PointerValue(port64),
				URL:                   endpoint.Model.Image.Tei.URL,
				MaxBatchTokens:        endpoint.Model.Image.Tei.MaxBatchTokens,
				MaxConcurrentRequests: endpoint.Model.Image.Tei.MaxConcurrentRequests,
				Pooling:               endpoint.Model.Image.Tei.Pooling,
			},
		}
	} else if endpoint.Model.Image.Vllm != nil {
		var port64 *int64
		port := endpoint.Model.Image.Vllm.Port
		if port == nil {
			port64 = nil
		} else {
			portInt64 := int64(*port)
			port64 = &portInt64
		}
		image = Image{
			Vllm: &Vllm{
				HealthRoute:         endpoint.Model.Image.Vllm.HealthRoute,
				Port:                types.Int64PointerValue(port64),
				URL:                 endpoint.Model.Image.Vllm.URL,
				KvCacheDtype:        endpoint.Model.Image.Vllm.KvCacheDtype,
				MaxNumBatchedTokens: endpoint.Model.Image.Vllm.MaxNumBatchedTokens,
				MaxNumSeqs:          endpoint.Model.Image.Vllm.MaxNumSeqs,
				TensorParallelSize:  endpoint.Model.Image.Vllm.TensorParallelSize,
				ServerArgs:          endpoint.Model.Image.Vllm.ServerArgs,
			},
		}
	}

	var timeout64 *int64
	timeout := endpoint.Compute.Scaling.ScaleToZeroTimeout
	if timeout == nil {
		timeout64 = nil
	} else {
		timeoutInt := int64(*timeout)
		timeout64 = &timeoutInt
	}

	var measure *Measure = nil
	if endpoint.Compute.Scaling.Measure != nil {
		measure = &Measure{
			HardwareUsage:   endpoint.Compute.Scaling.Measure.HardwareUsage,
			PendingRequests: endpoint.Compute.Scaling.Measure.PendingRequests,
		}
	}

	providerEndpoint := endpointResourceModel{
		AccountId: types.StringPointerValue(endpoint.AccountId),
		Compute: Compute{
			Accelerator:  endpoint.Compute.Accelerator,
			InstanceSize: endpoint.Compute.InstanceSize,
			InstanceType: endpoint.Compute.InstanceType,
			Scaling: Scaling{
				MaxReplica:         endpoint.Compute.Scaling.MaxReplica,
				MinReplica:         endpoint.Compute.Scaling.MinReplica,
				ScaleToZeroTimeout: types.Int64PointerValue(timeout64),
				Measure:            measure,
			},
		},
		Model: Model{
			Framework:  endpoint.Model.Framework,
			Image:      image,
			Repository: endpoint.Model.Repository,
			Revision:   types.StringPointerValue(endpoint.Model.Revision),
			Task:       types.StringPointerValue(endpoint.Model.Task),
			Env:        endpoint.Model.Env,
		},
		Name: types.StringValue(endpoint.Name),
		Cloud: Cloud{
			Region: endpoint.Provider.Region,
			Vendor: endpoint.Provider.Vendor,
		},
		Type: types.StringValue(endpoint.Type),
	}

	if endpoint.Model.Env == nil {
		providerEndpoint.Model.Env = make(map[string]string)
	}

	return providerEndpoint
}

func providerEndpointToCreateEndpointRequest(endpoint endpointResourceModel) huggingface.CreateEndpointRequest {
	var image huggingface.Image
	if endpoint.Model.Image.Huggingface != nil {
		image = huggingface.Image{
			Huggingface: &huggingface.Huggingface{},
		}
	} else if endpoint.Model.Image.Custom != nil {
		var port *int
		if endpoint.Model.Image.Custom.Port.IsUnknown() || endpoint.Model.Image.Custom.Port.IsNull() {
			port = nil
		} else {
			portInt := int(endpoint.Model.Image.Custom.Port.ValueInt64())
			port = &portInt
		}
		image = huggingface.Image{
			Custom: &huggingface.Custom{
				HealthRoute: endpoint.Model.Image.Custom.HealthRoute,
				Port:        port,
				URL:         endpoint.Model.Image.Custom.URL,
			},
		}
		if endpoint.Model.Image.Custom.Credentials != nil {
			image.Custom.Credentials = &huggingface.Credentials{
				Username: endpoint.Model.Image.Custom.Credentials.Username,
				Password: endpoint.Model.Image.Custom.Credentials.Password,
			}
		}
	} else if endpoint.Model.Image.Tgi != nil {
		var port *int
		if endpoint.Model.Image.Tgi.Port.IsUnknown() || endpoint.Model.Image.Tgi.Port.IsNull() {
			port = nil
		} else {
			portInt := int(endpoint.Model.Image.Tgi.Port.ValueInt64())
			port = &portInt
		}
		image = huggingface.Image{
			Tgi: &huggingface.Tgi{
				HealthRoute:           endpoint.Model.Image.Tgi.HealthRoute,
				Port:                  port,
				URL:                   endpoint.Model.Image.Tgi.URL,
				MaxBatchPrefillTokens: endpoint.Model.Image.Tgi.MaxBatchPrefillTokens,
				MaxBatchTotalTokens:   endpoint.Model.Image.Tgi.MaxBatchTotalTokens,
				MaxInputLength:        endpoint.Model.Image.Tgi.MaxInputLength,
				MaxTotalTokens:        endpoint.Model.Image.Tgi.MaxTotalTokens,
				DisableCustomKernels:  endpoint.Model.Image.Tgi.DisableCustomKernels,
				Quantize:              endpoint.Model.Image.Tgi.Quantize,
			},
		}
	} else if endpoint.Model.Image.TgiNeuron != nil {
		var port *int
		if endpoint.Model.Image.TgiNeuron.Port.IsUnknown() || endpoint.Model.Image.TgiNeuron.Port.IsNull() {
			port = nil
		} else {
			portInt := int(endpoint.Model.Image.TgiNeuron.Port.ValueInt64())
			port = &portInt
		}
		image = huggingface.Image{
			TgiNeuron: &huggingface.TgiNeuron{
				HealthRoute:           endpoint.Model.Image.TgiNeuron.HealthRoute,
				Port:                  port,
				URL:                   endpoint.Model.Image.TgiNeuron.URL,
				MaxBatchPrefillTokens: endpoint.Model.Image.TgiNeuron.MaxBatchPrefillTokens,
				MaxBatchTotalTokens:   endpoint.Model.Image.TgiNeuron.MaxBatchTotalTokens,
				MaxInputLength:        endpoint.Model.Image.TgiNeuron.MaxInputLength,
				MaxTotalTokens:        endpoint.Model.Image.TgiNeuron.MaxTotalTokens,
				HfAutoCastType:        endpoint.Model.Image.TgiNeuron.HfAutoCastType,
				HfNumCores:            endpoint.Model.Image.TgiNeuron.HfNumCores,
			},
		}
	} else if endpoint.Model.Image.Tei != nil {
		var port *int
		if endpoint.Model.Image.Tei.Port.IsUnknown() || endpoint.Model.Image.Tei.Port.IsNull() {
			port = nil
		} else {
			portInt := int(endpoint.Model.Image.Tei.Port.ValueInt64())
			port = &portInt
		}
		image = huggingface.Image{
			Tei: &huggingface.Tei{
				HealthRoute:           endpoint.Model.Image.Tei.HealthRoute,
				Port:                  port,
				URL:                   endpoint.Model.Image.Tei.URL,
				MaxBatchTokens:        endpoint.Model.Image.Tei.MaxBatchTokens,
				MaxConcurrentRequests: endpoint.Model.Image.Tei.MaxConcurrentRequests,
				Pooling:               endpoint.Model.Image.Tei.Pooling,
			},
		}
	} else if endpoint.Model.Image.Vllm != nil {
		var port *int
		if endpoint.Model.Image.Vllm.Port.IsUnknown() || endpoint.Model.Image.Vllm.Port.IsNull() {
			port = nil
		} else {
			portInt := int(endpoint.Model.Image.Vllm.Port.ValueInt64())
			port = &portInt
		}
		image = huggingface.Image{
			Vllm: &huggingface.Vllm{
				HealthRoute:         endpoint.Model.Image.Vllm.HealthRoute,
				Port:                port,
				URL:                 endpoint.Model.Image.Vllm.URL,
				KvCacheDtype:        endpoint.Model.Image.Vllm.KvCacheDtype,
				MaxNumBatchedTokens: endpoint.Model.Image.Vllm.MaxNumBatchedTokens,
				MaxNumSeqs:          endpoint.Model.Image.Vllm.MaxNumSeqs,
				TensorParallelSize:  endpoint.Model.Image.Vllm.TensorParallelSize,
				ServerArgs:          endpoint.Model.Image.Vllm.ServerArgs,
			},
		}
	}

	var timeout *int
	if endpoint.Compute.Scaling.ScaleToZeroTimeout.IsUnknown() || endpoint.Compute.Scaling.ScaleToZeroTimeout.IsNull() {
		timeout = nil
	} else {
		timeoutInt := int(endpoint.Compute.Scaling.ScaleToZeroTimeout.ValueInt64())
		timeout = &timeoutInt
	}

	var measure *huggingface.Measure = nil
	if endpoint.Compute.Scaling.Measure != nil {
		measure = &huggingface.Measure{
			HardwareUsage:   endpoint.Compute.Scaling.Measure.HardwareUsage,
			PendingRequests: endpoint.Compute.Scaling.Measure.PendingRequests,
		}
	}

	huggingfaceEndpoint := huggingface.CreateEndpointRequest{
		Name:      endpoint.Name.ValueString(),
		AccountId: endpoint.AccountId.ValueStringPointer(),
		Compute: huggingface.Compute{
			Accelerator:  endpoint.Compute.Accelerator,
			InstanceSize: endpoint.Compute.InstanceSize,
			InstanceType: endpoint.Compute.InstanceType,
			Scaling: huggingface.Scaling{
				MaxReplica:         endpoint.Compute.Scaling.MaxReplica,
				MinReplica:         endpoint.Compute.Scaling.MinReplica,
				ScaleToZeroTimeout: timeout,
				Measure:            measure,
			},
		},
		Model: huggingface.Model{
			Framework:  endpoint.Model.Framework,
			Image:      image,
			Repository: endpoint.Model.Repository,
			Revision:   endpoint.Model.Revision.ValueStringPointer(),
			Task:       endpoint.Model.Task.ValueStringPointer(),
			Env:        endpoint.Model.Env,
		},
		Provider: huggingface.Provider{
			Region: endpoint.Cloud.Region,
			Vendor: endpoint.Cloud.Vendor,
		},
		Type: endpoint.Type.ValueString(),
	}

	return huggingfaceEndpoint
}

func providerEndpointToUpdateEndpointRequest(endpoint endpointResourceModel) huggingface.UpdateEndpointRequest {
	var image huggingface.Image
	if endpoint.Model.Image.Huggingface != nil {
		image = huggingface.Image{
			Huggingface: &huggingface.Huggingface{},
		}
	} else if endpoint.Model.Image.Custom != nil {
		var port *int
		if endpoint.Model.Image.Custom.Port.IsUnknown() || endpoint.Model.Image.Custom.Port.IsNull() {
			port = nil
		} else {
			portInt := int(endpoint.Model.Image.Custom.Port.ValueInt64())
			port = &portInt
		}
		image = huggingface.Image{
			Custom: &huggingface.Custom{
				HealthRoute: endpoint.Model.Image.Custom.HealthRoute,
				Port:        port,
				URL:         endpoint.Model.Image.Custom.URL,
			},
		}
		if endpoint.Model.Image.Custom.Credentials != nil {
			image.Custom.Credentials = &huggingface.Credentials{
				Username: endpoint.Model.Image.Custom.Credentials.Username,
				Password: endpoint.Model.Image.Custom.Credentials.Password,
			}
		}
	} else if endpoint.Model.Image.Tgi != nil {
		var port *int
		if endpoint.Model.Image.Tgi.Port.IsUnknown() || endpoint.Model.Image.Tgi.Port.IsNull() {
			port = nil
		} else {
			portInt := int(endpoint.Model.Image.Tgi.Port.ValueInt64())
			port = &portInt
		}
		image = huggingface.Image{
			Tgi: &huggingface.Tgi{
				HealthRoute:           endpoint.Model.Image.Tgi.HealthRoute,
				Port:                  port,
				URL:                   endpoint.Model.Image.Tgi.URL,
				MaxBatchPrefillTokens: endpoint.Model.Image.Tgi.MaxBatchPrefillTokens,
				MaxBatchTotalTokens:   endpoint.Model.Image.Tgi.MaxBatchTotalTokens,
				MaxInputLength:        endpoint.Model.Image.Tgi.MaxInputLength,
				MaxTotalTokens:        endpoint.Model.Image.Tgi.MaxTotalTokens,
				DisableCustomKernels:  endpoint.Model.Image.Tgi.DisableCustomKernels,
				Quantize:              endpoint.Model.Image.Tgi.Quantize,
			},
		}
	} else if endpoint.Model.Image.TgiNeuron != nil {
		var port *int
		if endpoint.Model.Image.TgiNeuron.Port.IsUnknown() || endpoint.Model.Image.TgiNeuron.Port.IsNull() {
			port = nil
		} else {
			portInt := int(endpoint.Model.Image.TgiNeuron.Port.ValueInt64())
			port = &portInt
		}
		image = huggingface.Image{
			TgiNeuron: &huggingface.TgiNeuron{
				HealthRoute:           endpoint.Model.Image.TgiNeuron.HealthRoute,
				Port:                  port,
				URL:                   endpoint.Model.Image.TgiNeuron.URL,
				MaxBatchPrefillTokens: endpoint.Model.Image.TgiNeuron.MaxBatchPrefillTokens,
				MaxBatchTotalTokens:   endpoint.Model.Image.TgiNeuron.MaxBatchTotalTokens,
				MaxInputLength:        endpoint.Model.Image.TgiNeuron.MaxInputLength,
				MaxTotalTokens:        endpoint.Model.Image.TgiNeuron.MaxTotalTokens,
				HfAutoCastType:        endpoint.Model.Image.TgiNeuron.HfAutoCastType,
				HfNumCores:            endpoint.Model.Image.TgiNeuron.HfNumCores,
			},
		}
	} else if endpoint.Model.Image.Tei != nil {
		var port *int
		if endpoint.Model.Image.Tei.Port.IsUnknown() || endpoint.Model.Image.Tei.Port.IsNull() {
			port = nil
		} else {
			portInt := int(endpoint.Model.Image.Tei.Port.ValueInt64())
			port = &portInt
		}
		image = huggingface.Image{
			Tei: &huggingface.Tei{
				HealthRoute:           endpoint.Model.Image.Tei.HealthRoute,
				Port:                  port,
				URL:                   endpoint.Model.Image.Tei.URL,
				MaxBatchTokens:        endpoint.Model.Image.Tei.MaxBatchTokens,
				MaxConcurrentRequests: endpoint.Model.Image.Tei.MaxConcurrentRequests,
				Pooling:               endpoint.Model.Image.Tei.Pooling,
			},
		}
	} else if endpoint.Model.Image.Vllm != nil {
		var port *int
		if endpoint.Model.Image.Vllm.Port.IsUnknown() || endpoint.Model.Image.Vllm.Port.IsNull() {
			port = nil
		} else {
			portInt := int(endpoint.Model.Image.Vllm.Port.ValueInt64())
			port = &portInt
		}
		image = huggingface.Image{
			Vllm: &huggingface.Vllm{
				HealthRoute:         endpoint.Model.Image.Vllm.HealthRoute,
				Port:                port,
				URL:                 endpoint.Model.Image.Vllm.URL,
				KvCacheDtype:        endpoint.Model.Image.Vllm.KvCacheDtype,
				MaxNumBatchedTokens: endpoint.Model.Image.Vllm.MaxNumBatchedTokens,
				MaxNumSeqs:          endpoint.Model.Image.Vllm.MaxNumSeqs,
				TensorParallelSize:  endpoint.Model.Image.Vllm.TensorParallelSize,
				ServerArgs:          endpoint.Model.Image.Vllm.ServerArgs,
			},
		}
	}

	var timeout *int
	if endpoint.Compute.Scaling.ScaleToZeroTimeout.IsUnknown() || endpoint.Compute.Scaling.ScaleToZeroTimeout.IsNull() {
		timeout = nil
	} else {
		timeoutInt := int(endpoint.Compute.Scaling.ScaleToZeroTimeout.ValueInt64())
		timeout = &timeoutInt
	}

	var measure *huggingface.Measure = nil
	if endpoint.Compute.Scaling.Measure != nil {
		measure = &huggingface.Measure{
			HardwareUsage:   endpoint.Compute.Scaling.Measure.HardwareUsage,
			PendingRequests: endpoint.Compute.Scaling.Measure.PendingRequests,
		}
	}

	huggingfaceEndpoint := huggingface.UpdateEndpointRequest{
		Compute: &huggingface.Compute{
			Accelerator:  endpoint.Compute.Accelerator,
			InstanceSize: endpoint.Compute.InstanceSize,
			InstanceType: endpoint.Compute.InstanceType,
			Scaling: huggingface.Scaling{
				MaxReplica:         endpoint.Compute.Scaling.MaxReplica,
				MinReplica:         endpoint.Compute.Scaling.MinReplica,
				ScaleToZeroTimeout: timeout,
				Measure:            measure,
			},
		},
		Model: &huggingface.Model{
			Framework:  endpoint.Model.Framework,
			Image:      image,
			Repository: endpoint.Model.Repository,
			Revision:   endpoint.Model.Revision.ValueStringPointer(),
			Task:       endpoint.Model.Task.ValueStringPointer(),
			Env:        endpoint.Model.Env,
		},
		Type: endpoint.Type.ValueStringPointer(),
	}

	return huggingfaceEndpoint
}

func (r *endpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan endpointResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	existingEndpoints, err := r.client.ListEndpoints()
	if err != nil {
		resp.Diagnostics.AddError(
			"error listing endpoints",
			err.Error(),
		)
		return
	}

	useUpdate := false
	for _, existingEndpoint := range existingEndpoints {
		if existingEndpoint.Name == plan.Name.ValueString() {
			useUpdate = true
			break
		}
	}

	var createdEndpoint huggingface.EndpointDetails

	if useUpdate {
		updateEndpointRequest := providerEndpointToUpdateEndpointRequest(plan)
		createdEndpoint, err = r.client.UpdateEndpoint(plan.Name.ValueString(), updateEndpointRequest)
	} else {
		createEndpointRequest := providerEndpointToCreateEndpointRequest(plan)
		createdEndpoint, err = r.client.CreateEndpoint(createEndpointRequest)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"error creating endpoint",
			err.Error(),
		)
		return
	}

	plan = clientEndpointToProviderEndpoint(createdEndpoint)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *endpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state endpointResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint, err := r.client.GetEndpoint(state.Name.ValueString())
	if err != nil {
		httpErr, ok := err.(*huggingface.HTTPError)
		if ok && httpErr.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError(
				"error reading endpoint",
				"could not read endpoint named "+state.Name.ValueString()+": "+err.Error(),
			)
			return
		}
	}

	state = clientEndpointToProviderEndpoint(endpoint)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *endpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan endpointResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := providerEndpointToUpdateEndpointRequest(plan)

	updatedEndpoint, err := r.client.UpdateEndpoint(plan.Name.ValueString(), endpoint)
	if err != nil {
		resp.Diagnostics.AddError(
			"error updating endpoint",
			err.Error(),
		)
		return
	}

	plan = clientEndpointToProviderEndpoint(updatedEndpoint)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *endpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state endpointResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteEndpoint(state.Name.ValueString())
	if err != nil {
		if httpErr, ok := err.(*huggingface.HTTPError); ok && httpErr.StatusCode != http.StatusNotFound {
			resp.Diagnostics.AddError(
				"error deleting endpoint",
				err.Error(),
			)
			return
		}
	}

	resp.State.RemoveResource(ctx)
}
