#include <iostream>

size_t fib(size_t n)
{
    if (n == 1) {
        return 1;
    } else if (n == 2) {
        return 1;
    } else {
        return fib(n - 1) + fib(n - 2);
    }
}

int main()
{
    std::cout << fib(6) << std::endl;
}
