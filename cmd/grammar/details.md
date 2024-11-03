go generate grammar -i=<input_file>

where:
- `<input_file>` refers to the file where the grammar is contained.

NOTES:
- The input file must be a file having the ".txt" extension.


Once the command is run, the `test` directory will contain the following files:

```
test
├── ast.go
├── node_type.go
├── node.go
├── lexer.go
└── parsing.go
```

Here:
- `ast.go`: A Go file containing the abstract syntax tree (AST) of the grammar. This file must be modified by users to suit their needs.
- `node_type.go`: A Go file containing the types of node. This file can be edited and/or adjusted to one's needs and use cases.
   - This file also contains a `go:generate` line for when users prefer using the `stringer` generation tool. If this is the case, then users must also run the `go generate` command themselves. In any other case, users are recommended to delete the `go:generate` line.
- `node.go`: A Go file containing the bare minimum implementation of the node type required for running the SlParser. As such, users must not edit this file as, otherwise, things may not work as expected.
   - If new functions are required (for one reason or another), then users should create these functions in another file and pass the node as the first parameter rather than adding new methods.
- `lexer.go`: A Go file containing the layout of the lexer.
   - This file must be modified by users to suit their needs.
- `parsing.go`: A Go file containing the actual parser implementation (as well as the parser's basic tree writer).
   - This file must not be touched for any reason.

Once the files are successfully generated, users should follow the following steps:
