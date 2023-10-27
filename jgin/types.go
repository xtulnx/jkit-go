package jgin

import (
	"strconv"
	"strings"
)

// JIDs 用逗号分隔的整数列表
type JIDs string

func (id JIDs) GetInt() []int {
	ss1 := strings.Split(string(id), ",")
	id2 := make([]int, 0, len(ss1))
	for _, s := range ss1 {
		if n, ok := strconv.ParseInt(s, 10, 32); ok == nil {
			id2 = append(id2, int(n))
		}
	}
	return id2
}

func (id *JIDs) FromInt(ids []int) {
	elems := make([]string, len(ids))
	for i, v := range ids {
		elems[i] = strconv.Itoa(v)
	}
	id.FromTag(elems)
}

func (id *JIDs) FromID(ids []uint) {
	elems := make([]string, len(ids))
	for i, v := range ids {
		elems[i] = strconv.Itoa(int(v))
	}
	id.FromTag(elems)
}

func (id *JIDs) FromTag(elems []string) {
	if len(elems) == 0 {
		return
	}
	*id = JIDs(strings.Join(elems, ","))
}

func (id JIDs) GetID() []uint {
	if id == "" {
		return nil
	}
	ss1 := strings.Split(string(id), ",")
	id2 := make([]uint, 0, len(ss1))
	for _, s := range ss1 {
		if n, ok := strconv.ParseUint(s, 10, 32); ok == nil {
			id2 = append(id2, uint(n))
		}
	}
	return id2
}

func (id JIDs) GetIDUnique() []uint {
	if id == "" {
		return nil
	}
	ss1 := strings.Split(string(id), ",")
	m1 := make(map[uint]struct{})
	id2 := make([]uint, 0, len(ss1))
	for _, s := range ss1 {
		if s == "" {
			continue
		}
		if n, ok := strconv.ParseUint(s, 10, 32); ok == nil {
			var k = uint(n)
			if _, ok := m1[k]; ok {
				continue
			}
			m1[k] = struct{}{}
			id2 = append(id2, k)
		}
	}
	return id2
}

func (id JIDs) GetTag() []string {
	if id == "" {
		return nil
	}
	ss1 := strings.Split(string(id), ",")
	id2 := make([]string, 0, len(ss1))
	for _, s := range ss1 {
		if s != "" {
			id2 = append(id2, s)
		}
	}
	return id2
}

func (id JIDs) GetTagUnique() []string {
	if id == "" {
		return nil
	}
	ss1 := strings.Split(string(id), ",")
	m1 := make(map[string]struct{})
	id2 := make([]string, 0, len(ss1))
	for _, s := range ss1 {
		if s == "" {
			continue
		}
		if _, ok := m1[s]; ok {
			continue
		}
		m1[s] = struct{}{}
		id2 = append(id2, s)
	}
	return id2
}

//////////////////////////////////////////////////////////////////

//////////////////////////////////////////////////////////////////

type KeyText struct {
	Key  string `json:"key"`
	Text string `json:"text"`
}

type IdText struct {
	ID   uint   `json:"id" form:"id"`
	Text string `json:"text" form:"text"`
}

type IdExt interface {
	GetID() uint
	GetVal() interface{}
}
