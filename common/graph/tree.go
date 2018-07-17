/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package graph



type Iterator interface {
	
	
	Next() *TreeVertex
}


type TreeVertex struct {
	Id          string        
	Data        interface{}   
	Descendants []*TreeVertex 
	Threshold   int           
}


func NewTreeVertex(id string, data interface{}, descendants ...*TreeVertex) *TreeVertex {
	return &TreeVertex{
		Id:          id,
		Data:        data,
		Descendants: descendants,
	}
}


func (v *TreeVertex) IsLeaf() bool {
	return len(v.Descendants) == 0
}



func (v *TreeVertex) AddDescendant(u *TreeVertex) *TreeVertex {
	v.Descendants = append(v.Descendants, u)
	return u
}


func (v *TreeVertex) ToTree() *Tree {
	return &Tree{
		Root: v,
	}
}



func (v *TreeVertex) Find(id string) *TreeVertex {
	if v.Id == id {
		return v
	}
	for _, u := range v.Descendants {
		if r := u.Find(id); r != nil {
			return r
		}
	}
	return nil
}



func (v *TreeVertex) Exists(id string) bool {
	return v.Find(id) != nil
}


func (v *TreeVertex) Clone() *TreeVertex {
	var descendants []*TreeVertex
	for _, u := range v.Descendants {
		descendants = append(descendants, u.Clone())
	}
	copy := &TreeVertex{
		Id:          v.Id,
		Descendants: descendants,
		Data:        v.Data,
	}
	return copy
}



func (v *TreeVertex) replace(id string, r *TreeVertex) {
	if v.Id == id {
		v.Descendants = r.Descendants
		return
	}
	for _, u := range v.Descendants {
		u.replace(id, r)
	}
}


type Tree struct {
	Root *TreeVertex
}



func (t *Tree) Permute() []*Tree {
	return newTreePermutation(t.Root).permute()
}



func (t *Tree) BFS() Iterator {
	return newBFSIterator(t.Root)
}

type bfsIterator struct {
	*queue
}

func newBFSIterator(v *TreeVertex) *bfsIterator {
	return &bfsIterator{
		queue: &queue{
			arr: []*TreeVertex{v},
		},
	}
}



func (bfs *bfsIterator) Next() *TreeVertex {
	if len(bfs.arr) == 0 {
		return nil
	}
	v := bfs.dequeue()
	for _, u := range v.Descendants {
		bfs.enqueue(u)
	}
	return v
}


type queue struct {
	arr []*TreeVertex
}

func (q *queue) enqueue(v *TreeVertex) {
	q.arr = append(q.arr, v)
}

func (q *queue) dequeue() *TreeVertex {
	v := q.arr[0]
	q.arr = q.arr[1:]
	return v
}
