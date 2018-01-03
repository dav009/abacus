<img src="https://raw.githubusercontent.com/dav009/abacus/master/199506.svg?sanitize=true"  width="150px" >

# Abacus

Abacus let you count item frequencies in big datasets with a fixed amount of memory.

Unlike a regular counter it trades off accuracy for memory.
This is useful for particular tasks, for example in NLP/ML related tasks you might want to count millions of items
however approximate counts are good enough.

Example:

```go
counter := abacus.New(maxMemoryMB=10) # abacus will use max 10MB to store your counts
counter.Update([]string{"item1", "item2", "item2"})
counter.Counts("item1") # 1 , counts for "item1"
counter.Total() # 3 ,Total number of counts (sum of counts of all elements)
counter.Cardinality() # 2 , How many different items are there?
```

Abacus lets you define how much memory you want to use and you go from there counting items.
Of course there are some limitations, and if you set the memory threshold too low, you might get innacurate counts.

## Benchmarks

- Counting bigrams (words) from [Wiki corpus](http://www.cs.upc.edu/~nlp/wikicorpus/).
- Compared memory and accuracy of `Abacus` vs using a `map[string]int`


Corpus Data Structure Used Memory Accuracy

| Corpus  | Data Structure  | Used Memory     | Accuracy  |
|---------|-----------------|-----------------|-----------|
| Half of Wiki corpus (English)   | Abacus (1000MB) |  1.75GB    | 96%  |
| Half of Wiki corpus (English)   | Map       |  3.3GB    | 100%  |
| Complete Wiki corpus (English)  | Abacus (2200MB) |  3.63GB    | 98%  |
| Complete Wiki corpus (English)  | Abacus (500MB) |   741MB   | 15%  |
| Complete Wiki corpus (English)  | Map       |  10.46GB    | 100%  |

Note: This is me playing with Golang again, heavily based on [Bounter](https://github.com/RaRe-Technologies/bounter)




## Under the hood

- Count–min sketch
- HyperLogLog algorithm 


Icon made by (free-icon)[https://www.flaticon.com/free-icon/] from www.flaticon.com 
