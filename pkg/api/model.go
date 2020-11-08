package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	common "github.com/fhivemind/go-hastily/pkg/common"
	. "github.com/fhivemind/go-hastily/pkg/global"
	"github.com/ghodss/yaml"
	"github.com/imdario/mergo"
	"github.com/olekukonko/tablewriter"
	"github.com/r3labs/diff/v2"
)

// Model represents generic data model for backend API.
// Inherit this object to implement all API functionalities.
type Model struct {
	ID int `json:"id"`
}

// Filter defines which filters can be applied to Model.
type Filter struct {
	ID int `json:"id"`
}

// Meta holds Model object and its internal byte representation.
// This is important as some objects might contain nullable fields
// e.g. when loading incomplete object from file.
type Meta struct {
	Model Model
	Data  []byte
}

// Print prints filters to console.
func (filter *Filter) Print() {
	table := tablewriter.NewWriter(os.Stdout)
	keys, values := common.StructNonNullKeysAndValues(filter)

	// header
	table.SetHeader(keys)

	// configure table
	ttype := common.Tabler.Basic
	ttype.SetStyleForTable(table, len(keys))
	table.SetAutoWrapText(false)

	// row
	table.Append(values)

	// print
	table.Render()
}

// Print prints Model object to console.
func (model *Model) Print() {
	table := tablewriter.NewWriter(os.Stdout)
	keys, values := common.StructNonNullKeysAndValues(model)

	// configure table
	table.SetAutoWrapText(false)

	// update data
	for i := range keys {
		if i == 0 {
			table.SetHeader([]string{keys[i], values[i]})
		} else {
			table.Append([]string{keys[i], values[i]})
		}
	}

	// vertical style
	ttype := common.Tabler.Vertical
	ttype.SetStyleForTable(table, len(keys))

	// print
	table.Render()
}

// Update updates Model based on provided source.
func (model *Model) Update(source *Meta) common.Status {

	// update and override dest values with source values
	dest := *model
	dest.Merge(source)
	changes, err := diff.Diff(model, &dest)

	// verify merge
	if err != nil {
		return common.Status{
			Success:   false,
			Operation: fmt.Sprintf("%+v", err),
		}
	} else if len(changes) == 0 {
		return common.Status{
			Success:   false,
			Operation: "no change",
		}
	}

	// update
	*model = dest

	return common.Status{
		Success:   true,
		Operation: fmt.Sprintf("%+v", changes),
	}
}

// Merge adds and overrides everything on model from source.
func (model *Model) Merge(source *Meta) error {

	// convert to maps
	var sourceObject map[string]interface{}
	err := json.Unmarshal(source.Data, &sourceObject)
	if err != nil {
		return err
	}
	modelObject := common.ObjectToMap(model)

	// merge maps
	if err = mergo.Merge(&modelObject, sourceObject, mergo.WithOverride); err != nil {
		return err
	}

	// convert to Model
	var result Model
	byt, _ := json.Marshal(modelObject)
	json.Unmarshal(byt, &result)

	// update
	*model = result

	return nil
}

// ValidForFilter checks if a given Model is valid for a specific filter.
func (model *Model) ValidForFilter(filter *Filter) bool {

	// no filter applied, skip
	if filter == nil {
		return true
	}

	// apply custom filters
	modelMap := common.ObjectToMap(model)
	filterMap := common.ObjectToMap(filter)

	// remove unused filters here
	// delete(filterMap, "NAME")

	// filter
	for fKey, fVal := range filterMap {
		if !IsZero(fVal) && modelMap[fKey] != fVal {
			//fmt.Printf("Diff %+v:   expected %+v   got %+v\n", fKey, fVal, modelMap[fKey])
			return false
		}
	}

	return true
}

// FromFile parses yaml file into Meta object.
func (meta *Meta) FromFile(file string) error {

	// read file
	yml, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	// convert yaml to json
	byt, err := yaml.YAMLToJSON(yml)
	if err != nil {
		return err
	}

	// extract to Model
	var model Model
	err = yaml.Unmarshal(yml, &model)
	if err != nil {
		return err
	}

	// update
	meta.Model = model
	meta.Data = byt

	return nil
}

// Print prints Meta object to console.
func (meta *Meta) Print() {
	meta.Model.Print()
}

// IsJsonModel checks if bytes string represents a valid Model object.
func IsJsonModel(data []byte) bool {
	var model Model
	if json.Unmarshal(data, &model) != nil {
		return false
	}
	return true
}
