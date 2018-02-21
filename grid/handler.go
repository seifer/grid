package grid

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	html "github.com/antchfx/xquery/html"
)

type Handler struct {
}

type Request struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

type Response struct {
	Id  string `json:"id"`
	Url string `json:"url"`
	Err string `json:"err,omitempty"`

	Data *ResponseData `json:"data,omitempty"`
}

type ResponseData struct {
	Title string `json:"title"`
	Price string `json:"price"`
	Image string `json:"image"`

	InStock bool `json:"in_stock"`
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	requests := []*Request{}

	err := json.NewDecoder(r.Body).Decode(&requests)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(processRequests(requests))
}

func processRequests(requests []*Request) *[]*Response {
	mutex := &sync.Mutex{}
	responses := make([]*Response, 0, len(requests))

	wg := &sync.WaitGroup{}
	for i := range requests {
		wg.Add(1)

		go func(r *Request) {
			defer wg.Done()

			data, err := processRequest(r.Url)

			mutex.Lock()
			if err != nil {
				responses = append(responses, &Response{
					Id:  r.Id,
					Url: r.Url,
					Err: err.Error(),
				})
			} else {
				responses = append(responses, &Response{
					Id:   r.Id,
					Url:  r.Url,
					Data: data,
				})
			}
			mutex.Unlock()
		}(requests[i])
	}
	wg.Wait()

	return &responses
}

func processRequest(url string) (*ResponseData, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	node, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	title := extractTitle(node)
	if title == "" {
		return nil, errors.New("Node title is nil")
	}

	price := extractPrice(node)
	if price == "" {
		return nil, errors.New("Node price is nil")
	}

	image := extractImage(node)
	if image == "" {
		return nil, errors.New("Node image is nil")
	}

	return &ResponseData{
		Title:   title,
		Price:   price,
		Image:   image,
		InStock: extractIsInStock(node),
	}, nil
}
