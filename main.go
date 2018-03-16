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
	Data     map[string]interface{}
	Children []*Node
}

type WalkResult struct {
	Node  *Node
	Error error
	Size  int64
}

func main() {
	dir := os.Args[1]
	fileName := os.Args[2]

	result := walk(dir, getFileData)
	if result.Error != nil {
		fmt.Println(result.Error)
	} else {
		//printNode(node)
		b, errJSON := getJSON(result.Node)
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

func walk(dir string, fnData func(os.FileInfo) (map[string]interface{}, int64)) *WalkResult {
	walkResult := &WalkResult{}

	if dir == "" {
		walkResult.Error = errors.New("Directory is empty")
		return walkResult
	}

	info, err := ioutil.ReadDir(dir)
	if err != nil {
		walkResult.Error = err
		return walkResult
	}
	walkResult.Node = &Node{}

	for _, fileInfo := range info {
		var childResult *WalkResult
		var childSize int64
		var childData map[string]interface{}
		if fnData != nil {
			childData, childSize = fnData(fileInfo)
		}
		if fileInfo.IsDir() {
			childResult = walk(path.Join(dir, fileInfo.Name()), fnData)
			if childData != nil {
				childData["Size"] = childResult.Size
			}
			childResult.Node.Data = childData
		} else {
			childResult = &WalkResult{Node: &Node{Data: childData}}
		}
		if childResult.Node != nil {
			walkResult.Node.Children = append(walkResult.Node.Children, childResult.Node)
		}
		walkResult.Size += childSize
	}

	if walkResult.Node.Data == nil {
		walkResult.Node.Data = map[string]interface{}{}
		walkResult.Node.Data["Name"] = dir
		walkResult.Node.Data["Size"] = walkResult.Size
	}

	return walkResult
}

func getFileData(info os.FileInfo) (map[string]interface{}, int64) {
	if info == nil {
		return nil, 0
	}

	ret := map[string]interface{}{}

	ret["Name"] = info.Name()
	ret["IsDir"] = info.IsDir()
	size := info.Size()
	ret["Size"] = size

	return ret, size
}

func printNode(node *Node) {
	if node == nil {
		return
	}

	if name, nameOk := node.Data["Name"]; nameOk {
		fmt.Println(name)
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
