package rekognition

import (
	"encoding/json"
	"reflect"

	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/fatih/structs"

	intErr "github.com/easynetwork/aws-sdk-go-bindings/internal/error"
)

// CompareFacesInput contains parameters to be sent to CompareFaces
type CompareFacesInput struct {
	SourceImage []byte
	TargetImage []byte
	Similarity  float64
}

// CompareFacesOutput CompareFaces response
type CompareFacesOutput struct {
	FaceMatches    []FaceMatches `json:"FaceMatches"`
	UnmatchedFaces []FaceMatches `json:"UnmatchedFaces"`
}

// DetectFacesInput contains parameters to be sent to DetectFaces
type DetectFacesInput struct {
	SourceImage []byte
}

// DetectFacesOutput contains DetectFaces response
type DetectFacesOutput struct {
	FaceDetails []FaceDetails `json:"FaceDetails"`
}

// DetectTextInput contains parameters to be sent to DetectText
type DetectTextInput struct {
	SourceImage []byte
}

// DetectTextOutput contains DetectText response
type DetectTextOutput struct {
	TextDetections []TextDetection `json:"TextDetections"`
}

// FaceMatches contains CompareFaces Face Matches
type FaceMatches struct {
	Face       Face    `json:"Face,omitempty"`
	Similarity float64 `json:"Similarity,omitempty"`
}

// FaceDetails contains Face Details
type FaceDetails struct {
	AgeRange   AgeRange        `json:"AgeRange,omitempty"`
	Beard      BoolAttribute   `json:"Beard,omitempty"`
	EyeGlasses BoolAttribute   `json:"EyeGlasses,omitempty"`
	EyesOpen   BoolAttribute   `json:"EyesOpen,omitempty"`
	Gender     StringAttribute `json:"Gender,omitempty"`
	MouthOpen  BoolAttribute   `json:"MouthOpen,omitempty"`
	Mustache   BoolAttribute   `json:"Mustache,omitempty"`
	Confidence float64         `json:"Confidence,omitempty"`
	Quality    Quality         `json:"Quality,omitempty"`
	Sunglasses BoolAttribute   `json:"Sunglasses,omitempty"`
}

// TextDetection contains Text Detection output
type TextDetection struct {
	Confidence   float64 `json:"Confidence,omitempty"`
	DetectedText string  `json:"DetectedText,omitempty"`
	Id           int64   `json:"Id,omitempty"`
	Type         string  `json:"Type,omitempty"`
}

// AgeRange describes a low and max age
type AgeRange struct {
	High int64 `json:"High,omitempty"`
	Low  int64 `json:"Low,omitempty"`
}

// BoolAttribute is used to represent a rekognition boolean attribute
type BoolAttribute struct {
	Confidence float64 `json:"Confidence,omitempty"`
	Value      bool    `json:"Value,omitempty"`
}

// StringAttribute is used to represent a rekognition string attribute
type StringAttribute struct {
	Confidence float64 `json:"Confidence,omitempty"`
	Value      string  `json:"Value,omitempty"`
}

// Face describes a face analysis output
type Face struct {
	Confidence float64 `json:"Confidence,omitempty"`
	Quality    Quality `json:"Quality,omitempty"`
}

// Quality describes the quality of an analyzed picture
type Quality struct {
	Brightness float64 `json:"Brightness,omitempty"`
	Sharpness  float64 `json:"Sharpness,omitempty"`
}

// UnmarshalCompareFacesOutput unmarshals a *rekognition.CompareFacesOutput to *CompareFacesOutput
func UnmarshalCompareFacesOutput(input *rekognition.CompareFacesOutput, output *CompareFacesOutput) error {

	if reflect.DeepEqual(input, reflect.Zero(reflect.TypeOf(input)).Interface()) {
		return intErr.Format(Input, ErrEmptyParameter)
	}

	if reflect.ValueOf(output).Kind() != reflect.Ptr {
		return intErr.Format(Output, ErrNoPointerParameter)
	}

	err := unmarshalRekognitionOut(input, output)
	if err != nil {
		return err
	}

	return nil

}

