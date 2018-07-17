

package cldr


type LDMLBCP47 struct {
	Common
	Version *struct {
		Common
		Number string `xml:"number,attr"`
	} `xml:"version"`
	Generation *struct {
		Common
		Date string `xml:"date,attr"`
	} `xml:"generation"`
	Keyword []*struct {
		Common
		Key []*struct {
			Common
			Extension   string `xml:"extension,attr"`
			Name        string `xml:"name,attr"`
			Description string `xml:"description,attr"`
			Deprecated  string `xml:"deprecated,attr"`
			Preferred   string `xml:"preferred,attr"`
			Alias       string `xml:"alias,attr"`
			ValueType   string `xml:"valueType,attr"`
			Since       string `xml:"since,attr"`
			Type        []*struct {
				Common
				Name        string `xml:"name,attr"`
				Description string `xml:"description,attr"`
				Deprecated  string `xml:"deprecated,attr"`
				Preferred   string `xml:"preferred,attr"`
				Alias       string `xml:"alias,attr"`
				Since       string `xml:"since,attr"`
			} `xml:"type"`
		} `xml:"key"`
	} `xml:"keyword"`
	Attribute []*struct {
		Common
		Name        string `xml:"name,attr"`
		Description string `xml:"description,attr"`
		Deprecated  string `xml:"deprecated,attr"`
		Preferred   string `xml:"preferred,attr"`
		Since       string `xml:"since,attr"`
	} `xml:"attribute"`
}



