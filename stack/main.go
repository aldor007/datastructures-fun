package main

import (
	"fmt"
	"os"
	"strings"
	"bufio"
)

type DataStatus int

const (
	SET    DataStatus = 0
	DELETE DataStatus = 1
)

type TransactionData struct {
	status DataStatus
	value  string
}

type Transaction struct {
	data    map[string]TransactionData
	counter map[string]int
	prev    *Transaction
}

func NewTransaction() *Transaction{
	return &Transaction{
		data: make(map[string]TransactionData),
		counter: make(map[string]int),
		prev: nil,
	}
}

type TransactionStack struct {
	top    *Transaction
}

func (tr *TransactionStack) Push(t *Transaction) {
	if tr.top == nil {
		tr.top = t
	} else {
		oldTop := tr.top
		tr.top = t
		tr.top.prev = oldTop
	}
}

func (tr *TransactionStack) Top() *Transaction {
	return tr.top
}

func (tr *TransactionStack) Pop() *Transaction {
	if tr.top == nil {
		return nil
	}
	top := tr.top
	prevTop := top.prev
	top.prev = nil
	tr.top = prevTop

	return top
}

type DB struct {
	tr      *TransactionStack
	data    map[string]string
	counter map[string]int
}

func NewDB() *DB {
	return &DB{
		tr:      &TransactionStack{top: nil},
		data:    map[string]string{},
		counter: make(map[string]int),
	}
}

func (d *DB) Get(key string) string {
	if tr := d.tr.Top(); tr != nil {
		if v, ok := tr.data[key]; ok {
			if v.status == DELETE {
				return "" 
			} else {
				return v.value
			}
		}
		for tr.prev != nil {
			tr = tr.prev
			if v, ok := tr.data[key]; ok {
				if v.status != DELETE {
					return v.value
				} else {
					return ""
				}
			}

		}
	}

	return d.data[key]
}

func (d *DB) Set(key, value string) {
	if tr := d.tr.Top(); tr != nil {
		if v, ok := tr.counter[value]; ok {
			tr.counter[value] = v + 1
		} else {
			tr.counter[value] = 1
		}
		tr.data[key] = TransactionData{
			value:  value,
			status: SET,
		}
	} else {
		d.data[key] = value
		if v, ok := d.counter[value]; ok {
			d.counter[value] = v + 1
		} else {
			d.counter[value] = 1
		}
	}
}

func (d *DB) Delete(key string) {
	if tr := d.tr.Top(); tr != nil {
		if v, ok := tr.data[key]; ok {
			if v.status != DELETE {
				v.status = DELETE
				tr.data[key] = v
				if v.value != "" {
					tr.counter[v.value] = tr.counter[v.value] - 1
					if tr.counter[v.value] < -1 {
						tr.counter[v.value] = -1
					}
				}
			}
		} else {
			tr.data[key] = TransactionData{
				status: DELETE,
				value: "",
			}
			if v.value != "" {
				tr.counter[v.value] = -1
			}
		}
	} else {
		if v, ok := d.data[key]; ok {
			d.counter[v] = d.counter[v] - 1
			if d.counter[v] < 0 {
				d.counter[v] = 0
			}
		}
		delete(d.data, key)
	}
}

func (d *DB) Count(value string) int {
	sum := 0
	if tr := d.tr.Top(); tr != nil {
		if v, ok := tr.counter[value]; ok {
			sum = sum + v
		} else {
			// maybe removed
			for _, v := range tr.data {
				if v.status == DELETE {
					sum = sum - 1
				}
			}
		}
		for tr.prev != nil {
			tr = tr.prev
			if v, ok := tr.counter[value]; ok {
				sum += v
			} else {
				for _, v := range tr.data {
					if v.status == DELETE {
						sum = sum - 1
					}
				}
			}
		}


	}
	
	sum += d.counter[value]

	return sum
}

func (d *DB) Begin() {
	tr := NewTransaction() 
	d.tr.Push(tr)
}

func (d *DB) Rollback() {
	d.tr.Pop()
}

func (d *DB) Commit() {
	top := d.tr.Pop()
	if top == nil {
		return
	}
	for k, v := range top.data {
		if top.prev != nil {
			top.prev.data[k] = v
		}
	}
	for k, v := range top.counter{
		if top.prev != nil {
			top.prev.counter[k] = v
		}
	}
	// bottom
	if  top.prev == nil {
		for k, v := range top.data {
			if v.status == DELETE {
				delete(d.data, k)
			} else {
				d.data[k] = v.value
			}
		}

		counter := make(map[string]int)
		for _, v := range d.data {
			if vC, ok := counter[v]; ok {
				counter[v] = vC + 1 
			} else {
				counter[v] = 1
			}
		}
		d.counter = counter

	}
}


func main() {
	db := NewDB()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		parts := strings.Split(line, " ")
		parts[0] = strings.ToLower(parts[0])
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		switch parts[0] {
		case "set":
			if len(parts) < 3 {
				fmt.Println("invalid params")
				continue
			}
			db.Set(parts[1], parts[2])
		case "get":
			if len(parts) < 2 {
				fmt.Println("invalid params")
				continue
			}
			v := db.Get(parts[1])
			if v == "" {
				fmt.Println("empty")
			} else {
				fmt.Println(v)
			}
		case "delete":
			if len(parts) < 2 {
				fmt.Println("invalid params")
				continue
			}
			db.Delete(parts[1])
		case "begin":
			db.Begin()
		case "rollback":
			db.Rollback()
		case "commit":
			db.Commit()
		case "count":
			fmt.Println(db.Count(parts[1]))
		case "quit":
			os.Exit(0)
		default:
			fmt.Println("unknown command")
			
		}
	}
}