// UnmarshalDetectFacesOutput unmarshals a *rekognition.DetectFacesOutput to *DetectFacesOutput
func UnmarshalDetectFacesOutput(input *rekognition.DetectFacesOutput, output *DetectFacesOutput) error {

	if reflect.DeepEqual(input, reflect.Zero(reflect.TypeOf(input)).Interface()) {
		return intErr.Format(Input, ErrEmptyParameter)
	}

	if reflect.ValueOf(output).Kind() != reflect.Ptr {
		return intErr.Format(Output, ErrNoPointerParameter)
	}
	err := unmarshalRekognitionOut(input, output)
	if err != nil {
		return err
	}
	return nil
}

// UnmarshalDetectTextOutput unmarshals a *rekognition.DetectTextOutput to *DetectTextOutput
func UnmarshalDetectTextOutput(input *rekognition.DetectTextOutput, output *DetectTextOutput) error {

	if reflect.DeepEqual(input, reflect.Zero(reflect.TypeOf(input)).Interface()) {
		return intErr.Format(Input, ErrEmptyParameter)
	}

	if reflect.ValueOf(output).Kind() != reflect.Ptr {
		return intErr.Format(Output, ErrNoPointerParameter)
	}

	err := unmarshalRekognitionOut(input, output)
	if err != nil {
		return err
	}

	return nil

}

// unmarshalRekognitionOut unmarshals a *rekognition.Output to a given interface.
// Returns error if something went wrong.
func unmarshalRekognitionOut(input, output interface{}) error {

	if reflect.DeepEqual(input, reflect.Zero(reflect.TypeOf(input)).Interface()) {
		return intErr.Format(Input, ErrEmptyParameter)
	}

	if reflect.ValueOf(output).Kind() != reflect.Ptr {
		return intErr.Format(Output, ErrNoPointerParameter)
	}

	m := structs.Map(input)

	if len(m) == 0 {
		return intErr.Format(Map, ErrEmptyMap)
	}

	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, output)
	if err != nil {
		return err
	}

	return nil

}

// NewCompareFacesInput builds a *rekognition.CompareFacesInput starting from the two images, their bucket and a similarity threshold
func NewCompareFacesInput(source, target []byte, similarity float64) (*rekognition.CompareFacesInput, error) {

	if len(source) == 0 {
		return nil, intErr.Format(Source, ErrEmptyParameter)
	}

	if len(target) == 0 {
		return nil, intErr.Format(Target, ErrEmptyParameter)
	}

	if similarity == 0 {
		return nil, intErr.Format(Similarity, ErrBadSimilarityParameter)
	}

	newSourceInputImg, err := newInputImage(source)
	if err != nil {
		return nil, err
	}
	newTargetInputImg, err := newInputImage(target)
	if err != nil {
		return nil, err
	}

	out := &rekognition.CompareFacesInput{}
	out = out.SetSimilarityThreshold(similarity)
	out = out.SetSourceImage(newSourceInputImg)
	out = out.SetTargetImage(newTargetInputImg)

	return out, nil

}

// NewDetectFacesInput builds a *rekognition.DetectFacesInput starting from the image
func NewDetectFacesInput(source []byte) (*rekognition.DetectFacesInput, error) {

	if len(source) == 0 {
		return nil, intErr.Format(Source, ErrEmptyParameter)
	}

	newInputImg, err := newInputImage(source)
	if err != nil {
		return nil, err
	}

	out := &rekognition.DetectFacesInput{}
	out = out.SetImage(newInputImg)

	return out, nil

}

// NewDetectTextInput builds a *rekognition.DetectTextInput starting from the image
func NewDetectTextInput(source []byte) (*rekognition.DetectTextInput, error) {

	if len(source) == 0 {
		return nil, intErr.Format(Source, ErrEmptyParameter)
	}

	newInputImg, err := newInputImage(source)
	if err != nil {
		return nil, err
	}

	out := &rekognition.DetectTextInput{}
	out = out.SetImage(newInputImg)

	return out, nil

}

// newInputImage returns a *rekognition.Image given an S3 image []byte encoded
func newInputImage(image []byte) (*rekognition.Image, error) {

	if len(image) == 0 {
		return nil, intErr.Format(Image, ErrEmptyParameter)
	}

	out := &rekognition.Image{
		Bytes: image,
	}

	return out, nil

}
