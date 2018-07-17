package gomega

import (
	"time"

	"github.com/onsi/gomega/matchers"
	"github.com/onsi/gomega/types"
)




func Equal(expected interface{}) types.GomegaMatcher {
	return &matchers.EqualMatcher{
		Expected: expected,
	}
}





func BeEquivalentTo(expected interface{}) types.GomegaMatcher {
	return &matchers.BeEquivalentToMatcher{
		Expected: expected,
	}
}




func BeIdenticalTo(expected interface{}) types.GomegaMatcher {
	return &matchers.BeIdenticalToMatcher{
		Expected: expected,
	}
}


func BeNil() types.GomegaMatcher {
	return &matchers.BeNilMatcher{}
}


func BeTrue() types.GomegaMatcher {
	return &matchers.BeTrueMatcher{}
}


func BeFalse() types.GomegaMatcher {
	return &matchers.BeFalseMatcher{}
}





func HaveOccurred() types.GomegaMatcher {
	return &matchers.HaveOccurredMatcher{}
}












func Succeed() types.GomegaMatcher {
	return &matchers.SucceedMatcher{}
}








func MatchError(expected interface{}) types.GomegaMatcher {
	return &matchers.MatchErrorMatcher{
		Expected: expected,
	}
}












func BeClosed() types.GomegaMatcher {
	return &matchers.BeClosedMatcher{}
}




































func Receive(args ...interface{}) types.GomegaMatcher {
	var arg interface{}
	if len(args) > 0 {
		arg = args[0]
	}

	return &matchers.ReceiveMatcher{
		Arg: arg,
	}
}













func BeSent(arg interface{}) types.GomegaMatcher {
	return &matchers.BeSentMatcher{
		Arg: arg,
	}
}




func MatchRegexp(regexp string, args ...interface{}) types.GomegaMatcher {
	return &matchers.MatchRegexpMatcher{
		Regexp: regexp,
		Args:   args,
	}
}




func ContainSubstring(substr string, args ...interface{}) types.GomegaMatcher {
	return &matchers.ContainSubstringMatcher{
		Substr: substr,
		Args:   args,
	}
}




func HavePrefix(prefix string, args ...interface{}) types.GomegaMatcher {
	return &matchers.HavePrefixMatcher{
		Prefix: prefix,
		Args:   args,
	}
}




func HaveSuffix(suffix string, args ...interface{}) types.GomegaMatcher {
	return &matchers.HaveSuffixMatcher{
		Suffix: suffix,
		Args:   args,
	}
}




func MatchJSON(json interface{}) types.GomegaMatcher {
	return &matchers.MatchJSONMatcher{
		JSONToMatch: json,
	}
}




func MatchXML(xml interface{}) types.GomegaMatcher {
	return &matchers.MatchXMLMatcher{
		XMLToMatch: xml,
	}
}




func MatchYAML(yaml interface{}) types.GomegaMatcher {
	return &matchers.MatchYAMLMatcher{
		YAMLToMatch: yaml,
	}
}


func BeEmpty() types.GomegaMatcher {
	return &matchers.BeEmptyMatcher{}
}


func HaveLen(count int) types.GomegaMatcher {
	return &matchers.HaveLenMatcher{
		Count: count,
	}
}


func HaveCap(count int) types.GomegaMatcher {
	return &matchers.HaveCapMatcher{
		Count: count,
	}
}


func BeZero() types.GomegaMatcher {
	return &matchers.BeZeroMatcher{}
}








func ContainElement(element interface{}) types.GomegaMatcher {
	return &matchers.ContainElementMatcher{
		Element: element,
	}
}
















func ConsistOf(elements ...interface{}) types.GomegaMatcher {
	return &matchers.ConsistOfMatcher{
		Elements: elements,
	}
}





func HaveKey(key interface{}) types.GomegaMatcher {
	return &matchers.HaveKeyMatcher{
		Key: key,
	}
}






func HaveKeyWithValue(key interface{}, value interface{}) types.GomegaMatcher {
	return &matchers.HaveKeyWithValueMatcher{
		Key:   key,
		Value: value,
	}
}












func BeNumerically(comparator string, compareTo ...interface{}) types.GomegaMatcher {
	return &matchers.BeNumericallyMatcher{
		Comparator: comparator,
		CompareTo:  compareTo,
	}
}





func BeTemporally(comparator string, compareTo time.Time, threshold ...time.Duration) types.GomegaMatcher {
	return &matchers.BeTemporallyMatcher{
		Comparator: comparator,
		CompareTo:  compareTo,
		Threshold:  threshold,
	}
}







func BeAssignableToTypeOf(expected interface{}) types.GomegaMatcher {
	return &matchers.AssignableToTypeOfMatcher{
		Expected: expected,
	}
}



func Panic() types.GomegaMatcher {
	return &matchers.PanicMatcher{}
}



func BeAnExistingFile() types.GomegaMatcher {
	return &matchers.BeAnExistingFileMatcher{}
}



func BeARegularFile() types.GomegaMatcher {
	return &matchers.BeARegularFileMatcher{}
}



func BeADirectory() types.GomegaMatcher {
	return &matchers.BeADirectoryMatcher{}
}






func And(ms ...types.GomegaMatcher) types.GomegaMatcher {
	return &matchers.AndMatcher{Matchers: ms}
}



func SatisfyAll(matchers ...types.GomegaMatcher) types.GomegaMatcher {
	return And(matchers...)
}






func Or(ms ...types.GomegaMatcher) types.GomegaMatcher {
	return &matchers.OrMatcher{Matchers: ms}
}



func SatisfyAny(matchers ...types.GomegaMatcher) types.GomegaMatcher {
	return Or(matchers...)
}





func Not(matcher types.GomegaMatcher) types.GomegaMatcher {
	return &matchers.NotMatcher{Matcher: matcher}
}







func WithTransform(transform interface{}, matcher types.GomegaMatcher) types.GomegaMatcher {
	return matchers.NewWithTransformMatcher(transform, matcher)
}
