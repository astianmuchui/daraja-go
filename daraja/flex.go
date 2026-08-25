package daraja

import (
	"bytes"
	"encoding/json"
)

type FlexString string

func (f *FlexString) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		*f = ""
		return nil
	}

	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		*f = FlexString(s)
		return nil
	}

	/*
		Numbers and booleans are kept as their literal source bytes rather than
		routed through float64. BOCompletedTime is a 14-digit timestamp
		(20260825140947); a float64 round-trip would render it in exponent
		notation and lose the exact digits.
	*/
	*f = FlexString(data)
	return nil
}



func (f FlexString) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(f))
}

func (f FlexString) String() string {
	return string(f)
}


type ResultParameter struct {
	Key   string     `json:"Key"`
	Value FlexString `json:"Value"`
}

type ResultParameterList []ResultParameter

func (l *ResultParameterList) UnmarshalJSON(data []byte) error {
	items, err := decodeOneOrMany[ResultParameter](data)
	if err != nil {
		return err
	}
	*l = items
	return nil
}
type ReferenceItem struct {
	Key   string     `json:"Key"`
	Value FlexString `json:"Value"`
}

type ReferenceItemList []ReferenceItem

func (l *ReferenceItemList) UnmarshalJSON(data []byte) error {
	items, err := decodeOneOrMany[ReferenceItem](data)
	if err != nil {
		return err
	}
	*l = items
	return nil
}

func decodeOneOrMany[T any](data []byte) ([]T, error) {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		return nil, nil
	}

	if data[0] == '[' {
		var many []T
		if err := json.Unmarshal(data, &many); err != nil {
			return nil, err
		}
		return many, nil
	}

	var one T
	if err := json.Unmarshal(data, &one); err != nil {
		return nil, err
	}
	return []T{one}, nil
}


func (l ResultParameterList) Get(key string) string {
	for _, p := range l {
		if p.Key == key {
			return string(p.Value)
		}
	}
	return ""
}


func (l ReferenceItemList) Get(key string) string {
	for _, item := range l {
		if item.Key == key {
			return string(item.Value)
		}
	}
	return ""
}
