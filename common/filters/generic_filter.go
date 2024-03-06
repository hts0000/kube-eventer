package filters

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	v1 "k8s.io/api/core/v1"
	log "k8s.io/klog"
)

type GenericFilter struct {
	field  string
	keys   []string
	regexp bool
}

func IsZero(v reflect.Value) bool {
	return !v.IsValid() || reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

func (gf *GenericFilter) Filter(event *v1.Event) (matched bool) {
	var field reflect.Value

	switch gf.field {
	case "Kind":
		field = reflect.Indirect(reflect.ValueOf(event)).FieldByNameFunc(func(name string) bool {
			return name == "InvolvedObject"
		}).FieldByName("Kind")
	case "Namespace":
		field = reflect.Indirect(reflect.ValueOf(event)).FieldByNameFunc(func(name string) bool {
			return name == "InvolvedObject"
		}).FieldByName("Namespace")
	case "Type":
		field = reflect.Indirect(reflect.ValueOf(event)).FieldByName("Type")
	case "Reason":
		field = reflect.Indirect(reflect.ValueOf(event)).FieldByName("Reason")
	case "Object":
		log.Error("@@@@@@@@@@ hit Object case")
		field = reflect.Indirect(reflect.ValueOf(event)).FieldByName("InvolvedObject")
		log.Infof("event: %#v, field: %#v\n", event, field.String())
		field = reflect.Indirect(reflect.ValueOf(event)).FieldByName("ObjectMeta")
		log.Infof("event: %#v, field: %#v, field: %#v\n", event, field.String(), field)
		for _, k := range gf.keys {
			// 包含子串时，希望过滤掉改子串
			filteFlag := !(k != "" && k[0] == '!')
			s := k[1:]
			if strings.Contains(field.String(), s) {
				// 包含子串，但是希望过滤掉
				return true && filteFlag
			}
		}
		return false
	}

	if IsZero(field) {
		log.Error("************ error field")
		return false
	}

	log.Infof("&&&&&&&&&&&&&&&&&& gf:keys: %#v\n", gf.keys)
	for _, k := range gf.keys {
		// enable regexp
		if gf.regexp {
			if ok, err := regexp.Match(k, []byte(field.String())); err == nil && ok {
				matched = true
				fmt.Printf("!!!!!!!!!!!!!!!!!! match fiele: %#v\n", field.String())
				return
			} else {
				if err != nil {
					fmt.Printf("################### not match fiele: %#v\n", field.String())
					log.Errorf("Failed to match pattern %s with %s,because of %v", k, field.String(), err)
				}
			}
		} else {
			if field.String() == k {
				matched = true
				return
			}
		}
	}
	return false
}

// Generic Filter
func NewGenericFilter(field string, keys []string, regexp bool) *GenericFilter {
	k := &GenericFilter{
		field:  field,
		regexp: regexp,
	}
	if keys != nil {
		k.keys = keys
		return k
	}
	k.keys = make([]string, 0)
	return k
}