type SupplementalData struct {
	Common
	Version *struct {
		Common
		Number string `xml:"number,attr"`
	} `xml:"version"`
	Generation *struct {
		Common
		Date string `xml:"date,attr"`
	} `xml:"generation"`
	CurrencyData *struct {
		Common
		Fractions []*struct {
			Common
			Info []*struct {
				Common
				Iso4217      string `xml:"iso4217,attr"`
				Digits       string `xml:"digits,attr"`
				Rounding     string `xml:"rounding,attr"`
				CashDigits   string `xml:"cashDigits,attr"`
				CashRounding string `xml:"cashRounding,attr"`
			} `xml:"info"`
		} `xml:"fractions"`
		Region []*struct {
			Common
			Iso3166  string `xml:"iso3166,attr"`
			Currency []*struct {
				Common
				Before       string `xml:"before,attr"`
				From         string `xml:"from,attr"`
				To           string `xml:"to,attr"`
				Iso4217      string `xml:"iso4217,attr"`
				Digits       string `xml:"digits,attr"`
				Rounding     string `xml:"rounding,attr"`
				CashRounding string `xml:"cashRounding,attr"`
				Tender       string `xml:"tender,attr"`
				Alternate    []*struct {
					Common
					Iso4217 string `xml:"iso4217,attr"`
				} `xml:"alternate"`
			} `xml:"currency"`
		} `xml:"region"`
	} `xml:"currencyData"`
	TerritoryContainment *struct {
		Common
		Group []*struct {
			Common
			Contains string `xml:"contains,attr"`
			Grouping string `xml:"grouping,attr"`
			Status   string `xml:"status,attr"`
		} `xml:"group"`
	} `xml:"territoryContainment"`
	SubdivisionContainment *struct {
		Common
		Subgroup []*struct {
			Common
			Subtype  string `xml:"subtype,attr"`
			Contains string `xml:"contains,attr"`
		} `xml:"subgroup"`
	} `xml:"subdivisionContainment"`
	LanguageData *struct {
		Common
		Language []*struct {
			Common
			Scripts     string `xml:"scripts,attr"`
			Territories string `xml:"territories,attr"`
			Variants    string `xml:"variants,attr"`
		} `xml:"language"`
	} `xml:"languageData"`
	TerritoryInfo *struct {
		Common
		Territory []*struct {
			Common
			Gdp                string `xml:"gdp,attr"`
			LiteracyPercent    string `xml:"literacyPercent,attr"`
			Population         string `xml:"population,attr"`
			LanguagePopulation []*struct {
				Common
				LiteracyPercent   string `xml:"literacyPercent,attr"`
				WritingPercent    string `xml:"writingPercent,attr"`
				PopulationPercent string `xml:"populationPercent,attr"`
				OfficialStatus    string `xml:"officialStatus,attr"`
			} `xml:"languagePopulation"`
		} `xml:"territory"`
	} `xml:"territoryInfo"`
	PostalCodeData *struct {
		Common
		PostCodeRegex []*struct {
			Common
			TerritoryId string `xml:"territoryId,attr"`
		} `xml:"postCodeRegex"`
	} `xml:"postalCodeData"`
	CalendarData *struct {
		Common
		Calendar []*struct {
			Common
			Territories    string  `xml:"territories,attr"`
			CalendarSystem *Common `xml:"calendarSystem"`
			Eras           *struct {
				Common
				Era []*struct {
					Common
					Start string `xml:"start,attr"`
					End   string `xml:"end,attr"`
				} `xml:"era"`
			} `xml:"eras"`
		} `xml:"calendar"`
	} `xml:"calendarData"`
	CalendarPreferenceData *struct {
		Common
		CalendarPreference []*struct {
			Common
			Territories string `xml:"territories,attr"`
			Ordering    string `xml:"ordering,attr"`
		} `xml:"calendarPreference"`
	} `xml:"calendarPreferenceData"`
	WeekData *struct {
		Common
		MinDays []*struct {
			Common
			Count       string `xml:"count,attr"`
			Territories string `xml:"territories,attr"`
		} `xml:"minDays"`
		FirstDay []*struct {
			Common
			Day         string `xml:"day,attr"`
			Territories string `xml:"territories,attr"`
		} `xml:"firstDay"`
		WeekendStart []*struct {
			Common
			Day         string `xml:"day,attr"`
			Territories string `xml:"territories,attr"`
		} `xml:"weekendStart"`
		WeekendEnd []*struct {
			Common
			Day         string `xml:"day,attr"`
			Territories string `xml:"territories,attr"`
		} `xml:"weekendEnd"`
		WeekOfPreference []*struct {
			Common
			Locales  string `xml:"locales,attr"`
			Ordering string `xml:"ordering,attr"`
		} `xml:"weekOfPreference"`
	} `xml:"weekData"`
	TimeData *struct {
		Common
		Hours []*struct {
			Common
			Allowed   string `xml:"allowed,attr"`
			Preferred string `xml:"preferred,attr"`
			Regions   string `xml:"regions,attr"`
		} `xml:"hours"`
	} `xml:"timeData"`
	MeasurementData *struct {
		Common
		MeasurementSystem []*struct {
			Common
			Category    string `xml:"category,attr"`
			Territories string `xml:"territories,attr"`
		} `xml:"measurementSystem"`
		PaperSize []*struct {
			Common
			Territories string `xml:"territories,attr"`
		} `xml:"paperSize"`
	} `xml:"measurementData"`
	UnitPreferenceData *struct {
		Common
		UnitPreferences []*struct {
			Common
			Category       string `xml:"category,attr"`
			Usage          string `xml:"usage,attr"`
			Scope          string `xml:"scope,attr"`
			UnitPreference []*struct {
				Common
				Regions string `xml:"regions,attr"`
			} `xml:"unitPreference"`
		} `xml:"unitPreferences"`
	} `xml:"unitPreferenceData"`
	TimezoneData *struct {
		Common
		MapTimezones []*struct {
			Common
			OtherVersion string `xml:"otherVersion,attr"`
			TypeVersion  string `xml:"typeVersion,attr"`
			MapZone      []*struct {
				Common
				Other     string `xml:"other,attr"`
				Territory string `xml:"territory,attr"`
			} `xml:"mapZone"`
		} `xml:"mapTimezones"`
		ZoneFormatting []*struct {
			Common
			Multizone   string `xml:"multizone,attr"`
			TzidVersion string `xml:"tzidVersion,attr"`
			ZoneItem    []*struct {
				Common
				Territory string `xml:"territory,attr"`
				Aliases   string `xml:"aliases,attr"`
			} `xml:"zoneItem"`
		} `xml:"zoneFormatting"`
	} `xml:"timezoneData"`
	Characters *struct {
		Common
		CharacterFallback []*struct {
			Common
			Character []*struct {
				Common
				Value      string    `xml:"value,attr"`
				Substitute []*Common `xml:"substitute"`
			} `xml:"character"`
		} `xml:"character-fallback"`
	} `xml:"characters"`
	Transforms *struct {
		Common
		Transform []*struct {
			Common
			Source        string    `xml:"source,attr"`
			Target        string    `xml:"target,attr"`
			Variant       string    `xml:"variant,attr"`
			Direction     string    `xml:"direction,attr"`
			Alias         string    `xml:"alias,attr"`
			BackwardAlias string    `xml:"backwardAlias,attr"`
			Visibility    string    `xml:"visibility,attr"`
			Comment       []*Common `xml:"comment"`
			TRule         []*Common `xml:"tRule"`
		} `xml:"transform"`
	} `xml:"transforms"`
	Metadata *struct {
		Common
		AttributeOrder *Common `xml:"attributeOrder"`
		ElementOrder   *Common `xml:"elementOrder"`
		SerialElements *Common `xml:"serialElements"`
		Suppress       *struct {
			Common
			Attributes []*struct {
				Common
				Element        string `xml:"element,attr"`
				Attribute      string `xml:"attribute,attr"`
				AttributeValue string `xml:"attributeValue,attr"`
			} `xml:"attributes"`
		} `xml:"suppress"`
		Validity *struct {
			Common
			Variable []*struct {
				Common
				Id string `xml:"id,attr"`
			} `xml:"variable"`
			AttributeValues []*struct {
				Common
				Dtds       string `xml:"dtds,attr"`
				Elements   string `xml:"elements,attr"`
				Attributes string `xml:"attributes,attr"`
				Order      string `xml:"order,attr"`
			} `xml:"attributeValues"`
		} `xml:"validity"`
		Alias *struct {
			Common
			LanguageAlias []*struct {
				Common
				Replacement string `xml:"replacement,attr"`
				Reason      string `xml:"reason,attr"`
			} `xml:"languageAlias"`
			ScriptAlias []*struct {
				Common
				Replacement string `xml:"replacement,attr"`
				Reason      string `xml:"reason,attr"`
			} `xml:"scriptAlias"`
			TerritoryAlias []*struct {
				Common
				Replacement string `xml:"replacement,attr"`
				Reason      string `xml:"reason,attr"`
			} `xml:"territoryAlias"`
			SubdivisionAlias []*struct {
				Common
				Replacement string `xml:"replacement,attr"`
				Reason      string `xml:"reason,attr"`
			} `xml:"subdivisionAlias"`
			VariantAlias []*struct {
				Common
				Replacement string `xml:"replacement,attr"`
				Reason      string `xml:"reason,attr"`
			} `xml:"variantAlias"`
			ZoneAlias []*struct {
				Common
				Replacement string `xml:"replacement,attr"`
				Reason      string `xml:"reason,attr"`
			} `xml:"zoneAlias"`
		} `xml:"alias"`
		Deprecated *struct {
			Common
			DeprecatedItems []*struct {
				Common
				Elements   string `xml:"elements,attr"`
				Attributes string `xml:"attributes,attr"`
				Values     string `xml:"values,attr"`
			} `xml:"deprecatedItems"`
		} `xml:"deprecated"`
		Distinguishing *struct {
			Common
			DistinguishingItems []*struct {
				Common
				Exclude    string `xml:"exclude,attr"`
				Elements   string `xml:"elements,attr"`
				Attributes string `xml:"attributes,attr"`
			} `xml:"distinguishingItems"`
		} `xml:"distinguishing"`
		Blocking *struct {
			Common
			BlockingItems []*struct {
				Common
				Elements string `xml:"elements,attr"`
			} `xml:"blockingItems"`
		} `xml:"blocking"`
		CoverageAdditions *struct {
			Common
			LanguageCoverage []*struct {
				Common
				Values string `xml:"values,attr"`
			} `xml:"languageCoverage"`
			ScriptCoverage []*struct {
				Common
				Values string `xml:"values,attr"`
			} `xml:"scriptCoverage"`
			TerritoryCoverage []*struct {
				Common
				Values string `xml:"values,attr"`
			} `xml:"territoryCoverage"`
			CurrencyCoverage []*struct {
				Common
				Values string `xml:"values,attr"`
			} `xml:"currencyCoverage"`
			TimezoneCoverage []*struct {
				Common
				Values string `xml:"values,attr"`
			} `xml:"timezoneCoverage"`
		} `xml:"coverageAdditions"`
		SkipDefaultLocale *struct {
			Common
			Services string `xml:"services,attr"`
		} `xml:"skipDefaultLocale"`
		DefaultContent *struct {
			Common
			Locales string `xml:"locales,attr"`
		} `xml:"defaultContent"`
	} `xml:"metadata"`
	CodeMappings *struct {
		Common
		LanguageCodes []*struct {
			Common
			Alpha3 string `xml:"alpha3,attr"`
		} `xml:"languageCodes"`
		TerritoryCodes []*struct {
			Common
			Numeric  string `xml:"numeric,attr"`
			Alpha3   string `xml:"alpha3,attr"`
			Fips10   string `xml:"fips10,attr"`
			Internet string `xml:"internet,attr"`
		} `xml:"territoryCodes"`
		CurrencyCodes []*struct {
			Common
			Numeric string `xml:"numeric,attr"`
		} `xml:"currencyCodes"`
	} `xml:"codeMappings"`
	ParentLocales *struct {
		Common
		ParentLocale []*struct {
			Common
			Parent  string `xml:"parent,attr"`
			Locales string `xml:"locales,attr"`
		} `xml:"parentLocale"`
	} `xml:"parentLocales"`
	LikelySubtags *struct {
		Common
		LikelySubtag []*struct {
			Common
			From string `xml:"from,attr"`
			To   string `xml:"to,attr"`
		} `xml:"likelySubtag"`
	} `xml:"likelySubtags"`
	MetazoneInfo *struct {
		Common
		Timezone []*struct {
			Common
			UsesMetazone []*struct {
				Common
				From  string `xml:"from,attr"`
				To    string `xml:"to,attr"`
				Mzone string `xml:"mzone,attr"`
			} `xml:"usesMetazone"`
		} `xml:"timezone"`
	} `xml:"metazoneInfo"`
	Plurals []*struct {
		Common
		PluralRules []*struct {
			Common
			Locales    string `xml:"locales,attr"`
			PluralRule []*struct {
				Common
				Count string `xml:"count,attr"`
			} `xml:"pluralRule"`
		} `xml:"pluralRules"`
		PluralRanges []*struct {
			Common
			Locales     string `xml:"locales,attr"`
			PluralRange []*struct {
				Common
				Start  string `xml:"start,attr"`
				End    string `xml:"end,attr"`
				Result string `xml:"result,attr"`
			} `xml:"pluralRange"`
		} `xml:"pluralRanges"`
	} `xml:"plurals"`
	TelephoneCodeData *struct {
		Common
		CodesByTerritory []*struct {
			Common
			Territory            string `xml:"territory,attr"`
			TelephoneCountryCode []*struct {
				Common
				Code string `xml:"code,attr"`
				From string `xml:"from,attr"`
				To   string `xml:"to,attr"`
			} `xml:"telephoneCountryCode"`
		} `xml:"codesByTerritory"`
	} `xml:"telephoneCodeData"`
	NumberingSystems *struct {
		Common
		NumberingSystem []*struct {
			Common
			Id     string `xml:"id,attr"`
			Radix  string `xml:"radix,attr"`
			Digits string `xml:"digits,attr"`
			Rules  string `xml:"rules,attr"`
		} `xml:"numberingSystem"`
	} `xml:"numberingSystems"`
	Bcp47KeywordMappings *struct {
		Common
		MapKeys *struct {
			Common
			KeyMap []*struct {
				Common
				Bcp47 string `xml:"bcp47,attr"`
			} `xml:"keyMap"`
		} `xml:"mapKeys"`
		MapTypes []*struct {
			Common
			TypeMap []*struct {
				Common
				Bcp47 string `xml:"bcp47,attr"`
			} `xml:"typeMap"`
		} `xml:"mapTypes"`
	} `xml:"bcp47KeywordMappings"`
	Gender *struct {
		Common
		PersonList []*struct {
			Common
			Locales string `xml:"locales,attr"`
		} `xml:"personList"`
	} `xml:"gender"`
	References *struct {
		Common
		Reference []*struct {
			Common
			Uri string `xml:"uri,attr"`
		} `xml:"reference"`
	} `xml:"references"`
	LanguageMatching *struct {
		Common
		LanguageMatches []*struct {
			Common
			ParadigmLocales []*struct {
				Common
				Locales string `xml:"locales,attr"`
			} `xml:"paradigmLocales"`
			MatchVariable []*struct {
				Common
				Id    string `xml:"id,attr"`
				Value string `xml:"value,attr"`
			} `xml:"matchVariable"`
			LanguageMatch []*struct {
				Common
				Desired   string `xml:"desired,attr"`
				Supported string `xml:"supported,attr"`
				Percent   string `xml:"percent,attr"`
				Distance  string `xml:"distance,attr"`
				Oneway    string `xml:"oneway,attr"`
			} `xml:"languageMatch"`
		} `xml:"languageMatches"`
	} `xml:"languageMatching"`
	DayPeriodRuleSet []*struct {
		Common
		DayPeriodRules []*struct {
			Common
			Locales       string `xml:"locales,attr"`
			DayPeriodRule []*struct {
				Common
				At     string `xml:"at,attr"`
				After  string `xml:"after,attr"`
				Before string `xml:"before,attr"`
				From   string `xml:"from,attr"`
				To     string `xml:"to,attr"`
			} `xml:"dayPeriodRule"`
		} `xml:"dayPeriodRules"`
	} `xml:"dayPeriodRuleSet"`
	MetaZones *struct {
		Common
		MetazoneInfo *struct {
			Common
			Timezone []*struct {
				Common
				UsesMetazone []*struct {
					Common
					From  string `xml:"from,attr"`
					To    string `xml:"to,attr"`
					Mzone string `xml:"mzone,attr"`
				} `xml:"usesMetazone"`
			} `xml:"timezone"`
		} `xml:"metazoneInfo"`
		MapTimezones *struct {
			Common
			OtherVersion string `xml:"otherVersion,attr"`
			TypeVersion  string `xml:"typeVersion,attr"`
			MapZone      []*struct {
				Common
				Other     string `xml:"other,attr"`
				Territory string `xml:"territory,attr"`
			} `xml:"mapZone"`
		} `xml:"mapTimezones"`
	} `xml:"metaZones"`
	PrimaryZones *struct {
		Common
		PrimaryZone []*struct {
			Common
			Iso3166 string `xml:"iso3166,attr"`
		} `xml:"primaryZone"`
	} `xml:"primaryZones"`
	WindowsZones *struct {
		Common
		MapTimezones *struct {
			Common
			OtherVersion string `xml:"otherVersion,attr"`
			TypeVersion  string `xml:"typeVersion,attr"`
			MapZone      []*struct {
				Common
				Other     string `xml:"other,attr"`
				Territory string `xml:"territory,attr"`
			} `xml:"mapZone"`
		} `xml:"mapTimezones"`
	} `xml:"windowsZones"`
	CoverageLevels *struct {
		Common
		ApprovalRequirements *struct {
			Common
			ApprovalRequirement []*struct {
				Common
				Votes   string `xml:"votes,attr"`
				Locales string `xml:"locales,attr"`
				Paths   string `xml:"paths,attr"`
			} `xml:"approvalRequirement"`
		} `xml:"approvalRequirements"`
		CoverageVariable []*struct {
			Common
			Key   string `xml:"key,attr"`
			Value string `xml:"value,attr"`
		} `xml:"coverageVariable"`
		CoverageLevel []*struct {
			Common
			InLanguage  string `xml:"inLanguage,attr"`
			InScript    string `xml:"inScript,attr"`
			InTerritory string `xml:"inTerritory,attr"`
			Value       string `xml:"value,attr"`
			Match       string `xml:"match,attr"`
		} `xml:"coverageLevel"`
	} `xml:"coverageLevels"`
	IdValidity *struct {
		Common
		Id []*struct {
			Common
			IdStatus string `xml:"idStatus,attr"`
		} `xml:"id"`
	} `xml:"idValidity"`
	RgScope *struct {
		Common
		RgPath []*struct {
			Common
			Path string `xml:"path,attr"`
		} `xml:"rgPath"`
	} `xml:"rgScope"`
	LanguageGroups *struct {
		Common
		LanguageGroup []*struct {
			Common
			Parent string `xml:"parent,attr"`
		} `xml:"languageGroup"`
	} `xml:"languageGroups"`
}


