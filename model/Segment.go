package model

import "time"

type Segment struct {
	ContentDisposition   ContentDisposition               `json:"contentDisposition"`
	Entity               map[string]interface{}           `json:"entity"`
	Headers              map[string][]string              `json:"headers"`
	MediaType            MediaType                        `json:"mediaType"`
	MessageBodyWorkers   map[string]interface{}           `json:"messageBodyWorkers"`
	Parent               Parent                           `json:"parent"`
	Providers            map[string]interface{}           `json:"providers"`
	BodyParts            []BodyPart                       `json:"bodyParts"`
	Fields               map[string][]Field               `json:"fields"`
	ParameterizedHeaders map[string][]ParameterizedHeader `json:"parameterizedHeaders"`
}

// ContentDisposition represents the structure for content disposition.
type ContentDisposition struct {
	Type             string            `json:"type"`
	Parameters       map[string]string `json:"parameters"`
	FileName         string            `json:"fileName"`
	CreationDate     time.Time         `json:"creationDate"`
	ModificationDate time.Time         `json:"modificationDate"`
	ReadDate         time.Time         `json:"readDate"`
	Size             int               `json:"size"`
}

// MediaType represents the structure for media type.
type MediaType struct {
	Type            string            `json:"type"`
	Subtype         string            `json:"subtype"`
	Parameters      map[string]string `json:"parameters"`
	WildcardType    bool              `json:"wildcardType"`
	WildcardSubtype bool              `json:"wildcardSubtype"`
}

// Parent represents the structure for parent, similar to Root but includes Parent and BodyParts fields.
type Parent struct {
	ContentDisposition   ContentDisposition               `json:"contentDisposition"`
	Entity               map[string]interface{}           `json:"entity"`
	Headers              map[string][]string              `json:"headers"`
	MediaType            MediaType                        `json:"mediaType"`
	MessageBodyWorkers   map[string]interface{}           `json:"messageBodyWorkers"`
	Parent               string                           `json:"parent"`
	Providers            map[string]interface{}           `json:"providers"`
	BodyParts            []BodyPart                       `json:"bodyParts"`
	ParameterizedHeaders map[string][]ParameterizedHeader `json:"parameterizedHeaders"`
}

// BodyPart represents the structure for body parts within the JSON.
type BodyPart struct {
	ContentDisposition   ContentDisposition               `json:"contentDisposition"`
	Entity               map[string]interface{}           `json:"entity"`
	Headers              map[string][]string              `json:"headers"`
	MediaType            MediaType                        `json:"mediaType"`
	MessageBodyWorkers   map[string]interface{}           `json:"messageBodyWorkers"`
	Parent               string                           `json:"parent"`
	Providers            map[string]interface{}           `json:"providers"`
	ParameterizedHeaders map[string][]ParameterizedHeader `json:"parameterizedHeaders"`
}

// Field represents a more detailed structure within "fields" which can contain nested structures similar to BodyPart.
type Field struct {
	ContentDisposition         ContentDisposition               `json:"contentDisposition"`
	Entity                     map[string]interface{}           `json:"entity"`
	Headers                    map[string][]string              `json:"headers"`
	MediaType                  MediaType                        `json:"mediaType"`
	MessageBodyWorkers         map[string]interface{}           `json:"messageBodyWorkers"`
	Parent                     Parent                           `json:"parent"`
	Providers                  map[string]interface{}           `json:"providers"`
	Name                       string                           `json:"name"`
	Value                      string                           `json:"value"`
	Simple                     bool                             `json:"simple"`
	FormDataContentDisposition ContentDisposition               `json:"formDataContentDisposition"`
	ParameterizedHeaders       map[string][]ParameterizedHeader `json:"parameterizedHeaders"`
}

// ParameterizedHeader represents the structure for parameterized headers.
type ParameterizedHeader struct {
	Value      string            `json:"value"`
	Parameters map[string]string `json:"parameters"`
}
