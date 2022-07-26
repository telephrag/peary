package bot_errors

import "errors"

var (
	ErrFailedGiveRole     = errors.New("failed to give role")
	ErrFailedTakeRole     = errors.New("failed to take role")
	ErrFailedSendResponse = errors.New("failed to respond to player")
	ErrFailedToRecover    = errors.New("failed to recover from error")
	ErrSomewhereElse      = errors.New("error occured somewhere else and bot failed to recover")
	ErrHandlerTimeout     = errors.New("command handler execution timeout")
)

const (
	CmdDeploy         = "cmd_deploy"
	CmdDeployDo       = "cmd_deploy_do"
	CmdDeployRollback = "cmd_deploy_rollback"

	CmdReturn         = "cmd_return"
	CmdReturnDo       = "cmd_return_do"
	CmdReturnRollback = "cmd_return_rollback"

	DBChangeStream = "db_changestream"
	DBInsert       = "db_insert"
	DBDelete       = "db_delete"
	DBRoleExpire   = "db_role_expire"

	NotifyUsr = "usr_notify"

	CtxCancel = "ctx_cancel"
)
