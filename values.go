package club

import "context"

//-----------------------------------------------------------------------------

// Global app context management
var (
	AppPool   *Group
	AppCtx    context.Context
	AppCancel context.CancelFunc
)

//-----------------------------------------------------------------------------

// Constants
const (
	ErrTimeout = Error(`TIMEOUT`)
)

//-----------------------------------------------------------------------------
