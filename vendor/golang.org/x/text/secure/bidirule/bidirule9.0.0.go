





package bidirule

func (t *Transformer) isFinal() bool {
	if !t.isRTL() {
		return true
	}
	return t.state == ruleLTRFinal || t.state == ruleRTLFinal || t.state == ruleInitial
}
