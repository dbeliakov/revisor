#include <iostream>
#include <vector>

size_t fib(size_t n)
{
    assert(n > 0);
    std::vector<size_t> res = {1, 1};
    for (size_t i = 2; i < n; ++i) {
        res.push_back(res[i - 1] + res[i - 2]);
    }
    return res[n - 1];
}

int main()
{
    std::cout << fib(6) << std::endl;
}