type LDML struct {
	Common
	Version  string `xml:"version,attr"`
	Identity *struct {
		Common
		Version *struct {
			Common
			Number string `xml:"number,attr"`
		} `xml:"version"`
		Generation *struct {
			Common
			Date string `xml:"date,attr"`
		} `xml:"generation"`
		Language  *Common `xml:"language"`
		Script    *Common `xml:"script"`
		Territory *Common `xml:"territory"`
		Variant   *Common `xml:"variant"`
	} `xml:"identity"`
	LocaleDisplayNames *LocaleDisplayNames `xml:"localeDisplayNames"`
	Layout             *struct {
		Common
		Orientation []*struct {
			Common
			Characters     string    `xml:"characters,attr"`
			Lines          string    `xml:"lines,attr"`
			CharacterOrder []*Common `xml:"characterOrder"`
			LineOrder      []*Common `xml:"lineOrder"`
		} `xml:"orientation"`
		InList []*struct {
			Common
			Casing string `xml:"casing,attr"`
		} `xml:"inList"`
		InText []*Common `xml:"inText"`
	} `xml:"layout"`
	ContextTransforms *struct {
		Common
		ContextTransformUsage []*struct {
			Common
			ContextTransform []*Common `xml:"contextTransform"`
		} `xml:"contextTransformUsage"`
	} `xml:"contextTransforms"`
	Characters *struct {
		Common
		ExemplarCharacters []*Common `xml:"exemplarCharacters"`
		Ellipsis           []*Common `xml:"ellipsis"`
		MoreInformation    []*Common `xml:"moreInformation"`
		Stopwords          []*struct {
			Common
			StopwordList []*Common `xml:"stopwordList"`
		} `xml:"stopwords"`
		IndexLabels []*struct {
			Common
			IndexSeparator           []*Common `xml:"indexSeparator"`
			CompressedIndexSeparator []*Common `xml:"compressedIndexSeparator"`
			IndexRangePattern        []*Common `xml:"indexRangePattern"`
			IndexLabelBefore         []*Common `xml:"indexLabelBefore"`
			IndexLabelAfter          []*Common `xml:"indexLabelAfter"`
			IndexLabel               []*struct {
				Common
				IndexSource string `xml:"indexSource,attr"`
				Priority    string `xml:"priority,attr"`
			} `xml:"indexLabel"`
		} `xml:"indexLabels"`
		Mapping []*struct {
			Common
			Registry string `xml:"registry,attr"`
		} `xml:"mapping"`
		ParseLenients []*struct {
			Common
			Scope        string `xml:"scope,attr"`
			Level        string `xml:"level,attr"`
			ParseLenient []*struct {
				Common
				Sample string `xml:"sample,attr"`
			} `xml:"parseLenient"`
		} `xml:"parseLenients"`
	} `xml:"characters"`
	Delimiters *struct {
		Common
		QuotationStart          []*Common `xml:"quotationStart"`
		QuotationEnd            []*Common `xml:"quotationEnd"`
		AlternateQuotationStart []*Common `xml:"alternateQuotationStart"`
		AlternateQuotationEnd   []*Common `xml:"alternateQuotationEnd"`
	} `xml:"delimiters"`
	Measurement *struct {
		Common
		MeasurementSystem []*Common `xml:"measurementSystem"`
		PaperSize         []*struct {
			Common
			Height []*Common `xml:"height"`
			Width  []*Common `xml:"width"`
		} `xml:"paperSize"`
	} `xml:"measurement"`
	Dates *struct {
		Common
		LocalizedPatternChars []*Common `xml:"localizedPatternChars"`
		DateRangePattern      []*Common `xml:"dateRangePattern"`
		Calendars             *struct {
			Common
			Calendar []*Calendar `xml:"calendar"`
		} `xml:"calendars"`
		Fields *struct {
			Common
			Field []*struct {
				Common
				DisplayName []*struct {
					Common
					Count string `xml:"count,attr"`
				} `xml:"displayName"`
				Relative     []*Common `xml:"relative"`
				RelativeTime []*struct {
					Common
					RelativeTimePattern []*struct {
						Common
						Count string `xml:"count,attr"`
					} `xml:"relativeTimePattern"`
				} `xml:"relativeTime"`
				RelativePeriod []*Common `xml:"relativePeriod"`
			} `xml:"field"`
		} `xml:"fields"`
		TimeZoneNames *TimeZoneNames `xml:"timeZoneNames"`
	} `xml:"dates"`
	Numbers *Numbers `xml:"numbers"`
	Units   *struct {
		Common
		Unit []*struct {
			Common
			DisplayName []*struct {
				Common
				Count string `xml:"count,attr"`
			} `xml:"displayName"`
			UnitPattern []*struct {
				Common
				Count string `xml:"count,attr"`
			} `xml:"unitPattern"`
			PerUnitPattern []*Common `xml:"perUnitPattern"`
		} `xml:"unit"`
		UnitLength []*struct {
			Common
			CompoundUnit []*struct {
				Common
				CompoundUnitPattern []*Common `xml:"compoundUnitPattern"`
			} `xml:"compoundUnit"`
			Unit []*struct {
				Common
				DisplayName []*struct {
					Common
					Count string `xml:"count,attr"`
				} `xml:"displayName"`
				UnitPattern []*struct {
					Common
					Count string `xml:"count,attr"`
				} `xml:"unitPattern"`
				PerUnitPattern []*Common `xml:"perUnitPattern"`
			} `xml:"unit"`
			CoordinateUnit []*struct {
				Common
				CoordinateUnitPattern []*Common `xml:"coordinateUnitPattern"`
			} `xml:"coordinateUnit"`
		} `xml:"unitLength"`
		DurationUnit []*struct {
			Common
			DurationUnitPattern []*Common `xml:"durationUnitPattern"`
		} `xml:"durationUnit"`
	} `xml:"units"`
	ListPatterns *struct {
		Common
		ListPattern []*struct {
			Common
			ListPatternPart []*Common `xml:"listPatternPart"`
		} `xml:"listPattern"`
	} `xml:"listPatterns"`
	Collations *struct {
		Common
		Version          string       `xml:"version,attr"`
		DefaultCollation *Common      `xml:"defaultCollation"`
		Collation        []*Collation `xml:"collation"`
	} `xml:"collations"`
	Posix *struct {
		Common
		Messages []*struct {
			Common
			Yesstr  []*Common `xml:"yesstr"`
			Nostr   []*Common `xml:"nostr"`
			Yesexpr []*Common `xml:"yesexpr"`
			Noexpr  []*Common `xml:"noexpr"`
		} `xml:"messages"`
	} `xml:"posix"`
	CharacterLabels *struct {
		Common
		CharacterLabelPattern []*struct {
			Common
			Count string `xml:"count,attr"`
		} `xml:"characterLabelPattern"`
		CharacterLabel []*Common `xml:"characterLabel"`
	} `xml:"characterLabels"`
	Segmentations *struct {
		Common
		Segmentation []*struct {
			Common
			Variables *struct {
				Common
				Variable []*struct {
					Common
					Id string `xml:"id,attr"`
				} `xml:"variable"`
			} `xml:"variables"`
			SegmentRules *struct {
				Common
				Rule []*struct {
					Common
					Id string `xml:"id,attr"`
				} `xml:"rule"`
			} `xml:"segmentRules"`
			Exceptions *struct {
				Common
				Exception []*Common `xml:"exception"`
			} `xml:"exceptions"`
			Suppressions *struct {
				Common
				Suppression []*Common `xml:"suppression"`
			} `xml:"suppressions"`
		} `xml:"segmentation"`
	} `xml:"segmentations"`
	Rbnf *struct {
		Common
		RulesetGrouping []*struct {
			Common
			Ruleset []*struct {
				Common
				Access        string `xml:"access,attr"`
				AllowsParsing string `xml:"allowsParsing,attr"`
				Rbnfrule      []*struct {
					Common
					Value  string `xml:"value,attr"`
					Radix  string `xml:"radix,attr"`
					Decexp string `xml:"decexp,attr"`
				} `xml:"rbnfrule"`
			} `xml:"ruleset"`
		} `xml:"rulesetGrouping"`
	} `xml:"rbnf"`
	Annotations *struct {
		Common
		Annotation []*struct {
			Common
			Cp  string `xml:"cp,attr"`
			Tts string `xml:"tts,attr"`
		} `xml:"annotation"`
	} `xml:"annotations"`
	Metadata *struct {
		Common
		CasingData *struct {
			Common
			CasingItem []*struct {
				Common
				Override   string `xml:"override,attr"`
				ForceError string `xml:"forceError,attr"`
			} `xml:"casingItem"`
		} `xml:"casingData"`
	} `xml:"metadata"`
	References *struct {
		Common
		Reference []*struct {
			Common
			Uri string `xml:"uri,attr"`
		} `xml:"reference"`
	} `xml:"references"`
}





