package apitool

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Ronmi/rtoolkit/jsonapi"
)

// ParamGreeting represents parameters of Greeting API
type ParamGreeting struct {
	Name    string
	Surname string
}

// RespGreeting represents returned type of Greeting API
type RespGreeting struct {
	Name    string
	Surname string
	Greeted bool
}

// greeting is handler of Greeting API
func Greeting(
	dec *json.Decoder,
	r *http.Request,
	w http.ResponseWriter,
) (interface{}, error) {
	var p ParamGreeting
	if err := dec.Decode(&p); err != nil {
		return nil, jsonapi.APPERR.SetData(
			"parameter format error",
		).SetCode("EParamFormat")
	}

	return RespGreeting{
		Name:    p.Name,
		Surname: p.Surname,
		Greeted: true,
	}, nil
}

// RunAPIServer creates and runs an API server at :9527
func RunAPIServer() *http.Server {
	http.Handle("/greeting", jsonapi.Handler(Greeting))
	server := &http.Server{Addr: ":9527"}
	go server.ListenAndServe()

	return server
}

func ExampleClient() {
	// start the API server
	server := RunAPIServer()
	defer server.Close()

	client := Call("POST", "http://localhost:9527/greeting", nil)

	var resp RespGreeting
	err := client.Exec(ParamGreeting{Name: "John", Surname: "Doe"}, &resp)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf(
		"Are we greeted to %s %s? %v",
		resp.Name, resp.Surname, resp.Greeted,
	)

	// output: Are we greeted to John Doe? true
}
