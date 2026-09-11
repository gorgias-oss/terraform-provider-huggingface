package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Endpoint struct {
	AccountId *string `tfsdk:"account_id"`
	Compute   Compute `tfsdk:"compute"`
	Model     Model   `tfsdk:"model"`
	Name      string  `tfsdk:"name"`
	Cloud     Cloud   `tfsdk:"cloud"`
	Status    *Status `tfsdk:"status"`
	Type      string  `tfsdk:"type"`
}

type Compute struct {
	Accelerator  string  `tfsdk:"accelerator"`
	InstanceSize string  `tfsdk:"instance_size"`
	InstanceType string  `tfsdk:"instance_type"`
	Scaling      Scaling `tfsdk:"scaling"`
}

type Measure struct {
	HardwareUsage   *float64 `tfsdk:"hardware_usage"`
	PendingRequests *float64 `tfsdk:"pending_requests"`
}

type Scaling struct {
	MaxReplica         int         `tfsdk:"max_replica"`
	MinReplica         int         `tfsdk:"min_replica"`
	ScaleToZeroTimeout types.Int64 `tfsdk:"scale_to_zero_timeout"`
	Measure            *Measure    `tfsdk:"measure"`
}

type Model struct {
	Framework  string            `tfsdk:"framework"`
	Image      Image             `tfsdk:"image"`
	Repository string            `tfsdk:"repository"`
	Revision   types.String      `tfsdk:"revision"`
	Task       types.String      `tfsdk:"task"`
	Env        map[string]string `tfsdk:"env"`
}

type Image struct {
	Huggingface *Huggingface `tfsdk:"huggingface"`
	Custom      *Custom      `tfsdk:"custom"`
	Tei         *Tei         `tfsdk:"tei"`
	Tgi         *Tgi         `tfsdk:"tgi"`
	TgiNeuron   *TgiNeuron   `tfsdk:"tgi_neuron"`
	Llamacpp    *Llamacpp    `tfsdk:"llamacpp"`
	Vllm        *Vllm        `tfsdk:"vllm"`
}

type Tei struct {
	HealthRoute           *string     `tfsdk:"health_route"`
	Port                  types.Int64 `tfsdk:"port"`
	URL                   string      `tfsdk:"url"`
	MaxBatchTokens        *int        `tfsdk:"max_batch_tokens"`
	MaxConcurrentRequests *int        `tfsdk:"max_concurrent_requests"`
	Pooling               *string     `tfsdk:"pooling"`
}

type Llamacpp struct {
	HealthRoute *string     `tfsdk:"health_route"`
	Port        types.Int64 `tfsdk:"port"`
	URL         string      `tfsdk:"url"`
	CtxSize     *int        `tfsdk:"ctx_size"`
	Embeddings  *bool       `tfsdk:"embeddings"`
	ModelPath   string      `tfsdk:"modelPath"`
	NParallel   *int        `tfsdk:"n_parallel"`
	ThreadsHttp *int        `tfsdk:"threads_http"`
}

type TgiNeuron struct {
	HealthRoute           *string     `tfsdk:"health_route"`
	Port                  types.Int64 `tfsdk:"port"`
	URL                   string      `tfsdk:"url"`
	MaxBatchPrefillTokens *int        `tfsdk:"max_batch_prefill_tokens"`
	MaxBatchTotalTokens   *int        `tfsdk:"max_batch_total_tokens"`
	MaxInputLength        *int        `tfsdk:"max_input_length"`
	MaxTotalTokens        *int        `tfsdk:"max_total_tokens"`
	HfAutoCastType        *string     `tfsdk:"hf_auto_cast_type"`
	HfNumCores            *int        `tfsdk:"hf_num_cores"`
}

type Tgi struct {
	HealthRoute           *string     `tfsdk:"health_route"`
	Port                  types.Int64 `tfsdk:"port"`
	URL                   string      `tfsdk:"url"`
	MaxBatchPrefillTokens *int        `tfsdk:"max_batch_prefill_tokens"`
	MaxBatchTotalTokens   *int        `tfsdk:"max_batch_total_tokens"`
	MaxInputLength        *int        `tfsdk:"max_input_length"`
	MaxTotalTokens        *int        `tfsdk:"max_total_tokens"`
	DisableCustomKernels  *bool       `tfsdk:"disable_custom_kernels"`
	Quantize              *string     `tfsdk:"quantize"`
}

type Custom struct {
	Credentials *Credentials `tfsdk:"credentials"`
	HealthRoute *string      `tfsdk:"health_route"`
	Port        types.Int64  `tfsdk:"port"`
	URL         string       `tfsdk:"url"`
}

type Credentials struct {
	Password string `tfsdk:"password"`
	Username string `tfsdk:"username"`
}

type Huggingface struct{}

type Cloud struct {
	Region string `tfsdk:"region"`
	Vendor string `tfsdk:"vendor"`
}

type Status struct {
	CreatedAt     string  `tfsdk:"created_at"`
	CreatedBy     User    `tfsdk:"created_by"`
	ErrorMessage  string  `tfsdk:"error_message"`
	Message       string  `tfsdk:"message"`
	Private       Private `tfsdk:"private"`
	ReadyReplica  int     `tfsdk:"ready_replica"`
	State         string  `tfsdk:"state"`
	TargetReplica int     `tfsdk:"target_replica"`
	UpdatedAt     string  `tfsdk:"updated_at"`
	UpdatedBy     User    `tfsdk:"updated_by"`
	URL           string  `tfsdk:"url"`
}

type User struct {
	ID   string `tfsdk:"id"`
	Name string `tfsdk:"name"`
}

type Private struct {
	ServiceName string `tfsdk:"service_name"`
}

type Vllm struct {
	HealthRoute         *string     `tfsdk:"health_route"`
	Port                types.Int64 `tfsdk:"port"`
	URL                 string      `tfsdk:"url"`
	KvCacheDtype        *string     `tfsdk:"kv_cache_dtype"`
	MaxNumBatchedTokens *int        `tfsdk:"max_num_batched_tokens"`
	MaxNumSeqs          *int        `tfsdk:"max_num_seqs"`
	TensorParallelSize  *int        `tfsdk:"tensor_parallel_size"`
	ServerArgs          []string    `tfsdk:"server_args"`
}
