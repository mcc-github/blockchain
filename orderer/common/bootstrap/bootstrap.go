


package bootstrap

import (
	ab "github.com/mcc-github/blockchain-protos-go/common"
)


type Helper interface {
	
	
	GenesisBlock() *ab.Block
}







type Replacer interface {
	
	
	
	
	
	
	ReplaceGenesisBlockFile(block *ab.Block) error

	
	
	
	
	
	CheckReadWrite() error
}
