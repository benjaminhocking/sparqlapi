package utils

import (
	"strings"
	"errors"
	"os"
	"log"
	"bufio"
	"fmt"
	"slices"
)

type Triple struct{
	Subject string
	Predicate string
	Object string
}

type QueryTriple struct{
	Subject string
	Predicate string
	Object string
}

type QueryType int

func openDB(filePath string) *os.File {
	file, err := os.Open(filePath)
	if err != nil{
		log.Fatal(err)
	}
	return file
}


func SaveTriple(t Triple, filePath string){
	cont := []byte(t.Subject + " " + t.Predicate + " " + t.Object + " .\n")
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err !=nil{
		log.Fatal(err)
	}
	defer f.Close()
	if _, err := f.Write(cont); err != nil{
		log.Fatal(err)
	}
	if err := f.Close(); err != nil{
		log.Fatal(err)
	}
}

func HasRecord(s string, filePath string) bool{
	file := openDB(filePath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan(){
		lineText := scanner.Text()
		if lineText == s{
			return true
		}
	}
	return false
}

func ValidateNewRecord(s string, filePath string) bool{
	tokens := strings.Split(s, " ")
	if len(tokens)!=4{
		return false
	}
	if tokens[3]!="."{
		return false
	}
	if HasRecord(s, filePath){
		return false
	}
	return true
}

func getVars(s []string) []string {
	res := []string{}
	for _, v := range s{
		if strings.ContainsRune(v, '?'){
			res = append(res, v)
		}
	}
	return res
}

func CreateTriple(s string, filePath string) (Triple, error){
	if !ValidateNewRecord(s, filePath){
		return Triple{"","",""}, errors.New("invalid")
	}
	tokens := strings.Split(s, " ")
	t := Triple{tokens[0], tokens[1], tokens[2]}
	return t, nil
}

func getVarMap(vars []string, tokens []string) map[string]int{
	//work out where each of the vars appears in the triple, return map is of the form:
	/*	
		{
			"varName": [0,2],
			"varName1": [0,2],		
		}
	*/
	m := map[string]int{}
	for _, v := range vars{
		for j, t := range tokens{
			if v==t{
				m[v]=j
				//fmt.Printf(v + " %d\n", j)
			}
		}
	}
	return m
}

func getQueryTriple(s []string) QueryTriple{
	sub := s[2]
	pred := s[3]
	obj := s[4]
	qt := QueryTriple{sub, pred, obj}
	return qt
}

func query(qt QueryTriple, vars []string, varMap map[string]int, filePath string) []map[string]string{
	file := openDB(filePath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	//matches will be of the form:
	/*
	[
		{"?var1": "res1", "?var2": "res2"},
		{"?var1": "res1", "?var2": "res2"},
	]
	*/
	matches := []map[string]string{}

	for scanner.Scan(){
		lineText := scanner.Text()
		//fmt.Println(lineText)
		lineTokens := strings.Split(lineText, " ")
		if match(lineText, qt){
			//if this line matches the QueryTriple, we must collect the vars into a map and append this to the matches slice
			m := map[string]string{}
			for _, v := range vars{
				//vars will be of form: ?varName, and lineText will be of form: <sub> <pred> <obj> .
				m[strings.ReplaceAll(v, "?", "")]=lineTokens[varMap[v]]
			}
			matches = append(matches, m)
		}
	}
	return matches
}

func match(s string, qt QueryTriple) bool{
	//QueryTriple will contain the values that a match must have:
	//	for any ?var value, this can match anything
	tokens := strings.Split(s, " ")
	if !strings.ContainsRune(qt.Subject, '?') && qt.Subject!=tokens[0]{
		return false
	}
	if !strings.ContainsRune(qt.Predicate, '?') && qt.Predicate!=tokens[1]{
		return false
	}
	if !strings.ContainsRune(qt.Object, '?') && qt.Object!=tokens[2]{
		return false
	}
	return true
}

func PrintRes(res []map[string]string){
	for _, v := range res{
		for k, val := range v{
			fmt.Println(k, ":", val)
		}
	}
}
func StringRes(res []map[string]string) string{
	s := ""
	for _, v := range res{
		for k, val := range v{
			s += k + ": " + val + "\n"
		}
	}
	return s
}

func lineToTriple(line string) Triple{
	tokens := strings.Split(line, " ")
	t := Triple{tokens[0], tokens[1], tokens[2]}
	return t
}

func GetResourceMap(filePath string) map[int]string{
	file := openDB(filePath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	resources := []string{}
	m := map[int]string{}
	i := 0
	for scanner.Scan(){
		lineText := scanner.Text()
		t := lineToTriple(lineText)
		if !slices.Contains(resources, t.Subject){
			resources = append(resources, t.Subject)
			m[i] = t.Subject
			i+=1
		}
		if !slices.Contains(resources, t.Object){
			resources = append(resources, t.Object)
			m[i] = t.Object
			i+=1
		}
	}
	return m
}

func GetNodes(resources map[int]string) []map[string]any{
	nodes := []map[string]any{}
	for k, v := range resources{
		t := map[string]any{"name":v, "id":k}
		nodes = append(nodes, t)
	}
	return nodes
}

func getResourceId(resourceName string, nodes []map[string]any) int{
	for _, node := range nodes{
		if resourceName == node["name"]{
			r, _ := node["id"].(int)
			return r
		}
	}
	return -1
}

func GetLinks(nodes []map[string]any, filePath string) []map[string]any{
	file := openDB(filePath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	links := []map[string]any{}
	for scanner.Scan(){
		lineText := scanner.Text()
		t := lineToTriple(lineText)
		m := map[string]any{"source":getResourceId(t.Subject, nodes), "target":getResourceId(t.Object, nodes), "type":t.Predicate}
		links = append(links, m)
	}
	return links
}

func MakeQuery(s string, filePath string) []map[string]string{
	lines := strings.Split(s, "\n")
	line1 := lines[0]
	line2 := lines[1]


	tokensLine1 := strings.Split(line1, " ")
	tokensLine2 := strings.Split(line2, " ")

	vars := getVars(tokensLine1)

	varMap := getVarMap(vars, tokensLine2[2:])
	qt := getQueryTriple(tokensLine2)

	return query(qt, vars, varMap, filePath)
}