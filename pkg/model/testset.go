package model

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
)

func NewTestset(name string) *Testset {
	return &Testset{Name: name}
}

type Testset struct {
	Name        string
	Status      int
	Evaluations []*Evaluation
}

type Evaluation struct {
	Left  []Field
	Right interface{} // string / float64 / null
}

type Field struct {
	Value   string
	IsIndex bool
}

func (t Testset) ToString() string {
	var ts string
	for i, e := range t.Evaluations {
		if i == 0 {
			ts += `var data = pm.response.json(); `
		}
		// parse Left
		var left string
		for i := len(e.Left) - 1; i >= 0; i-- {
			if e.Left[i].IsIndex {
				left += fmt.Sprintf(`[%s]`, e.Left[i].Value)
			} else {
				left += fmt.Sprintf(`.%s`, e.Left[i].Value)
			}
		}
		// parse Right
		var right string
		switch r := e.Right.(type) {
		case string:
			right = fmt.Sprintf(`"%s"`, r)
		case int:
			right = fmt.Sprintf(`%d`, r)
		case float32, float64:
			right = fmt.Sprintf(`%g`, r)
		case nil:
			right = "null"
		}
		ts += fmt.Sprintf(`pm.expect(data%s).to.eql(%v); `, left, right)
	}
	return fmt.Sprintf(`pm.test("%s", function(){ responseCode.code === %d; %s})`, t.Name, t.Status, ts)
}

func (t *Testset) BindResponseStatus(i int) error {
	t.Status = i
	return nil
}

func (t *Testset) BindResponseBody(b []byte) error {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	t.Evaluations = t.parseTest(data)
	return nil
}

func (t Testset) parseTest(data interface{}) []*Evaluation {
	var result []*Evaluation

	parseKey := func(key interface{}) Field {
		switch k := key.(type) {
		case int:
			return Field{strconv.Itoa(k), true}
		case string:
			return Field{k, false}
		default:
			panic(fmt.Errorf("TODO"))
		}
	}
	sortKeys := func(keys map[string]interface{}) []string {
		result := make([]string, len(keys))
		idx := 0
		for key := range keys {
			result[idx] = key
			idx++
		}
		sort.Strings(result)
		return result
	}

	var parse func(data, key interface{}) (children []*Evaluation)
	parse = func(data, key interface{}) (children []*Evaluation) {
		switch d := data.(type) {
		case map[string]interface{}:
			var x []*Evaluation
			sortedKeys := sortKeys(d)
			for _, k := range sortedKeys {
				cr := parse(d[k], k)
				for _, c := range cr {
					c.Left = append(c.Left, parseKey(key))
				}
				x = append(x, cr...)
			}
			return x
		case []interface{}:
			var x []*Evaluation
			for i, v := range d {
				cr := parse(v, i)
				for _, c := range cr {
					c.Left = append(c.Left, parseKey(key))
				}
				x = append(x, cr...)
			}
			return x
		default:
			f := parseKey(key)
			children = append([]*Evaluation{}, &Evaluation{
				Left:  []Field{f},
				Right: d,
			})
			return
		}
	}

	switch d := data.(type) {
	case map[string]interface{}:
		sortedKeys := sortKeys(d)
		for _, k := range sortedKeys {
			es := parse(d[k], k)
			result = append(result, es...)
		}
	case []interface{}:
		for i, v := range d {
			es := parse(v, i)
			result = append(result, es...)
		}
	default:
		panic(fmt.Errorf("TODO"))
	}
	return result
}
