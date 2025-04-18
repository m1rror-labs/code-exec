package pkg

type Dependencies struct {
	Repo          Repository
	TsRuntime     CodeExecutor
	RustRuntime   CodeExecutor
	AnchorRuntime ProgramBuilder
	RpcEngine     RpcEngine
}

type Repository interface {
	TransactionRepo
}
