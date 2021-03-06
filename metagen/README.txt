How to run

Require Python 2.7.x

Input: 
- csv file: frequencies of k-mer in Microbial genomes and in Reads
  - Column 1: k-mer ID
  - Column G_1, G_2,...., G_m: Microbial genome 1, 2,...., m (matrix F)
  - Last column: b 
  - n rows for number of k-mer
- Cell: number of k-mer in G_i

Output:
- logfile.txt: contains running time of each step.
- sol.csv: solution for variables x_1, x_2,...,xm (and y_1,y_2,....y_n)

1. metagen_lp.py
- Linear programming: Find x and minimal value of y_1+y_2+....+y_n for which Fx+y=b. x is the abundance of the genomes in the sample, y is the abundance of other unknown genomes.
- Require Gurobi installation at least 5.6.3 version (http://www.gurobi.com/)
- Command: python metagen_lp.py inputfile.csv logfile.txt sol.csv

2. metagen_le.py
- Linear equation: solve for x in Fx=b
- Require Gurobi (http://www.gurobi.com/)
- Command: python metagen_le.py inputfile.csv logfile.txt sol.csv

3. metagen_l1a.py
- L1-approximation: find x to minimize |Fx-b|
- Require CVXOPT installation (http://cvxopt.org/).
- metagen_l1a.py calls the function l1(P, q) (http://cvxopt.org/examples/mlbook/l1.html)
- Command: python metagen_l1a.py inputfile.csv logfile.txt sol.csv
