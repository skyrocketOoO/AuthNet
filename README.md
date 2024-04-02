# Authnet

Authnet is a graph based auth service.


## Features

- [x] Get all edges
- [x] Create edge
- [x] Delete edge
- [x] Delete edges by conditions
- [x] Batch create or delete operations
- [x] Get all namespaces
- [x] Check relation
- [x] Get shortest path
- [x] Get all paths
- [x] Get all object relations
- [x] Get all subject relations
- [x] Get tree

## Relation

The `Relation` struct represents a relationship like edge in DAG between objects and subjects. It is defined as follows:

```go
// This means: Subject has a relation on Object
type Relation struct {
    ObjectNamespace  string
    ObjectName       string 
    Relation         string 
    SubjectNamespace string 
    SubjectName      string 
    SubjectRelation  string 
}
```

## How to use

1. Run postgres on docker(without docker, see ./docker-compose.yaml to get config)

    ```bash
    docker compose up -d postgres
    ```

2. Run the main server

    ```bash
    go run .
    ```

## Example

[HRBAC](https://github.com/skyrocketOoO/hrbac/tree/main)

## Reserved words

%

## Development benchmark

[Link](https://docs.google.com/spreadsheets/d/1qZiRE_kkno1mM0LzWiUnvX4cuYQRnep2NcNb4fPud-k/edit#gid=0)

## Abbreviation

- Namespace -> Ns
- Relation -> Rel
- Object -> Obj
- Subject -> Sbj
- Condition -> Cond
- Authority -> Auth

## Something...
### Why only store edges?
Store only edges can reduce the storage space usage, our app only concern about
who has access to the other instead of the vertex's infomation. So we can focus
on access management to reduce the other requirement