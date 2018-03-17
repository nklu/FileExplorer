package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
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

	syncHelp := &syncHelp{Wg: &sync.WaitGroup{}, Mtx: &sync.Mutex{}}

	result, err := walk(dir, nil, syncHelp)

	syncHelp.Wg.Wait()

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

func walk(dir string, fnData func(os.FileInfo) *map[string]interface{}, syncHelp *syncHelp) (node *Node, err error) {
	syncHelp.Wg.Add(1)
	defer syncHelp.Wg.Done()

	info, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	node = &Node{Name: dir}

	for _, fileInfo := range info {
		var childNode *Node
		if fileInfo.IsDir() {
			fullPath := path.Join(dir, fileInfo.Name())
			go func() {
				childNode, _ = walk(fullPath, fnData, syncHelp)
			}()
			if childNode == nil {
				continue
			}
		} else {
			childNode = &Node{}
		}

		childNode.GetFileBaseData(fileInfo)
		if fnData != nil {
			childNode.Data = fnData(fileInfo)
		}
		node.Size += childNode.Size
		syncHelp.Mtx.Lock()
		node.Children = append(node.Children, childNode)
		syncHelp.Mtx.Unlock()
	}

	return
}

func (node *Node) GetFileBaseData(info os.FileInfo) {
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
