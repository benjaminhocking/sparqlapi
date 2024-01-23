package main

import (
	"fmt"
	"net/http"
	"log"
	"sparqlAPI/utils"
	"strings"
	"encoding/json"
)

const (
	filePath = "finance.ttl"
	writeFilePath = "graph_2.json"
)


type handlerFunc func(http.ResponseWriter, *http.Request)

func indexHandler(w http.ResponseWriter, req *http.Request){
	createGraph(w, req)
	fmt.Fprintf(w, utils.ReadFile("pages/index.html"))
}

func index1Handler(w http.ResponseWriter, req *http.Request){
	fmt.Fprintf(w, utils.ReadFile("pages/index_1.html"))
}

func createGraph(w http.ResponseWriter, req *http.Request){
	/*
	Want:
		Nodes:
			[
				{"name": <name>, "id": <unique_id>},...
			]
		Links:
			[
				{"source": <source>, "target": <target>, "type": <predicate>,}
			]
	*/
	//get resources
	resources := utils.GetResourceMap(filePath) //type: map[int]string
	nodes := utils.GetNodes(resources) //type: []map[string]interface
	//get relations between
	links := utils.GetLinks(nodes, filePath) //type: []map[string]interface
	total := map[string][]map[string]any{"nodes":nodes, "links":links}
	res, _ := json.Marshal(total)
	utils.WriteFile(writeFilePath, string(res))
	fmt.Println("new graph made")
}


func newRecHandler(w http.ResponseWriter, req *http.Request){
	textQuery := req.URL.Query()["text"]
	text := textQuery[0]
	t, err := utils.CreateTriple(text, filePath)
	if err != nil{
		fmt.Fprintf(w, "Invalid Record")
		return
	}
	utils.SaveTriple(t, filePath)
	fmt.Fprintf(w, "Record Saved")
}

func queryHandler(w http.ResponseWriter, req *http.Request){
	textQuery := req.URL.Query()["text"]
	text := strings.ReplaceAll(textQuery[0], "<br>", "\n")

	fmt.Println(text)
	res := utils.MakeQuery(text, filePath)
	utils.PrintRes(res)
	resStr := utils.StringRes(res)
	resp := map[string]string{"response":resStr}
	b, _ := json.Marshal(resp)
	fmt.Fprintf(w, string(b))
}

func dataHandler(w http.ResponseWriter, req *http.Request){
	cont := utils.ReadFile("data/data.json")
	fmt.Fprintf(w, cont)
}

func graphHandler(w http.ResponseWriter, req *http.Request){
	cont := utils.ReadFile("data/graph.json")
	fmt.Fprintf(w, cont)
}

func graph1Handler(w http.ResponseWriter, req *http.Request){
	cont := utils.ReadFile(writeFilePath)
	fmt.Fprintf(w, cont)
}

func setupRoutes(m map[string]handlerFunc, mux *http.ServeMux){
	for k,v := range m{
		mux.HandleFunc(k,v)
	}
}

func main(){
	//port to listen on
	listenAddr := ":8080"

	//create new mux
	mux := http.NewServeMux()

	//list routes
	routes := map[string]handlerFunc{
		"/": indexHandler,
		"/new_rec": newRecHandler,
		"/query": queryHandler,
		"/data.json": dataHandler,
		"/index_1": index1Handler,
		"/graph.json": graphHandler,
		"/createGraph":createGraph,
		"/graph_1.json": graph1Handler,
	}

	//setup routes	
	setupRoutes(routes, mux)
	
	//listen and serve on this port using this mux, log on end
	log.Fatal(http.ListenAndServe(listenAddr, mux))
}