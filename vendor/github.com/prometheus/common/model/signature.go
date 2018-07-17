












package model

import (
	"sort"
)




const SeparatorByte byte = 255

var (
	
	emptyLabelSignature = hashNew()
)




func LabelsToSignature(labels map[string]string) uint64 {
	if len(labels) == 0 {
		return emptyLabelSignature
	}

	labelNames := make([]string, 0, len(labels))
	for labelName := range labels {
		labelNames = append(labelNames, labelName)
	}
	sort.Strings(labelNames)

	sum := hashNew()
	for _, labelName := range labelNames {
		sum = hashAdd(sum, labelName)
		sum = hashAddByte(sum, SeparatorByte)
		sum = hashAdd(sum, labels[labelName])
		sum = hashAddByte(sum, SeparatorByte)
	}
	return sum
}



func labelSetToFingerprint(ls LabelSet) Fingerprint {
	if len(ls) == 0 {
		return Fingerprint(emptyLabelSignature)
	}

	labelNames := make(LabelNames, 0, len(ls))
	for labelName := range ls {
		labelNames = append(labelNames, labelName)
	}
	sort.Sort(labelNames)

	sum := hashNew()
	for _, labelName := range labelNames {
		sum = hashAdd(sum, string(labelName))
		sum = hashAddByte(sum, SeparatorByte)
		sum = hashAdd(sum, string(ls[labelName]))
		sum = hashAddByte(sum, SeparatorByte)
	}
	return Fingerprint(sum)
}




func labelSetToFastFingerprint(ls LabelSet) Fingerprint {
	if len(ls) == 0 {
		return Fingerprint(emptyLabelSignature)
	}

	var result uint64
	for labelName, labelValue := range ls {
		sum := hashNew()
		sum = hashAdd(sum, string(labelName))
		sum = hashAddByte(sum, SeparatorByte)
		sum = hashAdd(sum, string(labelValue))
		result ^= sum
	}
	return Fingerprint(result)
}





func SignatureForLabels(m Metric, labels ...LabelName) uint64 {
	if len(labels) == 0 {
		return emptyLabelSignature
	}

	sort.Sort(LabelNames(labels))

	sum := hashNew()
	for _, label := range labels {
		sum = hashAdd(sum, string(label))
		sum = hashAddByte(sum, SeparatorByte)
		sum = hashAdd(sum, string(m[label]))
		sum = hashAddByte(sum, SeparatorByte)
	}
	return sum
}




func SignatureWithoutLabels(m Metric, labels map[LabelName]struct{}) uint64 {
	if len(m) == 0 {
		return emptyLabelSignature
	}

	labelNames := make(LabelNames, 0, len(m))
	for labelName := range m {
		if _, exclude := labels[labelName]; !exclude {
			labelNames = append(labelNames, labelName)
		}
	}
	if len(labelNames) == 0 {
		return emptyLabelSignature
	}
	sort.Sort(labelNames)

	sum := hashNew()
	for _, labelName := range labelNames {
		sum = hashAdd(sum, string(labelName))
		sum = hashAddByte(sum, SeparatorByte)
		sum = hashAdd(sum, string(m[labelName]))
		sum = hashAddByte(sum, SeparatorByte)
	}
	return sum
}
