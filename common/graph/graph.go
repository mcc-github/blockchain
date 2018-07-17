/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package graph


type Vertex struct {
	Id        string
	Data      interface{}
	neighbors map[string]*Vertex
}


func NewVertex(id string, data interface{}) *Vertex {
	return &Vertex{
		Id:        id,
		Data:      data,
		neighbors: make(map[string]*Vertex),
	}
}



func (v *Vertex) NeighborById(id string) *Vertex {
	return v.neighbors[id]
}


func (v *Vertex) Neighbors() []*Vertex {
	var res []*Vertex
	for _, u := range v.neighbors {
		res = append(res, u)
	}
	return res
}



func (v *Vertex) AddNeighbor(u *Vertex) {
	v.neighbors[u.Id] = u
	u.neighbors[v.Id] = v
}
