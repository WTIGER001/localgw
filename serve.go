package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/blmayer/awslambdarpc/client"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware/denco"
	"github.com/go-openapi/spec"
)

var router *denco.Router

type LocalGWConfig struct {
	paths map[string]string
}

func Serve(swaggerFile string) {

	fmt.Printf("SERVing %v\n", swaggerFile)
	// Read and parse
	specDoc, err := loads.Spec(swaggerFile)
	if err != nil {
		panic(err)
	}
	a := analysis.New(specDoc.Spec())

	// Build the router
	router = denco.New()
	var records []denco.Record
	for k, v := range a.AllPaths() {
		k = strings.ReplaceAll(k, "{", ":")
		k = strings.ReplaceAll(k, "}", "")
		fmt.Printf("Adding Route: %v\n", k)
		records = append(records, denco.NewRecord(k, v))
	}

	if err := router.Build(records); err != nil {
		log.Fatalf("Could not build router, %v", err)
	}

	fmt.Printf("\nREADY\n\n")

	// Set up hanlder
	http.HandleFunc("/", handleRequests)

	err = http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func handleRequests(w http.ResponseWriter, r *http.Request) {
	path := r.RequestURI

	data, params, found := router.Lookup(path)

	if !found {
		w.WriteHeader(500)

		fmt.Printf("Route not found: %v\n", path)
		io.WriteString(w, string("ROUTE NOT FOUND"))
		return
	}

	item := data.(spec.PathItem)

	pathParams := make(map[string]string)
	for _, p := range params {
		pathParams[p.Name] = p.Value
	}

	method := r.Method
	// Determine which port to send to

	// Create API GW Event
	e := GeneratePayload(r)
	e.PathParameters = pathParams
	payload, err := json.Marshal(e)
	if err != nil {
		w.WriteHeader(500)

		fmt.Printf("BAD PAYLOAD: %v\n", path)
		io.WriteString(w, string("BAD PAYLOAD"))
		return
	}

	var op *spec.Operation

	switch method {
	case "GET":
		op = item.Get
	case "DELETE":
		op = item.Delete
	case "HEAD":
		op = item.Head
	case "OPTIONS":
		op = item.Options
	case "PATCH":
		op = item.Patch
	case "POST":
		op = item.Post
	case "PUT":
		op = item.Put
	}

	// Get the extension object
	ext := GetExtension(op)
	if ext == nil {
		fmt.Printf("NO EXTENSON FOR: %v\n", path)
		w.WriteHeader(500)

		io.WriteString(w, string("NO OEXTENSON"))
		return
	}
	addr := fmt.Sprintf("localhost:%v", ext.Port)
	fmt.Printf("Handling %v --> %v\t", path, addr)

	res, err := client.Invoke(addr, payload)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		w.WriteHeader(500)
		io.WriteString(w, string(err.Error()))
	}

	fmt.Println("OK")

	fmt.Printf("\n\n%+v\n\n", string(res))

	var resp events.APIGatewayProxyResponse
	err = json.Unmarshal(res, &resp)
	if err != nil {
		fmt.Printf("UnMarshall Err: %v\n", err)
	}
	fmt.Printf("\n\n%+v\n\n", resp)
	w.WriteHeader(resp.StatusCode)
	for k, h := range resp.Headers {
		w.Header().Add(k, h)
	}

	w.Header().Add("content-type", "application/json")
	io.WriteString(w, resp.Body)
}

func GeneratePayload(r *http.Request) *events.APIGatewayProxyRequest {
	var body string
	if r.Body != nil {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("CANNOT Gen payload., %v", err)
		}
		body = string(data)
	}

	headers := make(map[string]string)
	mheaders := make(map[string][]string)
	qparams := make(map[string]string)
	mqparams := make(map[string][]string)

	for k, v := range r.Header {
		if len(v) == 1 {
			headers[k] = v[0]
		} else {
			mheaders[k] = v
		}
	}

	for k, v := range r.URL.Query() {
		if len(v) == 1 {
			qparams[k] = v[0]
		} else {
			mqparams[k] = v
		}
	}

	e := events.APIGatewayProxyRequest{
		Resource:                        "/{proxy+}",
		Path:                            r.RequestURI,
		HTTPMethod:                      r.Method,
		Body:                            body,
		Headers:                         headers,
		MultiValueHeaders:               mheaders,
		QueryStringParameters:           qparams,
		MultiValueQueryStringParameters: mqparams,
	}

	return &e
}
