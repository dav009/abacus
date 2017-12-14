# Abacus

Abacus let you count  item frequencies in big datasets.
Unlike a regular counter it trades off accuracy for memory.
This is useful for particular tasks, for example in NLP related tasks you want to get estimates of counts.

Example:

```go
counter := abacus.New(maxMemoryMB=10) # abacus will use max 10MB to store your counts
counter.Update([]string{"item1", "item2", "item2"})
counter.Counts("item1") # 1
counter.Total() # 3
counter.Cardinality() # 2
```

Abacus lets you define how much memory you want to use and you go from there counting items.
Of course there are some limitations, and if you set the memory threshold too low, you might get innacurate counts.

## Benchmarks

Someday.. ;)

Note: This is me playing with Golang again, heavily based on [Bounter](https://github.com/RaRe-Technologies/bounter)



