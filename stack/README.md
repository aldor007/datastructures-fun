# simple database with transactions system

Data structure used: stack
# Usage

```bash
$ go run main.go
>> set a 2
>> get a  
2
>> count 2
1
>> begin
>> delete a
>> count 2
0
>> get a
empty
>> rollback
>> get a
2
>> begin
>> delete a
>> get a
empty
>> commit
>> get a
empty
```

# Tests

```bash
$ go test ./... -v
```