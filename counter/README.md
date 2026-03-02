# Worker Sum Task

## Task Description

This is a short exercise intended to take about **10--15 minutes**.

Write a program that:

1.  **Reads lines from standard input (stdin).**
2.  Each line contains **decimal numbers separated by spaces**.
3.  The program should **dispatch lines to worker(s)**.
4.  Each worker is responsible for **calculating the sum of numbers in a
    single line**.
5.  After processing all input, the program should **output the total
    sum of all numbers** across all lines.

------------------------------------------------------------------------

## Example Input

    1 2 3 4
    1000 200 300
    1
    100000 2000 3000
    500 2 3 1

------------------------------------------------------------------------

## Expected Behavior

-   Every line is processed independently.
-   Workers compute partial sums.
-   The main program aggregates results and prints the final total.
