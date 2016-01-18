package thunder

import (
	"fmt"
	"reflect"

	"github.com/omeid/thunder/codec"
)

type setter func(k, v []byte) error
type Mattcher func(interface{}) bool

func mapSetter(kc, vc codec.Codec, mv reflect.Value, match Mattcher) (setter, error) {

	mv = reflect.Indirect(mv)
	mvt := mv.Type()
	kt := mvt.Key()
	vt := mvt.Elem()

	if mv.IsNil() {
		return nil, fmt.Errorf("Thunder: Bucket.Where expects non nil map.")
	}

	return func(k, v []byte) error {

		key, value := reflect.New(kt), reflect.New(vt)
		err := unmarshal(kc, vc, k, v, key.Interface(), value.Interface())
		if err != nil {
			return err
		}
		if kt.Kind() != reflect.Ptr {
			key = reflect.Indirect(key)
		}
		if vt.Kind() != reflect.Ptr {
			value = reflect.Indirect(value)
		}

		if match(value.Interface()) {
			mv.SetMapIndex(key, value)
		}
		return nil
	}, nil
}

func sliceSetter(kc, vc codec.Codec, svp reflect.Value, match Mattcher) (setter, error) {

	//svp: Slice Value Pointer
	//sv Slice Value
	//vt Value Type

	if svp.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("Thunder:Where Expects Slice to be a pointer. %s", svp.Kind())
	}

	sv := svp.Elem()
	vt := sv.Type().Elem()
	return func(k, v []byte) error {
		value := reflect.New(vt)

		err := vc.Unmarshaler(value.Interface()).UnmarshalBinary(v)
		if err != nil {
			return err
		}

		if vt.Kind() != reflect.Ptr {
			value = reflect.Indirect(value)
		}

		if match(value.Interface()) {
			sv.Set(reflect.Append(sv, value))
		}
		return nil
	}, nil
}

func (b *Bucket) Where(match Mattcher, kvm interface{}) error {

	if b.err != nil {
		return b.err
	}

	mv := reflect.ValueOf(kvm)

	var (
		apply setter
		err   error
	)

	switch mv.Type().Elem().Kind() {
	case reflect.Map:
		apply, err = mapSetter(b.kc, b.vc, mv, match)
	case reflect.Slice:
		apply, err = sliceSetter(b.kc, b.vc, mv, match)
	default:
		err = fmt.Errorf("Thunder:Where Expects a Map or Slice. %T given.", kvm)
	}

	if err != nil {
		return err
	}
	return b.bucket.ForEach(func(k, v []byte) error {
		return apply(k, v)

	})
}
