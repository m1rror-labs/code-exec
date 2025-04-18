package pkg

type Dependencies struct {
	TsRuntime     CodeExecutor
	RustRuntime   CodeExecutor
	AnchorRuntime ProgramBuilder
	RpcEngine     RpcEngine
}
