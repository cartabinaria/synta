package json

import (
	"encoding/json"

	"github.com/cartabinaria/synta"
)

func Convert(syn synta.Synta) (s Synta) {
	s.Definitions = map[string]Definition{}
	for id, def := range syn.Definitions {
		if len(def.Comments) == 0 {
			s.Definitions[string(id)] = Definition{
				Comments: []string{}, Regexp: def.Regexp.String()}
		} else {
			s.Definitions[string(id)] = Definition{
				Comments: def.Comments, Regexp: def.Regexp.String()}
		}
	}
	for _, e := range syn.Filename.Segments {
		seg := Segment{}
		switch e.Kind {
		case synta.SegmentTypeIdentifier:
			seg.Value = string(*e.Value)
			seg.Kind = uint(e.Kind)
			seg.Subsegments = []Segment{}
		case synta.SegmentTypeOptional:
			seg.Value = ""
			seg.Kind = uint(e.Kind)
			seg.Subsegments = getSubSegments(e)
		}
		s.Filename.Segments = append(s.Filename.Segments, seg)
	}
	s.Filename.Extension = string(syn.Filename.Extension)
	return
}

func ToJson(synta synta.Synta) (buf []byte, err error) {
	buf, err = json.Marshal(Convert(synta))
	return
}

func getSubSegments(segment synta.Segment) (subSegments []Segment) {
	for _, e := range segment.Subsegments {
		seg := Segment{}
		switch e.Kind {
		case synta.SegmentTypeIdentifier:
			seg.Value = string(*e.Value)
			seg.Kind = uint(e.Kind)
			seg.Subsegments = []Segment{}
		case synta.SegmentTypeOptional:
			seg.Value = ""
			seg.Kind = uint(e.Kind)
			seg.Subsegments = getSubSegments(e)
		}
		subSegments = append(subSegments, seg)
	}
	return
}
