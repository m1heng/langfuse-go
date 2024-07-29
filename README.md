# Langfuse Go SDK


[![GoDoc](https://godoc.org/github.com/m1heng/langfuse-go?status.svg)](https://godoc.org/github.com/m1heng/langfuse-go) [![Go Report Card](https://goreportcard.com/badge/github.com/m1heng/langfuse-go)](https://goreportcard.com/report/github.com/m1heng/langfuse-go) [![GitHub release](https://img.shields.io/github/release/henomis/langfuse-go.svg)](https://github.com/m1heng/langfuse-go/releases)

This is [Langfuse](https://langfuse.com)'s **unofficial** Go client, designed to enable you to use Langfuse's services easily from your own applications.

## Langfuse

[Langfuse](https://langfuse.com) traces, evals, prompt management and metrics to debug and improve your LLM application.


## API support

| **Index Operations**  | **Status** |
| --- | --- |
| Trace | 🟢 | 
| Generation | 🟢 |
| Span | 🟢 |
| Event | 🟢 |
| Score | 🟢 |
| Prompt | 🔧 |




## Getting started

### Installation

You can load langfuse-go into your project by using:
```
go get github.com/m1heng/langfuse-go
```


### Configuration
You can config during the construction of the Langfuse client:
```go
l, _ := langfuse.New(context.Background(), &langfuse.Config{
	ApiClientConfig: &langfuse.APIConfig{
		LangfuseHost: "https://cloud.langfuse.com",
		PublicKey:    "public-key",
		SecretKey:    "secret-key",
	},
	AutoFlushInterval: 500 * time.Millisecond,
})
```

Or leave `ApiClientConfig` empty and just like the official Python SDK, these three environment variables will be used to configure the Langfuse client:

- `LANGFUSE_HOST`: The host of the Langfuse service.
- `LANGFUSE_PUBLIC_KEY`: Your public key for the Langfuse service.
- `LANGFUSE_SECRET_KEY`: Your secret key for the Langfuse service.


### Usage
Please refer to the [examples folder](examples/cmd/) to see how to use the SDK.


## Originally forked from
https://github.com/henomis/langfuse-go