type Collation struct {
	Common
	Visibility string  `xml:"visibility,attr"`
	Base       *Common `xml:"base"`
	Import     []*struct {
		Common
		Source string `xml:"source,attr"`
	} `xml:"import"`
	Settings *struct {
		Common
		Strength           string `xml:"strength,attr"`
		Alternate          string `xml:"alternate,attr"`
		Backwards          string `xml:"backwards,attr"`
		Normalization      string `xml:"normalization,attr"`
		CaseLevel          string `xml:"caseLevel,attr"`
		CaseFirst          string `xml:"caseFirst,attr"`
		HiraganaQuaternary string `xml:"hiraganaQuaternary,attr"`
		MaxVariable        string `xml:"maxVariable,attr"`
		Numeric            string `xml:"numeric,attr"`
		Private            string `xml:"private,attr"`
		VariableTop        string `xml:"variableTop,attr"`
		Reorder            string `xml:"reorder,attr"`
	} `xml:"settings"`
	SuppressContractions *Common   `xml:"suppress_contractions"`
	Optimize             *Common   `xml:"optimize"`
	Cr                   []*Common `xml:"cr"`
	rulesElem
}





type Calendar struct {
	Common
	Months *struct {
		Common
		MonthContext []*struct {
			Common
			MonthWidth []*struct {
				Common
				Month []*struct {
					Common
					Yeartype string `xml:"yeartype,attr"`
				} `xml:"month"`
			} `xml:"monthWidth"`
		} `xml:"monthContext"`
	} `xml:"months"`
	MonthNames *struct {
		Common
		Month []*struct {
			Common
			Yeartype string `xml:"yeartype,attr"`
		} `xml:"month"`
	} `xml:"monthNames"`
	MonthAbbr *struct {
		Common
		Month []*struct {
			Common
			Yeartype string `xml:"yeartype,attr"`
		} `xml:"month"`
	} `xml:"monthAbbr"`
	MonthPatterns *struct {
		Common
		MonthPatternContext []*struct {
			Common
			MonthPatternWidth []*struct {
				Common
				MonthPattern []*Common `xml:"monthPattern"`
			} `xml:"monthPatternWidth"`
		} `xml:"monthPatternContext"`
	} `xml:"monthPatterns"`
	Days *struct {
		Common
		DayContext []*struct {
			Common
			DayWidth []*struct {
				Common
				Day []*Common `xml:"day"`
			} `xml:"dayWidth"`
		} `xml:"dayContext"`
	} `xml:"days"`
	DayNames *struct {
		Common
		Day []*Common `xml:"day"`
	} `xml:"dayNames"`
	DayAbbr *struct {
		Common
		Day []*Common `xml:"day"`
	} `xml:"dayAbbr"`
	Quarters *struct {
		Common
		QuarterContext []*struct {
			Common
			QuarterWidth []*struct {
				Common
				Quarter []*Common `xml:"quarter"`
			} `xml:"quarterWidth"`
		} `xml:"quarterContext"`
	} `xml:"quarters"`
	Week *struct {
		Common
		MinDays []*struct {
			Common
			Count string `xml:"count,attr"`
		} `xml:"minDays"`
		FirstDay []*struct {
			Common
			Day string `xml:"day,attr"`
		} `xml:"firstDay"`
		WeekendStart []*struct {
			Common
			Day  string `xml:"day,attr"`
			Time string `xml:"time,attr"`
		} `xml:"weekendStart"`
		WeekendEnd []*struct {
			Common
			Day  string `xml:"day,attr"`
			Time string `xml:"time,attr"`
		} `xml:"weekendEnd"`
	} `xml:"week"`
	Am         []*Common `xml:"am"`
	Pm         []*Common `xml:"pm"`
	DayPeriods *struct {
		Common
		DayPeriodContext []*struct {
			Common
			DayPeriodWidth []*struct {
				Common
				DayPeriod []*Common `xml:"dayPeriod"`
			} `xml:"dayPeriodWidth"`
		} `xml:"dayPeriodContext"`
	} `xml:"dayPeriods"`
	Eras *struct {
		Common
		EraNames *struct {
			Common
			Era []*Common `xml:"era"`
		} `xml:"eraNames"`
		EraAbbr *struct {
			Common
			Era []*Common `xml:"era"`
		} `xml:"eraAbbr"`
		EraNarrow *struct {
			Common
			Era []*Common `xml:"era"`
		} `xml:"eraNarrow"`
	} `xml:"eras"`
	CyclicNameSets *struct {
		Common
		CyclicNameSet []*struct {
			Common
			CyclicNameContext []*struct {
				Common
				CyclicNameWidth []*struct {
					Common
					CyclicName []*Common `xml:"cyclicName"`
				} `xml:"cyclicNameWidth"`
			} `xml:"cyclicNameContext"`
		} `xml:"cyclicNameSet"`
	} `xml:"cyclicNameSets"`
	DateFormats *struct {
		Common
		DateFormatLength []*struct {
			Common
			DateFormat []*struct {
				Common
				Pattern []*struct {
					Common
					Numbers string `xml:"numbers,attr"`
					Count   string `xml:"count,attr"`
				} `xml:"pattern"`
				DisplayName []*struct {
					Common
					Count string `xml:"count,attr"`
				} `xml:"displayName"`
			} `xml:"dateFormat"`
		} `xml:"dateFormatLength"`
	} `xml:"dateFormats"`
	TimeFormats *struct {
		Common
		TimeFormatLength []*struct {
			Common
			TimeFormat []*struct {
				Common
				Pattern []*struct {
					Common
					Numbers string `xml:"numbers,attr"`
					Count   string `xml:"count,attr"`
				} `xml:"pattern"`
				DisplayName []*struct {
					Common
					Count string `xml:"count,attr"`
				} `xml:"displayName"`
			} `xml:"timeFormat"`
		} `xml:"timeFormatLength"`
	} `xml:"timeFormats"`
	DateTimeFormats *struct {
		Common
		DateTimeFormatLength []*struct {
			Common
			DateTimeFormat []*struct {
				Common
				Pattern []*struct {
					Common
					Numbers string `xml:"numbers,attr"`
					Count   string `xml:"count,attr"`
				} `xml:"pattern"`
				DisplayName []*struct {
					Common
					Count string `xml:"count,attr"`
				} `xml:"displayName"`
			} `xml:"dateTimeFormat"`
		} `xml:"dateTimeFormatLength"`
		AvailableFormats []*struct {
			Common
			DateFormatItem []*struct {
				Common
				Id    string `xml:"id,attr"`
				Count string `xml:"count,attr"`
			} `xml:"dateFormatItem"`
		} `xml:"availableFormats"`
		AppendItems []*struct {
			Common
			AppendItem []*struct {
				Common
				Request string `xml:"request,attr"`
			} `xml:"appendItem"`
		} `xml:"appendItems"`
		IntervalFormats []*struct {
			Common
			IntervalFormatFallback []*Common `xml:"intervalFormatFallback"`
			IntervalFormatItem     []*struct {
				Common
				Id                 string `xml:"id,attr"`
				GreatestDifference []*struct {
					Common
					Id string `xml:"id,attr"`
				} `xml:"greatestDifference"`
			} `xml:"intervalFormatItem"`
		} `xml:"intervalFormats"`
	} `xml:"dateTimeFormats"`
	Fields []*struct {
		Common
		Field []*struct {
			Common
			DisplayName []*struct {
				Common
				Count string `xml:"count,attr"`
			} `xml:"displayName"`
			Relative     []*Common `xml:"relative"`
			RelativeTime []*struct {
				Common
				RelativeTimePattern []*struct {
					Common
					Count string `xml:"count,attr"`
				} `xml:"relativeTimePattern"`
			} `xml:"relativeTime"`
			RelativePeriod []*Common `xml:"relativePeriod"`
		} `xml:"field"`
	} `xml:"fields"`
}
type TimeZoneNames struct {
	Common
	HourFormat           []*Common `xml:"hourFormat"`
	HoursFormat          []*Common `xml:"hoursFormat"`
	GmtFormat            []*Common `xml:"gmtFormat"`
	GmtZeroFormat        []*Common `xml:"gmtZeroFormat"`
	RegionFormat         []*Common `xml:"regionFormat"`
	FallbackFormat       []*Common `xml:"fallbackFormat"`
	FallbackRegionFormat []*Common `xml:"fallbackRegionFormat"`
	AbbreviationFallback []*Common `xml:"abbreviationFallback"`
	PreferenceOrdering   []*Common `xml:"preferenceOrdering"`
	SingleCountries      []*struct {
		Common
		List string `xml:"list,attr"`
	} `xml:"singleCountries"`
	Zone []*struct {
		Common
		Long []*struct {
			Common
			Generic  []*Common `xml:"generic"`
			Standard []*Common `xml:"standard"`
			Daylight []*Common `xml:"daylight"`
		} `xml:"long"`
		Short []*struct {
			Common
			Generic  []*Common `xml:"generic"`
			Standard []*Common `xml:"standard"`
			Daylight []*Common `xml:"daylight"`
		} `xml:"short"`
		CommonlyUsed []*struct {
			Common
			Used string `xml:"used,attr"`
		} `xml:"commonlyUsed"`
		ExemplarCity []*Common `xml:"exemplarCity"`
	} `xml:"zone"`
	Metazone []*struct {
		Common
		Long []*struct {
			Common
			Generic  []*Common `xml:"generic"`
			Standard []*Common `xml:"standard"`
			Daylight []*Common `xml:"daylight"`
		} `xml:"long"`
		Short []*struct {
			Common
			Generic  []*Common `xml:"generic"`
			Standard []*Common `xml:"standard"`
			Daylight []*Common `xml:"daylight"`
		} `xml:"short"`
		CommonlyUsed []*struct {
			Common
			Used string `xml:"used,attr"`
		} `xml:"commonlyUsed"`
	} `xml:"metazone"`
}



