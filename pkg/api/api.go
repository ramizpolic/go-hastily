package api

import (
	"errors"
	"os"
	"strconv"
	"sync"

	common "github.com/fhivemind/go-hastily/pkg/common"
	tablewriter "github.com/olekukonko/tablewriter"
)

// Tabler imports table controller.
var Tabler = common.Tabler

// ApiModel defines handler object for different backend models.
type ApiModel struct {
	Client *Client
	Name   string
}

// ExportModel defines generic output model.
type ExportModel struct {
	Data        []*Model
	ExtraFields map[string]*common.Generic
	Type        common.TableType
	IsWide      bool
	OutputFile  string
}

// API consumes backend API.
type API interface {
	// htpp get
	Get() ([]*Model, error)
	GetFiltered(*Filter) ([]*Model, error)
	// http create
	Create(*Model) error
	// http delete
	Delete(*Model) Response
	DeleteMany([]*Model) *ResponseList
	// http update
	Update(*Model) Response
	UpdateMany([]*Model, *common.StatusList) *ResponseList
	// object management
	ListFilter([]*Model, *Filter) []*Model
	ListUpdate([]*Model, *Meta) ([]*Model, *common.StatusList)
	// output
	Export(ExportModel) error
}

// NewAPI initializes a specific API.
func NewAPI(model string) ApiModel {
	return ApiModel{
		Client: NewClient(model),
		Name:   model,
	}
}

// Get fetches all objects from backend.
func (api *ApiModel) Get() ([]*Model, error) {
	return api.GetFiltered(nil)
}

// GetFiltered fetches objects from backend that satisfy a specific filter.
func (api *ApiModel) GetFiltered(modelFilter *Filter) ([]*Model, error) {

	// request form
	request := Request{}

	// do request
	var models []*Model
	resp := api.Client.Get(request, &models)
	if !resp.Success {
		return nil, errors.New(resp.Message)
	}

	// return filtered
	return filter(models, modelFilter), nil
}

// Create create provided object on backend.
func (api *ApiModel) Create(model *Model) error {

	// request form
	request := Request{
		Body: model,
	}

	// do request
	resp := api.Client.Post(request, nil)
	if !resp.Success {
		return errors.New(resp.Message)
	}

	// success
	return nil
}

// ListFilter filters objects that satisfy a specific filter.
func (api *ApiModel) ListFilter(models []*Model, modelFilter *Filter) []*Model {
	return filter(models, modelFilter)
}

// Export pretty prints the data to stdout or file.
// Also, if provided with a map based on ids, it will concat them as table columns.
func (api *ApiModel) Export(export ExportModel) error {

	var table = tablewriter.NewWriter(os.Stdout)
	if export.OutputFile != "" {
		// to file
		file, err := os.Create(export.OutputFile)
		if err != nil {
			return err
		}
		defer file.Close()
		table = tablewriter.NewWriter(file)
		export.IsWide = true
	}

	// header defs
	header := []string{"ID"}

	// custom logic for expanded model
	// if model.IsWide {
	// 	header = append(header)
	// }

	// expand default header data
	shouldExpand := export.ExtraFields != nil && len(export.ExtraFields) > 0
	if len(export.Data) > 0 && shouldExpand {
		if val, ok := export.ExtraFields[strconv.Itoa(export.Data[0].ID)]; ok {
			header = append(header, val.Keys...)
			table.SetHeader(header)
		}
	} else {
		table.SetHeader(header)
	}

	// configure table
	export.Type.SetStyleForTable(table, len(header))
	table.SetAutoWrapText(false)

	// populate data
	for _, model := range export.Data {

		// default data
		data := []string{strconv.Itoa(model.ID)}

		// if model.IsWide {
		// 	data = append(data)
		// }

		// expand default header and row data
		if shouldExpand {
			if val, ok := export.ExtraFields[strconv.Itoa(model.ID)]; ok {
				data = append(data, val.Values...)
			}
		}

		table.Append(data)
	}

	table.Render()
	return nil
}

// ListUpdate updates multiple objects based on provided source.
func (api *ApiModel) ListUpdate(models []*Model, source *Meta) ([]*Model, *common.StatusList) {

	// deep-copy models
	dests := make([]*Model, len(models))
	for id, elem := range models {
		copyElem := *elem
		dests[id] = &copyElem
	}

	// async
	var mutex sync.Mutex
	var wg sync.WaitGroup

	// perform object updates
	resp := common.NewStatusList()
	for _, d := range dests {
		wg.Add(1)
		go func(dest *Model) {
			defer wg.Done()
			res := dest.Update(source)
			mutex.Lock()
			defer mutex.Unlock()
			resp.Insert(strconv.Itoa(dest.ID), &res)
		}(d)
	}
	wg.Wait()

	return dests, resp
}

// Delete deletes a specific object in the backend API.
func (api *ApiModel) Delete(model *Model) Response {

	// request form
	request := Request{
		Id: strconv.Itoa(model.ID),
	}

	// do request
	return api.Client.Delete(request, nil)
}

// DeleteMany deletes multiple objects in the backend API.
func (api *ApiModel) DeleteMany(models []*Model) *ResponseList {

	// async
	var mutex sync.Mutex
	var wg sync.WaitGroup

	// perform http deletes
	resp := NewResponseList()
	for _, e := range models {
		wg.Add(1)
		go func(model *Model) {
			defer wg.Done()
			res := api.Delete(model)
			mutex.Lock()
			defer mutex.Unlock()
			resp.Insert(strconv.Itoa(model.ID), &res)
		}(e)
	}
	wg.Wait()
	return resp
}

// Update updates a specific object in the backend API.
func (api *ApiModel) Update(model *Model) Response {

	// request form
	request := Request{
		Id:   strconv.Itoa(model.ID),
		Body: model,
	}

	// do request
	return api.Client.Put(request, nil)
}

// UpdateMany updates multiple objects in the backend API.
func (api *ApiModel) UpdateMany(models []*Model, statuses *common.StatusList) *ResponseList {

	// async
	var mutex sync.Mutex
	var wg sync.WaitGroup

	// perform http updates
	resp := NewResponseList()
	for _, e := range models {
		wg.Add(1)
		go func(model *Model) {
			defer wg.Done()
			var (
				res       Response
				doRequest = true
			)
			// status checks
			if statuses != nil {
				status, ok := statuses.Get(strconv.Itoa(model.ID))
				if ok && !status.Success {
					res = api.Client.DefaultResponse("", errors.New(status.Operation))
					doRequest = false
				}
			}
			// do request
			if doRequest {
				res = api.Update(model)
			}
			mutex.Lock()
			defer mutex.Unlock()
			resp.Insert(strconv.Itoa(model.ID), &res)
		}(e)
	}
	wg.Wait()
	return resp
}

// filter returns the list of objects which satisfy the filtering options.
func filter(models []*Model, filter *Filter) (ret []*Model) {
	// process data
	for _, e := range models {
		if e.ValidForFilter(filter) {
			ret = append(ret, e)
		}
	}
	return
}

// filterAsync returns the list of objects which satisfy the filtering options.
// This is slightly slower on <= 4 threads than its sync sibling.
func filterAsync(models []*Model, filter *Filter) (ret []*Model) {

	// async
	var mutex sync.Mutex
	var wg sync.WaitGroup

	// filter
	for _, e := range models {
		wg.Add(1)
		go func(model *Model) {
			defer wg.Done()
			res := model.ValidForFilter(filter)
			mutex.Lock()
			defer mutex.Unlock()
			if res {
				// might be dangerous
				ret = append(ret, model)
			}
		}(e)
	}
	wg.Wait()

	return
}
