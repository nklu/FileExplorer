package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type Node struct {
	children []*Node
	data     interface{}
}

func main() {
	dir := "c:\\inetpub"
	node := &Node{}

	err := walk(dir, node, getFileData)
	if err != nil {
		fmt.Println(err)
	} else {
		printNode(node)
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
			childNode.data = fnData(fileInfo, childErr)
		}
		if childNode != nil {
			node.children = append(node.children, childNode)
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

	if dataMap, mapOk := node.data.(map[string]interface{}); mapOk {
		if name, nameOk := dataMap["Name"]; nameOk {
			fmt.Println(name)
		}
	}

	if node.children == nil {
		return
	}

	for _, child := range node.children {
		printNode(child)
	}
}
