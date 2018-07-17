/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package statecouchdb

import (
	"sync"
)


type batch interface {
	execute() error
}



func executeBatches(batches []batch) error {
	logger.Debugf("Executing batches = %s", batches)
	numBatches := len(batches)
	if numBatches == 0 {
		return nil
	}
	if numBatches == 1 {
		return batches[0].execute()
	}
	var batchWG sync.WaitGroup
	batchWG.Add(numBatches)
	errsChan := make(chan error, numBatches)
	defer close(errsChan)
	for _, b := range batches {
		go func(b batch) {
			defer batchWG.Done()
			if err := b.execute(); err != nil {
				errsChan <- err
				return
			}
		}(b)
	}
	batchWG.Wait()
	if len(errsChan) > 0 {
		return <-errsChan
	}
	return nil
}
