


package bootstrap

import (
	ab "github.com/mcc-github/blockchain/protos/common"
)


type Helper interface {
	
	
	GenesisBlock() *ab.Block
}







type Replacer interface {
	
	
	
	
	
	
	ReplaceGenesisBlockFile(block *ab.Block) error

	
	
	
	
	
	CheckReadWrite() error
}
