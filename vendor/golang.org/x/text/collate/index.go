



package collate

import "golang.org/x/text/internal/colltab"

const blockSize = 64

func getTable(t tableIndex) *colltab.Table {
	return &colltab.Table{
		Index: colltab.Trie{
			Index0:  mainLookup[:][blockSize*t.lookupOffset:],
			Values0: mainValues[:][blockSize*t.valuesOffset:],
			Index:   mainLookup[:],
			Values:  mainValues[:],
		},
		ExpandElem:     mainExpandElem[:],
		ContractTries:  colltab.ContractTrieSet(mainCTEntries[:]),
		ContractElem:   mainContractElem[:],
		MaxContractLen: 18,
		VariableTop:    varTop,
	}
}



type tableIndex struct {
	lookupOffset uint32
	valuesOffset uint32
}
