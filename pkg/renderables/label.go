package renderables

import (
	"strings"

	"joshuanario.com/hyperfacts/pkg/attr"
	"joshuanario.com/hyperfacts/pkg/hydratables"
)

func GetLabel(h *hydratables.Hydratable, href string) LabelPack {
	ret := LabelPack{}
	ret[Default] = make(LanguagePack)
	ret[Default][PureLabel] = href
	index := strings.IndexRune(href, '#')
	if index > -1 {
		ret[Default][BriefLabel] = href[index:]
	}
	ret = appendLabelModifiersFromHref(ret, h, href)
	return ret
}

func appendLabelModifiersFromHref(labelPack LabelPack, h *hydratables.Hydratable, href string) LabelPack {
	ret := labelPack
	for _, labels := range h.LabelLinkbases {
		for _, labelLink := range labels.LabelLink {
			for _, loc := range labelLink.Locs {
				if loc.Href == href {
					labelArcs := labelLink.LabelArcs
					for _, labelArc := range labelArcs {
						if labelArc.From == loc.Label && labelArc.Arcrole == attr.LabelArcrole {
							for _, labelLinkLabel := range labelLink.Labels {
								if labelLinkLabel.Label == labelArc.To {
									charData := labelLinkLabel.CharData
									switch labelLinkLabel.Role {
									case attr.Label:
										if _, found := ret[Default]; !found {
											ret[Default] = make(LanguagePack)
										}
										ret = appendLanguage(&labelLinkLabel, Default, charData, ret)
									case attr.TerseLabel:
										if _, found := ret[Terse]; !found {
											ret[Terse] = make(LanguagePack)
										}
										ret = appendLanguage(&labelLinkLabel, Terse, charData, ret)
									case attr.VerboseLabel:
										if _, found := ret[Verbose]; !found {
											ret[Verbose] = make(LanguagePack)
										}
										ret = appendLanguage(&labelLinkLabel, Verbose, charData, ret)
									default: //noop
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return ret
}

func appendLanguage(labelLinkLabel *hydratables.LabelLinkLabel, labelRole LabelRole, charData string, labelPack LabelPack) LabelPack {
	switch labelLinkLabel.Lang {
	case "en":
	case "en-US":
		labelPack[labelRole][English] = charData
	case "es":
		labelPack[labelRole][Español] = charData
	default: //noop
	}
	return labelPack
}
