# What is this?

My playground to learn Go, machine learning and prime factorization.

# Outcomes

* Learned some Go
   * Go routines, channels and waitGroups makes it easy to implement multi-threaded workloads
   * Performance tuning with `go test -cpuprofile cpu.prof` and `go tool pprof -http localhost:19123 cpu.prof` is also very easy
* Learned some Machine Learning
  * Focused on scikit-learn which seems the most easier tool to start on
  * Normalization of data with scaling to make raw data processable by ordinary python tools
  * That's always good to start by evaluating features/labels correlation with 'seaborn.pairplot'
  * Jupyter notebooks are very useful for quick iteration
* Learned about semi-primes factorization
  * Initially approach to use machine learning to estimate prime factors sum is not useful:
    * Although neural-network prediction score can be high, greater than 0.99, after de-normalization the actual error is very high
    * Trying to find the real prime factors sum based on neural-network prediction is way slower than using other methods
  * About the Fermat's factorization method
    * Tried a different approach, using only algorithms, to find the prime factor difference with a good performance
    * Afterwards re-checked Wikipedia and found that I just implemented Fermat's factorization :/