type LocaleDisplayNames struct {
	Common
	LocaleDisplayPattern *struct {
		Common
		LocalePattern        []*Common `xml:"localePattern"`
		LocaleSeparator      []*Common `xml:"localeSeparator"`
		LocaleKeyTypePattern []*Common `xml:"localeKeyTypePattern"`
	} `xml:"localeDisplayPattern"`
	Languages *struct {
		Common
		Language []*Common `xml:"language"`
	} `xml:"languages"`
	Scripts *struct {
		Common
		Script []*Common `xml:"script"`
	} `xml:"scripts"`
	Territories *struct {
		Common
		Territory []*Common `xml:"territory"`
	} `xml:"territories"`
	Subdivisions *struct {
		Common
		Subdivision []*Common `xml:"subdivision"`
	} `xml:"subdivisions"`
	Variants *struct {
		Common
		Variant []*Common `xml:"variant"`
	} `xml:"variants"`
	Keys *struct {
		Common
		Key []*Common `xml:"key"`
	} `xml:"keys"`
	Types *struct {
		Common
		Type []*struct {
			Common
			Key string `xml:"key,attr"`
		} `xml:"type"`
	} `xml:"types"`
	TransformNames *struct {
		Common
		TransformName []*Common `xml:"transformName"`
	} `xml:"transformNames"`
	MeasurementSystemNames *struct {
		Common
		MeasurementSystemName []*Common `xml:"measurementSystemName"`
	} `xml:"measurementSystemNames"`
	CodePatterns *struct {
		Common
		CodePattern []*Common `xml:"codePattern"`
	} `xml:"codePatterns"`
}


