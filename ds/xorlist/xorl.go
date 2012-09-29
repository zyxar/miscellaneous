package xorlist

import (
	"fmt"
	"sync"
	"unsafe"
)

type xnode struct {
	ptr  *xnode      "pointer"
	data interface{} "payload"
}

type XorList struct {
	head   *xnode
	tail   *xnode
	length int
	*sync.Mutex
}

func xor(p1, p2 *xnode) *xnode {
	return (*xnode)(unsafe.Pointer((uintptr(unsafe.Pointer(p1)) ^ uintptr(unsafe.Pointer(p2)))))
}

func New(data []interface{}) (*XorList, error) {
	length := len(data)
	var head, tail *xnode
	if length == 0 {
		head = nil
		tail = nil
		return &XorList{nil, nil, 0, new(sync.Mutex)}, nil
	} else {
		var ph, p, pt *xnode
		p = new(xnode)
		if p == nil {
			return nil, fmt.Errorf("insufficient memory.")
		}
		ph = p
		p.data = data[0]
		head = p
		pt = nil
		for i := 1; i < length; i++ {
			pt = new(xnode)
			pt.data = data[i]
			p.ptr = xor(ph, pt)
			ph = p
			p = pt
		}
		tail = p
		if pt != nil {
			p.ptr = xor(ph, pt)
		}
	}
	return &XorList{head, tail, length, new(sync.Mutex)}, nil
}

func (list *XorList) Append(data interface{}) error {
	list.Lock()
	defer list.Unlock()
	p := new(xnode)
	if p == nil {
		return fmt.Errorf("insufficient memory.")
	}
	p.data = data
	p1 := list.tail
	p0 := xor(p1, p1.ptr)
	p1.ptr = xor(p0, p)
	p.ptr = xor(p1, p)
	list.tail = p
	list.length++
	return nil
}

func (list *XorList) Trim() {
	list.Lock()
	defer list.Unlock()
	if list.length < 1 {
		return
	}
	p1 := xor(list.tail, list.tail.ptr)
	p0 := xor(list.tail, p1.ptr)
	p1.ptr = xor(p1, p0)
	list.length--
	list.tail = p1
}

func (list *XorList) Reverse() *XorList {
	list.Lock()
	defer list.Unlock()
	l := new(XorList)
	if l == nil {
		return nil
	}
	*l = *list
	l.head = list.tail
	l.tail = list.head
	return l
}

func (list *XorList) Len() int {
	list.Lock()
	defer list.Unlock()
	return list.length
}

func (list *XorList) Traverse() []interface{} {
	list.Lock()
	defer list.Unlock()
	if list.length == 0 {
		return nil
	}
	var p0, p1, p2 *xnode
	p0 = list.head
	p1 = p0
	ret := make([]interface{}, list.length)
	ret[0] = p1.data
	i := 1
	for p1 != list.tail {
		p2 = xor(p0, p1.ptr)
		ret[i] = p2.data
		p0 = p1
		p1 = p2
		i++
	}
	return ret
}

func (list *XorList) Get(num uint) interface{} {
	list.Lock()
	defer list.Unlock()
	if num >= uint(list.length) {
		return nil
	}
	return tr(list.head, num).data
}

func (list *XorList) Set(num uint, value interface{}) error {
	list.Lock()
	defer list.Unlock()
	if num >= uint(list.length) {
		return fmt.Errorf("Index out of range.")
	}
	tr(list.head, num).data = value
	return nil
}

func tr(head *xnode, num uint) *xnode {
	var p0, p1, p2 *xnode = head, head, nil
	for i := uint(0); i <= num; i++ {
		p2 = xor(p0, p1.ptr)
		p0 = p1
		p1 = p2
	}
	return p0
}

// func main() {
// 	a, err := New([]interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 44, 32})
// 	if err != nil {
// 		return
// 	}
// 	fmt.Println(a.Traverse())
// 	a.Append("data")
// 	fmt.Println(a.Reverse().Traverse())
// 	a.Trim()
// 	fmt.Println(a.Traverse())
// }
