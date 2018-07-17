package spec_iterator

func ParallelizedIndexRange(length int, parallelTotal int, parallelNode int) (startIndex int, count int) {
	if length == 0 {
		return 0, 0
	}

	
	if parallelTotal >= length {
		if parallelNode > length {
			return 0, 0
		} else {
			return parallelNode - 1, 1
		}
	}

	
	minTestsPerNode := length / parallelTotal

	
	
	
	maxTestsPerNode := minTestsPerNode
	if length%parallelTotal != 0 {
		maxTestsPerNode++
	}

	
	numMaxLoadNodes := length % parallelTotal

	
	var numPrecedingMaxLoadNodes int
	if parallelNode > numMaxLoadNodes {
		numPrecedingMaxLoadNodes = numMaxLoadNodes
	} else {
		numPrecedingMaxLoadNodes = parallelNode - 1
	}

	
	var numPrecedingMinLoadNodes int
	if parallelNode <= numMaxLoadNodes {
		numPrecedingMinLoadNodes = 0
	} else {
		numPrecedingMinLoadNodes = parallelNode - numMaxLoadNodes - 1
	}

	
	startIndex = numPrecedingMaxLoadNodes*maxTestsPerNode + numPrecedingMinLoadNodes*minTestsPerNode
	if parallelNode > numMaxLoadNodes {
		count = minTestsPerNode
	} else {
		count = maxTestsPerNode
	}
	return
}
