package hydratables

import (
	"joshuanario.com/hyperfacts/pkg/serializables"
)

type Hydratable struct {
	Folder *serializables.Folder
	UnitTypeRegistry
	Documents             []*Document
	Instances             map[string]Instance
	Schemas               map[string]Schema
	LabelLinkbases        map[string]LabelLinkbase
	PresentationLinkbases map[string]PresentationLinkbase
	DefinitionLinkbases   map[string]DefinitionLinkbase
	CalculationLinkbases  map[string]CalculationLinkbase
}

func Hydrate(folder *serializables.Folder) (*Hydratable, error) {
	ret := &Hydratable{
		Folder:                folder,
		Instances:             make(map[string]Instance),
		Schemas:               make(map[string]Schema),
		LabelLinkbases:        make(map[string]LabelLinkbase),
		PresentationLinkbases: make(map[string]PresentationLinkbase),
		DefinitionLinkbases:   make(map[string]DefinitionLinkbase),
		CalculationLinkbases:  make(map[string]CalculationLinkbase),
	}
	for filename, file := range folder.Schemas {
		entry, err := hydrateSchema(&file, filename)
		if err != nil {
			return nil, err
		}
		ret.Schemas[filename] = *entry
		embeddedPresentation, err := hydrateEmbeddedPresentationLinkbase(&file, filename)
		if err == nil && embeddedPresentation != nil {
			ret.PresentationLinkbases[filename] = *embeddedPresentation
		}
		embeddedDefinition, err := hydrateEmbeddedDefinitionLinkbase(&file, filename)
		if err == nil && embeddedDefinition != nil {
			ret.DefinitionLinkbases[filename] = *embeddedDefinition
		}
		embeddedCalculation, err := hydrateEmbeddedCalculationLinkbase(&file, filename)
		if err == nil && embeddedCalculation != nil {
			ret.CalculationLinkbases[filename] = *embeddedCalculation
		}
		embeddedLabel, err := hydrateEmbeddedLabelLinkbase(&file, filename)
		if err == nil && embeddedLabel != nil {
			ret.LabelLinkbases[filename] = *embeddedLabel
		}
	}
	for filename, file := range folder.PresentationLinkbases {
		entry, err := hydratePresentationLinkbase(&file, filename)
		if err != nil {
			return nil, err
		}
		ret.PresentationLinkbases[filename] = *entry
	}
	for filename, file := range folder.DefinitionLinkbases {
		entry, err := hydrateDefinitionLinkbase(&file, filename)
		if err != nil {
			return nil, err
		}
		ret.DefinitionLinkbases[filename] = *entry
	}
	for filename, file := range folder.CalculationLinkbases {
		entry, err := hydrateCalculationLinkbase(&file, filename)
		if err != nil {
			return nil, err
		}
		ret.CalculationLinkbases[filename] = *entry
	}
	for filename, file := range folder.LabelLinkbases {
		entry, err := hydrateLabelLinkbase(&file, filename)
		if err != nil {
			return nil, err
		}
		ret.LabelLinkbases[filename] = *entry
	}
	for filename, file := range folder.Instances {
		entry, err := hydrateInstance(&file, filename, ret)
		if err != nil {
			return nil, err
		}
		ret.Instances[filename] = *entry
	}
	utr, err := HydrateUnitTypeRegistry()
	if err != nil {
		return nil, err
	}
	ret.UnitTypeRegistry = *utr
	docs, err := hydrateDocuments(folder)
	if err != nil {
		return nil, err
	}
	ret.Documents = docs
	return ret, nil
}
