package bot_errors

const (
	ErrFailedGiveRole     = "failed to give role"
	ErrFailedTakeRole     = "failed to take role"
	ErrFailedSendResponse = "failed to respond to player"
	ErrFailedToRecover    = "failed to recover from error"
	ErrSomewhereElse      = "error occured somewhere else and bot failed to recover"
	ErrHandlerTimeout     = "command handler execution timeout"
)

const (
	CmdDeploy         = "cmd_deploy"
	CmdDeployDo       = "cmd_deploy_do"
	CmdDeployRollback = "cmd_deploy_rollback"

	CmdReturn        = "cmd_return"
	CmdReturnDo      = "cmd_return_do"
	CmdReturnRolback = "cmd_return_rollback"

	DBRoleExpire = "db_role_expire"
	NotifyUsr    = "notify_user"
)
