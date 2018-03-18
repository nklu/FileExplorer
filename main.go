package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type Node struct {
	Size     int64
	IsDir    bool
	Name     string
	Modified time.Time
	Data     *map[string]interface{}
	Children []*Node
}

type syncHelp struct {
	Wg  *sync.WaitGroup
	Mtx *sync.Mutex
}

func main() {
	dir := os.Args[1]
	fileName := os.Args[2]
	var doSync bool
	var argErr error
	if doSync, argErr = strconv.ParseBool(os.Args[3]); argErr != nil {
		doSync = false
	}

	var result *Node
	var err error

	if doSync {
		result, err = walkSync(dir, nil)
	} else {
		var wg sync.WaitGroup
		result, err = walk(dir, &wg, nil)
		wg.Wait()
		fmt.Println("Async")
	}

	if err != nil {
		panic(err)
	}

	b, errJSON := getJSON(result)
	if errJSON != nil {
		panic(errJSON)
	}

	errFile := writeFileToCurrentDir(b, fileName)
	if errFile != nil {
		panic(errFile)
	}

	fmt.Println("Successful")
}

var mtx sync.Mutex

func walk(dir string, wg *sync.WaitGroup, fnData func(os.FileInfo) *map[string]interface{}) (node *Node, err error) {

	info, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	node = &Node{Name: dir}
	for _, fileInfo := range info {
		if fileInfo.IsDir() {
			wg.Add(1)
			go func(curNode *Node, curInfo os.FileInfo) {
				fullPath := path.Join(dir, curInfo.Name())
				childNode, _ := walk(fullPath, wg, fnData)
				addNodeAndData(curNode, childNode, curInfo)
				wg.Done()
			}(node, fileInfo)
		} else {
			childNode := &Node{}
			addNodeAndData(node, childNode, fileInfo)
		}
	}
	return
}

func addNodeAndData(parentNode *Node, childNode *Node, fileInfo os.FileInfo) {
	childNode.getFileBaseData(fileInfo)
	mtx.Lock()
	parentNode.Size += childNode.Size
	parentNode.Children = append(parentNode.Children, childNode)
	mtx.Unlock()
}

func (node *Node) getSize(mtx *sync.Mutex) (size int64) {
	mtx.Lock()
	size = node.Size
	mtx.Unlock()
	return
}

func walkSync(dir string, fnData func(os.FileInfo) *map[string]interface{}) (node *Node, err error) {

	info, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	node = &Node{Name: dir}

	for _, fileInfo := range info {
		var childNode *Node
		if fileInfo.IsDir() {
			fullPath := path.Join(dir, fileInfo.Name())
			childNode, _ = walkSync(fullPath, fnData)
			if childNode == nil {
				continue
			}
		} else {
			childNode = &Node{}
		}

		childNode.getFileBaseData(fileInfo)
		if fnData != nil {
			childNode.Data = fnData(fileInfo)
		}
		node.Size += childNode.Size
		node.Children = append(node.Children, childNode)
	}

	return
}

func (node *Node) getFileBaseData(info os.FileInfo) {
	node.IsDir = info.IsDir()
	node.Name = info.Name()
	node.Modified = info.ModTime()
	if !node.IsDir {
		node.Size = info.Size()
	}
}

func printNode(node *Node) {
	if node == nil {
		return
	}
	fmt.Println(node.Name)

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
