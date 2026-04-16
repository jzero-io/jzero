package config

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"unicode"

	"github.com/iancoleman/orderedmap"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// MarshalDebugJSONForCommand marshals the effective config for the active command path.
func MarshalDebugJSONForCommand(cfg Config, cmd *cobra.Command) ([]byte, error) {
	value := debugJSONForCommand(cfg, cmd)

	raw, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	if err := json.Indent(&out, raw, "", "  "); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func debugJSONForCommand(cfg Config, cmd *cobra.Command) *orderedmap.OrderedMap {
	root := orderedmap.New()

	for _, scope := range debugFlagScopes(cmd) {
		scope := scope
		scope.flags.VisitAll(func(flag *pflag.Flag) {
			fullPath := append(append([]string{}, scope.path...), flag.Name)

			value, ok := configValueByPath(reflect.ValueOf(cfg), fullPath)
			if !ok {
				return
			}

			setDebugJSONValue(root, fullPath, debugJSONValue(value))
		})
	}

	return root
}

type debugFlagScope struct {
	path  []string
	flags *pflag.FlagSet
}

func debugFlagScopes(cmd *cobra.Command) []debugFlagScope {
	chain := commandChain(cmd)
	if len(chain) == 0 {
		return nil
	}

	scopes := []debugFlagScope{
		{flags: chain[0].PersistentFlags()},
	}

	for i := 1; i < len(chain); i++ {
		scope := debugFlagScope{
			path: commandNames(chain[1 : i+1]),
		}

		if i == len(chain)-1 {
			scope.flags = chain[i].NonInheritedFlags()
		} else {
			scope.flags = chain[i].PersistentFlags()
		}

		scopes = append(scopes, scope)
	}

	return scopes
}

func commandChain(cmd *cobra.Command) []*cobra.Command {
	var chain []*cobra.Command
	for current := cmd; current != nil; current = current.Parent() {
		chain = append(chain, current)
	}

	for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
		chain[i], chain[j] = chain[j], chain[i]
	}

	return chain
}

func commandNames(commands []*cobra.Command) []string {
	names := make([]string, 0, len(commands))
	for _, command := range commands {
		names = append(names, command.Name())
	}
	return names
}

func configValueByPath(v reflect.Value, path []string) (reflect.Value, bool) {
	current := v

	for _, segment := range path {
		for current.Kind() == reflect.Interface || current.Kind() == reflect.Pointer {
			if current.IsNil() {
				return reflect.Value{}, false
			}
			current = current.Elem()
		}

		if current.Kind() != reflect.Struct {
			return reflect.Value{}, false
		}

		next, ok := findConfigField(current, segment)
		if !ok {
			return reflect.Value{}, false
		}

		current = next
	}

	return current, true
}

func findConfigField(v reflect.Value, segment string) (reflect.Value, bool) {
	t := v.Type()
	want := normalizeDebugSegment(segment)

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		for _, candidate := range debugFieldKeys(field) {
			if normalizeDebugSegment(candidate) == want {
				return v.Field(i), true
			}
		}
	}

	return reflect.Value{}, false
}

func debugFieldKeys(field reflect.StructField) []string {
	keys := make([]string, 0, 2)

	tag := field.Tag.Get("mapstructure")
	if tag != "" {
		key := strings.Split(tag, ",")[0]
		if key == "-" {
			return nil
		}
		if key != "" {
			keys = append(keys, key)
		}
	}

	keys = append(keys, lowerFirst(field.Name))
	return keys
}

func normalizeDebugSegment(s string) string {
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, "_", "")
	return strings.ToLower(s)
}

func setDebugJSONValue(root *orderedmap.OrderedMap, path []string, value any) {
	current := root

	for i, segment := range path {
		if i == len(path)-1 {
			current.Set(segment, value)
			return
		}

		next, ok := current.Get(segment)
		if ok {
			if child, ok := next.(*orderedmap.OrderedMap); ok {
				current = child
				continue
			}
		}

		child := orderedmap.New()
		current.Set(segment, child)
		current = child
	}
}

func debugJSONValue(v reflect.Value) any {
	if !v.IsValid() {
		return nil
	}

	switch v.Kind() {
	case reflect.Interface, reflect.Pointer:
		if v.IsNil() {
			return nil
		}
		return debugJSONValue(v.Elem())
	case reflect.Struct:
		obj := orderedmap.New()
		t := v.Type()

		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			if !field.IsExported() {
				continue
			}

			keys := debugFieldKeys(field)
			if len(keys) == 0 {
				continue
			}

			obj.Set(keys[0], debugJSONValue(v.Field(i)))
		}

		return obj
	case reflect.Slice, reflect.Array:
		if v.Kind() == reflect.Slice && v.IsNil() {
			return nil
		}

		items := make([]any, v.Len())
		for i := 0; i < v.Len(); i++ {
			items[i] = debugJSONValue(v.Index(i))
		}
		return items
	case reflect.Map:
		if v.IsNil() {
			return nil
		}

		if v.Type().Key().Kind() != reflect.String {
			return v.Interface()
		}

		obj := orderedmap.New()
		for _, key := range v.MapKeys() {
			obj.Set(key.String(), debugJSONValue(v.MapIndex(key)))
		}
		return obj
	default:
		return v.Interface()
	}
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}
