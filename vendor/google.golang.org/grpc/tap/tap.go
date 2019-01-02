



package tap

import (
	"context"
)


type Info struct {
	
	
	FullMethodName string
	
}

















type ServerInHandle func(ctx context.Context, info *Info) (context.Context, error)
