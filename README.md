# MindScript
Writing a compiler in Golang for fictional MindScript programming language based on the OSes in the culture series.

## Parts of a compiler
1. Lexical Analysis
    Also known as "scanning". We should break the input into a series of tokens. A token is a sequence of characters that can be treated as a single logical entity.
2. Syntax Analysis
3. Semantic Analysis
4. Intermediate Code Generation
5. Code Optimization
6. Code Generation

# Usage
```bash
make

./bin/msc -i ./examples/example.ms
```

Output:
```
Input file:  ./examples/example.ms
Output file:  ./examples/example.mind
{Type:KEYWORD Literal:agent}
{Type:KEYWORD Literal:DataProcessor}
{Type:LBRACE Literal:{}
{Type:KEYWORD Literal:goal}
```

# References
- https://www.geeksforgeeks.org/phases-of-a-compiler/
- https://github.com/kitasuke/monkey-go