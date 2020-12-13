package utils

import (
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	pref "google.golang.org/protobuf/reflect/protoreflect"
	"math/bits"
	"sort"
	"strconv"
	"unicode/utf8"
)

type jsonWriter struct {
	buf []byte
}

func (w *jsonWriter) write(s string) {
	w.buf = append(w.buf, s...)
}

func (w *jsonWriter) writeString(s string) error {
	var err error
	if w.buf, err = appendString(w.buf, s); err != nil {
		return err
	}
	return nil
}

func Marshal(m proto.Message) ([]byte, error) {
	w := jsonWriter{}
	err := w.marshalMessage(m.ProtoReflect())
	return w.buf, err
}

// marshalMessage marshals the given protoreflect.Message.
func (w *jsonWriter) marshalMessage(m pref.Message) error {
	if err := w.marshalFields(m); err != nil {
		return err
	}

	return nil
}

// marshalFields marshals the fields in the given protoreflect.Message.
func (w *jsonWriter) marshalFields(m pref.Message) error {
	messageDesc := m.Descriptor()

	w.write("{")
	defer w.write("}")
	firstField := true

	// Marshal out known fields.
	fieldDescs := messageDesc.Fields()
	for i := 0; i < fieldDescs.Len(); {
		fd := fieldDescs.Get(i)
		if od := fd.ContainingOneof(); od != nil {
			fd = m.WhichOneof(od)
			i += od.Fields().Len()
			if fd == nil {
				continue // unpopulated oneofs are not affected by EmitUnpopulated
			}
		} else {
			i++
		}

		val := m.Get(fd)
		if !firstField {
			w.write(",")
		}
		if err := w.marshalField(val, fd); err != nil {
			return err
		}
		firstField = false
	}
	return nil
}

func (w *jsonWriter) marshalField(val pref.Value, fd pref.FieldDescriptor) error {
	w.write(`"` + fd.JSONName() + `":`)
	return w.marshalValue(val, fd)
}

// marshalValue marshals the given protoreflect.Value.
func (w *jsonWriter) marshalValue(val pref.Value, fd pref.FieldDescriptor) error {
	switch {
	case fd.IsList():
		return w.marshalList(val.List(), fd)
	case fd.IsMap():
		return w.marshalMap(val.Map(), fd)
	default:
		return w.marshalSingular(val, fd)
	}
}

// marshalList marshals the given protoreflect.List.
func (w *jsonWriter) marshalList(list pref.List, fd pref.FieldDescriptor) error {
	w.write("[")
	defer w.write("]")

	comma := ""
	for i := 0; i < list.Len(); i++ {
		w.write(comma)
		item := list.Get(i)
		if err := w.marshalSingular(item, fd); err != nil {
			return err
		}
		comma = ","
	}
	return nil
}

type mapEntry struct {
	key   pref.MapKey
	value pref.Value
}

// marshalMap marshals given protoreflect.Map.
func (w *jsonWriter) marshalMap(mmap pref.Map, fd pref.FieldDescriptor) error {
	// Get a sorted list based on keyType first.
	entries := make([]mapEntry, 0, mmap.Len())
	mmap.Range(func(key pref.MapKey, val pref.Value) bool {
		entries = append(entries, mapEntry{key: key, value: val})
		return true
	})
	sortMap(fd.MapKey().Kind(), entries)

	w.write(`{`)
	defer w.write(`}`)
	comma := ""

	// Write out sorted list.
	for _, entry := range entries {
		w.write(comma)
		w.write(`"` + entry.key.String() + `":`)
		if err := w.marshalSingular(entry.value, fd.MapValue()); err != nil {
			return err
		}
		comma = ","
	}
	return nil
}

// sortMap orders list based on value of key field for deterministic ordering.
func sortMap(keyKind pref.Kind, values []mapEntry) {
	sort.Slice(values, func(i, j int) bool {
		switch keyKind {
		case pref.Int32Kind, pref.Sint32Kind, pref.Sfixed32Kind,
			pref.Int64Kind, pref.Sint64Kind, pref.Sfixed64Kind:
			return values[i].key.Int() < values[j].key.Int()

		case pref.Uint32Kind, pref.Fixed32Kind,
			pref.Uint64Kind, pref.Fixed64Kind:
			return values[i].key.Uint() < values[j].key.Uint()
		}
		return values[i].key.String() < values[j].key.String()
	})
}

// marshalSingular marshals the given non-repeated field value. This includes
// all scalar types, enums, messages, and groups.
func (w *jsonWriter) marshalSingular(val pref.Value, fd pref.FieldDescriptor) error {
	if !val.IsValid() {
		return nil
	}

	switch kind := fd.Kind(); kind {
	case pref.Int64Kind:
		w.write(val.String())

	case pref.FloatKind:
		w.write(val.String())

	case pref.BytesKind:
		if err := w.writeString(string(val.Bytes())); err != nil {
			return err
		}

	case pref.MessageKind, pref.GroupKind:
		if err := w.marshalMessage(val.Message()); err != nil {
			return err
		}

	default:
		panic(fmt.Sprintf("%v has unknown kind: %v", fd.FullName(), kind))
	}
	return nil
}

// Sentinel error used for indicating invalid UTF-8.
var errInvalidUTF8 = errors.New("invalid UTF-8")

func appendString(out []byte, in string) ([]byte, error) {
	out = append(out, '"')
	i := indexNeedEscapeInString(in)
	in, out = in[i:], append(out, in[:i]...)
	for len(in) > 0 {
		switch r, n := utf8.DecodeRuneInString(in); {
		case r == utf8.RuneError && n == 1:
			return out, errInvalidUTF8
		case r < ' ' || r == '"' || r == '\\':
			out = append(out, '\\')
			switch r {
			case '"', '\\':
				out = append(out, byte(r))
			case '\b':
				out = append(out, 'b')
			case '\f':
				out = append(out, 'f')
			case '\n':
				out = append(out, 'n')
			case '\r':
				out = append(out, 'r')
			case '\t':
				out = append(out, 't')
			default:
				out = append(out, 'u')
				out = append(out, "0000"[1+(bits.Len32(uint32(r))-1)/4:]...)
				out = strconv.AppendUint(out, uint64(r), 16)
			}
			in = in[n:]
		default:
			i := indexNeedEscapeInString(in[n:])
			in, out = in[n+i:], append(out, in[:n+i]...)
		}
	}
	out = append(out, '"')
	return out, nil
}

// indexNeedEscapeInString returns the index of the character that needs
// escaping. If no characters need escaping, this returns the input length.
func indexNeedEscapeInString(s string) int {
	for i, r := range s {
		if r < ' ' || r == '\\' || r == '"' || r == utf8.RuneError {
			return i
		}
	}
	return len(s)
}
