package daraja

import (
	"bytes"
	"encoding/json"
)

/*
Tolerant decoders for Daraja's result callbacks.

Safaricom is not consistent about how it encodes these fields, and a mismatch
is not a partial failure — encoding/json aborts the whole document, so one
unexpected number discards an entire callback. The three types here absorb the
variations that show up in practice:

  ResultType / ResultCode   quoted on some callbacks, a bare number on others
  ResultParameter[].Value   a string for AccountBalance, a number for
                            BOCompletedTime, in the same array
  ResultParameter           an array when there are several, sometimes a lone
                            object when there is one
  ReferenceItem             the same object-or-array inconsistency

Use FlexString for any scalar Daraja might send unquoted.
*/

// FlexString holds a Daraja scalar that may arrive quoted or bare.
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

// MarshalJSON always emits a quoted string, so a decoded callback re-serialises
// to a stable shape.
func (f FlexString) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(f))
}

func (f FlexString) String() string {
	return string(f)
}

// ResultParameter is one Key/Value pair from a result callback.
type ResultParameter struct {
	Key   string     `json:"Key"`
	Value FlexString `json:"Value"`
}

// ResultParameterList decodes ResultParameter whether Daraja sent an array or
// a single object.
type ResultParameterList []ResultParameter

func (l *ResultParameterList) UnmarshalJSON(data []byte) error {
	items, err := decodeOneOrMany[ResultParameter](data)
	if err != nil {
		return err
	}
	*l = items
	return nil
}

// ReferenceItem is one Key/Value pair from a callback's ReferenceData.
type ReferenceItem struct {
	Key   string     `json:"Key"`
	Value FlexString `json:"Value"`
}

// ReferenceItemList decodes ReferenceItem whether Daraja sent an array or a
// single object. Account balance results send an object; reversal results send
// an array.
type ReferenceItemList []ReferenceItem

func (l *ReferenceItemList) UnmarshalJSON(data []byte) error {
	items, err := decodeOneOrMany[ReferenceItem](data)
	if err != nil {
		return err
	}
	*l = items
	return nil
}

// decodeOneOrMany accepts either a JSON array of T or a single T object.
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

// Get returns the value of the first parameter with the given key.
func (l ResultParameterList) Get(key string) string {
	for _, p := range l {
		if p.Key == key {
			return string(p.Value)
		}
	}
	return ""
}

// Get returns the value of the first reference item with the given key.
func (l ReferenceItemList) Get(key string) string {
	for _, item := range l {
		if item.Key == key {
			return string(item.Value)
		}
	}
	return ""
}
