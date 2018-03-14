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

func main() {
	dir := os.Args[1]
	fileName := os.Args[2]

	node, err, _ := walk(dir, getFileData)
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

func walk(dir string, fnData func(os.FileInfo, error) (map[string]interface{}, int64)) (node *Node, err error, size int64) {

	if dir == "" {
		return node, errors.New("Directory is empty"), size
	}

	info, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	node = &Node{}

	for _, fileInfo := range info {
		var childNode *Node
		var childErr error
		var childSize int64
		var childData map[string]interface{}
		if fnData != nil {
			childData, childSize = fnData(fileInfo, childErr)
		}
		if fileInfo.IsDir() {
			childNode, childErr, childSize = walk(path.Join(dir, fileInfo.Name()), fnData)
			if childData != nil {
				childData["Size"] = childSize
			}
			childNode.Data = childData
		} else {
			childNode = &Node{Data: childData}
		}
		if childNode != nil {
			node.Children = append(node.Children, childNode)
		}
		size += childSize
	}

	if node.Data == nil {
		node.Data = map[string]interface{}{}
		node.Data["Name"] = dir
		node.Data["Size"] = size
	}

	return
}

func getFileData(info os.FileInfo, err error) (map[string]interface{}, int64) {
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
