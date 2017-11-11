# Request count

The task is to design a module that counts the number of requests in a given time interval. The module should be thread-safe and efficient as possible.

The task also necessitates a garbage collection strategy.

## Solution

The solution presented here stores timestamp counts in an ordered queue and garbage collects using binary search.

### Complexities

#### Space
The space complexity of the algorithm is `O(i)` where `i` is the number of potential timestamps in the interval.

#### Time
The time complexity of 

* `Inc(amount uint)`
    * best case `O(1)`, when there is no garbage collection
    * average case `O(1/i * log2(n))` if calls to `Inc()` are on different timestamps
    * average case is less if `Inc()` is called frequently and timestamps overlap
    * worst case `O(log2(n))`, when there is garbage collection
* `Count() uint64`
    * best case `O(log2(n))`
    * average case is `O(n)` if `Count()` is always called before the last stored timestamp reaches its TTL
    * average case can be far less than `O(n)` if `Count()` is sometimes called after the last stored timestamp reaches its TTL, approaching `O(log2(n))` the longer `Count()` is not called
    * worst case is `O(n)`
