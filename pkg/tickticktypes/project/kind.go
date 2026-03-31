package project

import (
	"encoding/json"
	"fmt"
)

type Kind string

const (
	KindTask    Kind = "TASK"
	KindNote    Kind = "NOTE"
	KindUnknown Kind = "UNKNOWN"
)

func (k *Kind) UnmarshalJSON(data []byte) error {
	var kind string
	if err := json.Unmarshal(data, &kind); err != nil {
		return err
	}
	switch kind {
	case string(KindTask), string(KindNote), string(KindUnknown):
		*k = Kind(kind)
	default:
		*k = KindUnknown
	}
	return nil
}

func (k Kind) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(k))
}

func (k Kind) String() string {
	switch k {
	case KindTask:
		return "📝Task"
	case KindNote:
		return "📖Note"
	default:
		return "🔧Unknown"
	}
}

func (k *Kind) Set(s string) error {
	switch s {
	case string(KindTask), string(KindNote):
		*k = Kind(s)
	default:
		return fmt.Errorf("invalid project kind %q", s)
	}
	return nil
}

func (k *Kind) Type() string {
	return "ProjectKind"
}
