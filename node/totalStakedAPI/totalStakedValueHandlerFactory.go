package totalStakedAPI

import (
	"time"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-go/marshal"
	"github.com/ElrondNetwork/elrond-go/node/external"
)

// ArgsTotalStakedValueHandler is struct that contains components that are needed to create a TotalStakedValueHandler
type ArgsTotalStakedValueHandler struct {
	ShardID                     uint32
	RoundDurationInMilliseconds uint64
	InternalMarshalizer         marshal.Marshalizer
	Accounts                    state.AccountsAdapter
}

const numOfRounds = 10

// CreateTotalStakedValueHandler wil create a new instance of TotalStakedValueHandler
func CreateTotalStakedValueHandler(args *ArgsTotalStakedValueHandler) (external.TotalStakedValueHandler, error) {
	if args.ShardID != core.MetachainShardId {
		return NewDisabledTotalStakedValueProcessor()
	}

	return NewTotalStakedValueProcessor(
		args.InternalMarshalizer,
		time.Duration(args.RoundDurationInMilliseconds)*time.Millisecond*numOfRounds,
		args.Accounts,
	)
}