type Numbers struct {
	Common
	DefaultNumberingSystem []*Common `xml:"defaultNumberingSystem"`
	OtherNumberingSystems  []*struct {
		Common
		Native      []*Common `xml:"native"`
		Traditional []*Common `xml:"traditional"`
		Finance     []*Common `xml:"finance"`
	} `xml:"otherNumberingSystems"`
	MinimumGroupingDigits []*Common `xml:"minimumGroupingDigits"`
	Symbols               []*struct {
		Common
		NumberSystem string `xml:"numberSystem,attr"`
		Decimal      []*struct {
			Common
			NumberSystem string `xml:"numberSystem,attr"`
		} `xml:"decimal"`
		Group []*struct {
			Common
			NumberSystem string `xml:"numberSystem,attr"`
		} `xml:"group"`
		List []*struct {
			Common
			NumberSystem string `xml:"numberSystem,attr"`
		} `xml:"list"`
		PercentSign []*struct {
			Common
			NumberSystem string `xml:"numberSystem,attr"`
		} `xml:"percentSign"`
		NativeZeroDigit []*struct {
			Common
			NumberSystem string `xml:"numberSystem,attr"`
		} `xml:"nativeZeroDigit"`
		PatternDigit []*struct {
			Common
			NumberSystem string `xml:"numberSystem,attr"`
		} `xml:"patternDigit"`
		PlusSign []*struct {
			Common
			NumberSystem string `xml:"numberSystem,attr"`
		} `xml:"plusSign"`
		MinusSign []*struct {
			Common
			NumberSystem string `xml:"numberSystem,attr"`
		} `xml:"minusSign"`
		Exponential []*struct {
			Common
			NumberSystem string `xml:"numberSystem,attr"`
		} `xml:"exponential"`
		SuperscriptingExponent []*Common `xml:"superscriptingExponent"`
		PerMille               []*struct {
			Common
			NumberSystem string `xml:"numberSystem,attr"`
		} `xml:"perMille"`
		Infinity []*struct {
			Common
			NumberSystem string `xml:"numberSystem,attr"`
		} `xml:"infinity"`
		Nan []*struct {
			Common
			NumberSystem string `xml:"numberSystem,attr"`
		} `xml:"nan"`
		CurrencyDecimal []*struct {
			Common
			NumberSystem string `xml:"numberSystem,attr"`
		} `xml:"currencyDecimal"`
		CurrencyGroup []*struct {
			Common
			NumberSystem string `xml:"numberSystem,attr"`
		} `xml:"currencyGroup"`
		TimeSeparator []*Common `xml:"timeSeparator"`
	} `xml:"symbols"`
	DecimalFormats []*struct {
		Common
		NumberSystem        string `xml:"numberSystem,attr"`
		DecimalFormatLength []*struct {
			Common
			DecimalFormat []*struct {
				Common
				Pattern []*struct {
					Common
					Numbers string `xml:"numbers,attr"`
					Count   string `xml:"count,attr"`
				} `xml:"pattern"`
			} `xml:"decimalFormat"`
		} `xml:"decimalFormatLength"`
	} `xml:"decimalFormats"`
	ScientificFormats []*struct {
		Common
		NumberSystem           string `xml:"numberSystem,attr"`
		ScientificFormatLength []*struct {
			Common
			ScientificFormat []*struct {
				Common
				Pattern []*struct {
					Common
					Numbers string `xml:"numbers,attr"`
					Count   string `xml:"count,attr"`
				} `xml:"pattern"`
			} `xml:"scientificFormat"`
		} `xml:"scientificFormatLength"`
	} `xml:"scientificFormats"`
	PercentFormats []*struct {
		Common
		NumberSystem        string `xml:"numberSystem,attr"`
		PercentFormatLength []*struct {
			Common
			PercentFormat []*struct {
				Common
				Pattern []*struct {
					Common
					Numbers string `xml:"numbers,attr"`
					Count   string `xml:"count,attr"`
				} `xml:"pattern"`
			} `xml:"percentFormat"`
		} `xml:"percentFormatLength"`
	} `xml:"percentFormats"`
	CurrencyFormats []*struct {
		Common
		NumberSystem    string `xml:"numberSystem,attr"`
		CurrencySpacing []*struct {
			Common
			BeforeCurrency []*struct {
				Common
				CurrencyMatch    []*Common `xml:"currencyMatch"`
				SurroundingMatch []*Common `xml:"surroundingMatch"`
				InsertBetween    []*Common `xml:"insertBetween"`
			} `xml:"beforeCurrency"`
			AfterCurrency []*struct {
				Common
				CurrencyMatch    []*Common `xml:"currencyMatch"`
				SurroundingMatch []*Common `xml:"surroundingMatch"`
				InsertBetween    []*Common `xml:"insertBetween"`
			} `xml:"afterCurrency"`
		} `xml:"currencySpacing"`
		CurrencyFormatLength []*struct {
			Common
			CurrencyFormat []*struct {
				Common
				Pattern []*struct {
					Common
					Numbers string `xml:"numbers,attr"`
					Count   string `xml:"count,attr"`
				} `xml:"pattern"`
			} `xml:"currencyFormat"`
		} `xml:"currencyFormatLength"`
		UnitPattern []*struct {
			Common
			Count string `xml:"count,attr"`
		} `xml:"unitPattern"`
	} `xml:"currencyFormats"`
	Currencies *struct {
		Common
		Currency []*struct {
			Common
			Pattern []*struct {
				Common
				Numbers string `xml:"numbers,attr"`
				Count   string `xml:"count,attr"`
			} `xml:"pattern"`
			DisplayName []*struct {
				Common
				Count string `xml:"count,attr"`
			} `xml:"displayName"`
			Symbol  []*Common `xml:"symbol"`
			Decimal []*struct {
				Common
				NumberSystem string `xml:"numberSystem,attr"`
			} `xml:"decimal"`
			Group []*struct {
				Common
				NumberSystem string `xml:"numberSystem,attr"`
			} `xml:"group"`
		} `xml:"currency"`
	} `xml:"currencies"`
	MiscPatterns []*struct {
		Common
		NumberSystem string `xml:"numberSystem,attr"`
		Pattern      []*struct {
			Common
			Numbers string `xml:"numbers,attr"`
			Count   string `xml:"count,attr"`
		} `xml:"pattern"`
	} `xml:"miscPatterns"`
	MinimalPairs []*struct {
		Common
		PluralMinimalPairs []*struct {
			Common
			Count string `xml:"count,attr"`
		} `xml:"pluralMinimalPairs"`
		OrdinalMinimalPairs []*struct {
			Common
			Ordinal string `xml:"ordinal,attr"`
		} `xml:"ordinalMinimalPairs"`
	} `xml:"minimalPairs"`
}


const Version = "32"
