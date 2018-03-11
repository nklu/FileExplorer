package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type Node struct {
	Children []*Node
	Data     interface{}
}

func main() {
	dir := os.Args[1]
	fileName := os.Args[2]
	node := &Node{}

	err := walk(dir, node, getFileData)
	if err != nil {
		fmt.Println(err)
	} else {
		//printNode(node)
		b, errJSON := getJSON(node)
		fmt.Println(b)
		if errJSON != nil {
			panic(errJSON)
		}
		errFile := writeFileToCurrentDir(b, fileName)
		if errFile != nil {
			panic(errFile)
		}
	}
}

func walk(dir string, node *Node, fnData func(os.FileInfo, error) interface{}) (err error) {
	if node == nil {
		return errors.New("Cannot Be Nil")
	}
	if dir == "" {
		return errors.New("Directory is empty")
	}

	info, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, fileInfo := range info {
		childNode := &Node{}
		var childErr error
		if fileInfo.IsDir() {
			childErr = walk(path.Join(dir, fileInfo.Name()), childNode, fnData)
		}
		if fnData != nil {
			childNode.Data = fnData(fileInfo, childErr)
		}
		if childNode != nil {
			node.Children = append(node.Children, childNode)
		}
	}
	return
}

func getFileData(info os.FileInfo, err error) interface{} {
	if info == nil {
		return nil
	}

	ret := map[string]interface{}{}

	ret["Name"] = info.Name()
	ret["Size"] = info.Size()
	ret["IsDir"] = info.IsDir()

	return ret
}

func printNode(node *Node) {
	if node == nil {
		return
	}

	if dataMap, mapOk := node.Data.(map[string]interface{}); mapOk {
		if name, nameOk := dataMap["Name"]; nameOk {
			fmt.Println(name)
		}
	}

	if node.Children == nil {
		return
	}

	for _, child := range node.Children {
		printNode(child)
	}
}

func getJSON(node *Node) (b []byte, err error) {
	b, err = json.Marshal(node)
	return
}

func writeFileToCurrentDir(json []byte, fileName string) (err error) {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	err = ioutil.WriteFile(path.Join(dir, fileName), json, 0644)
	return
}
