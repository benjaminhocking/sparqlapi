package main

import (
	"fmt"
	"net/http"
	"log"
	"sparqlAPI/utils"
	"strings"
	"encoding/json"
)

//turtle file and custom graph file locations
const (
	filePath = "finance.ttl"
	writeFilePath = "graph_2.json"
	graphDemo1 = "data/dataset1.json"
	graphDemo2 = "data/dataset2.json"
	indexPage = "pages/index.html"
)

//new func type to enable setupRoutes()
type handlerFunc func(http.ResponseWriter, *http.Request)

//serve index page
func indexHandler(w http.ResponseWriter, req *http.Request){
	createGraph(w, req)
	fmt.Fprintf(w, utils.ReadFile(indexPage))
}

//refresh the graph from the current .ttl file
func createGraph(w http.ResponseWriter, req *http.Request){
	/*
	Want nodes and links of form:
		Nodes:
			[
				{"name": <name>, "id": <unique_id>},...
			]
		Links:
			[
				{"source": <source>, "target": <target>, "type": <predicate>,}
			]
	*/
	//get all resources with an associated id
	resources := utils.GetResourceMap(filePath) //type: map[int]string
	//format these resources to match d3.js requirements
	nodes := utils.GetNodes(resources) //type: []map[string]interface
	//get relations between resources
	links := utils.GetLinks(nodes, filePath) //type: []map[string]interface
	//combine into one map
	total := map[string][]map[string]any{"nodes":nodes, "links":links}
	//jsonify this Go map
	res, _ := json.Marshal(total)
	//save the result to the graph file
	utils.WriteFile(writeFilePath, string(res))
}


func newRecHandler(w http.ResponseWriter, req *http.Request){
	//get triple
	textQuery := req.URL.Query()["text"]
	text := textQuery[0]
	//create triple if textQuery was valid, else return err
	t, err := utils.CreateTriple(text, filePath)
	if err != nil{
		fmt.Fprintf(w, "Invalid Record")
		return
	}
	//save this triple to .ttl file
	utils.SaveTriple(t, filePath)
	fmt.Fprintf(w, "Record Saved")
}

func queryHandler(w http.ResponseWriter, req *http.Request){
	//get query
	textQuery := req.URL.Query()["text"]
	text := strings.ReplaceAll(textQuery[0], "<br>", "\n")
	//get results of query using utils
	res := utils.MakeQuery(text, filePath)
	//convert results to string
	resStr := utils.StringRes(res)
	resp := map[string]string{"response":resStr}
	//jsonify result
	b, _ := json.Marshal(resp)
	//write result to response
	fmt.Fprintf(w, string(b))
}


func graph1Handler(w http.ResponseWriter, req *http.Request){
	//serve graph data
	cont := utils.ReadFile(writeFilePath)
	fmt.Fprintf(w, cont)
}

func setupRoutes(m map[string]handlerFunc, mux *http.ServeMux){
	//register paths to their handler functions
	for k,v := range m{
		mux.HandleFunc(k,v)
	}
}

func dataset1Handler(w http.ResponseWriter, req *http.Request){
	//serve demo dataset 1
	cont := utils.ReadFile(graphDemo1)
	fmt.Fprintf(w, cont)
}

func dataset2Handler(w http.ResponseWriter, req *http.Request){
	//serve demo dataset 2
	cont := utils.ReadFile(graphDemo2)
	fmt.Fprintf(w, cont)
}

func main(){
	//listen on 8080
	listenAddr := ":8080"

	//create new serve mux
	mux := http.NewServeMux()
	
	//list routes
	routes := map[string]handlerFunc{
		"/": indexHandler,
		"/new_rec": newRecHandler,
		"/query": queryHandler,
		"/createGraph":createGraph,
		"/graph_1.json": graph1Handler,
		"/dataset1.json": dataset1Handler,
		"/dataset2.json": dataset2Handler,
	}

	//setup routes	
	setupRoutes(routes, mux)
	
	//listen and serve on this port using this mux, log on end
	log.Fatal(http.ListenAndServe(listenAddr, mux))
}