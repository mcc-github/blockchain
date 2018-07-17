



package language





func MustParse(s string) Tag {
	t, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return t
}



func (c CanonType) MustParse(s string) Tag {
	t, err := c.Parse(s)
	if err != nil {
		panic(err)
	}
	return t
}



func MustParseBase(s string) Base {
	b, err := ParseBase(s)
	if err != nil {
		panic(err)
	}
	return b
}



func MustParseScript(s string) Script {
	scr, err := ParseScript(s)
	if err != nil {
		panic(err)
	}
	return scr
}



func MustParseRegion(s string) Region {
	r, err := ParseRegion(s)
	if err != nil {
		panic(err)
	}
	return r
}

var (
	und = Tag{}

	Und Tag = Tag{}

	Afrikaans            Tag = Tag{lang: _af}                
	Amharic              Tag = Tag{lang: _am}                
	Arabic               Tag = Tag{lang: _ar}                
	ModernStandardArabic Tag = Tag{lang: _ar, region: _001}  
	Azerbaijani          Tag = Tag{lang: _az}                
	Bulgarian            Tag = Tag{lang: _bg}                
	Bengali              Tag = Tag{lang: _bn}                
	Catalan              Tag = Tag{lang: _ca}                
	Czech                Tag = Tag{lang: _cs}                
	Danish               Tag = Tag{lang: _da}                
	German               Tag = Tag{lang: _de}                
	Greek                Tag = Tag{lang: _el}                
	English              Tag = Tag{lang: _en}                
	AmericanEnglish      Tag = Tag{lang: _en, region: _US}   
	BritishEnglish       Tag = Tag{lang: _en, region: _GB}   
	Spanish              Tag = Tag{lang: _es}                
	EuropeanSpanish      Tag = Tag{lang: _es, region: _ES}   
	LatinAmericanSpanish Tag = Tag{lang: _es, region: _419}  
	Estonian             Tag = Tag{lang: _et}                
	Persian              Tag = Tag{lang: _fa}                
	Finnish              Tag = Tag{lang: _fi}                
	Filipino             Tag = Tag{lang: _fil}               
	French               Tag = Tag{lang: _fr}                
	CanadianFrench       Tag = Tag{lang: _fr, region: _CA}   
	Gujarati             Tag = Tag{lang: _gu}                
	Hebrew               Tag = Tag{lang: _he}                
	Hindi                Tag = Tag{lang: _hi}                
	Croatian             Tag = Tag{lang: _hr}                
	Hungarian            Tag = Tag{lang: _hu}                
	Armenian             Tag = Tag{lang: _hy}                
	Indonesian           Tag = Tag{lang: _id}                
	Icelandic            Tag = Tag{lang: _is}                
	Italian              Tag = Tag{lang: _it}                
	Japanese             Tag = Tag{lang: _ja}                
	Georgian             Tag = Tag{lang: _ka}                
	Kazakh               Tag = Tag{lang: _kk}                
	Khmer                Tag = Tag{lang: _km}                
	Kannada              Tag = Tag{lang: _kn}                
	Korean               Tag = Tag{lang: _ko}                
	Kirghiz              Tag = Tag{lang: _ky}                
	Lao                  Tag = Tag{lang: _lo}                
	Lithuanian           Tag = Tag{lang: _lt}                
	Latvian              Tag = Tag{lang: _lv}                
	Macedonian           Tag = Tag{lang: _mk}                
	Malayalam            Tag = Tag{lang: _ml}                
	Mongolian            Tag = Tag{lang: _mn}                
	Marathi              Tag = Tag{lang: _mr}                
	Malay                Tag = Tag{lang: _ms}                
	Burmese              Tag = Tag{lang: _my}                
	Nepali               Tag = Tag{lang: _ne}                
	Dutch                Tag = Tag{lang: _nl}                
	Norwegian            Tag = Tag{lang: _no}                
	Punjabi              Tag = Tag{lang: _pa}                
	Polish               Tag = Tag{lang: _pl}                
	Portuguese           Tag = Tag{lang: _pt}                
	BrazilianPortuguese  Tag = Tag{lang: _pt, region: _BR}   
	EuropeanPortuguese   Tag = Tag{lang: _pt, region: _PT}   
	Romanian             Tag = Tag{lang: _ro}                
	Russian              Tag = Tag{lang: _ru}                
	Sinhala              Tag = Tag{lang: _si}                
	Slovak               Tag = Tag{lang: _sk}                
	Slovenian            Tag = Tag{lang: _sl}                
	Albanian             Tag = Tag{lang: _sq}                
	Serbian              Tag = Tag{lang: _sr}                
	SerbianLatin         Tag = Tag{lang: _sr, script: _Latn} 
	Swedish              Tag = Tag{lang: _sv}                
	Swahili              Tag = Tag{lang: _sw}                
	Tamil                Tag = Tag{lang: _ta}                
	Telugu               Tag = Tag{lang: _te}                
	Thai                 Tag = Tag{lang: _th}                
	Turkish              Tag = Tag{lang: _tr}                
	Ukrainian            Tag = Tag{lang: _uk}                
	Urdu                 Tag = Tag{lang: _ur}                
	Uzbek                Tag = Tag{lang: _uz}                
	Vietnamese           Tag = Tag{lang: _vi}                
	Chinese              Tag = Tag{lang: _zh}                
	SimplifiedChinese    Tag = Tag{lang: _zh, script: _Hans} 
	TraditionalChinese   Tag = Tag{lang: _zh, script: _Hant} 
	Zulu                 Tag = Tag{lang: _zu}                
)
