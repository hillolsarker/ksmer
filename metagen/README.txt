How to run:

Require Python 2.7.x

Input: 
- csv file: 
Column 1: k-mer ID
Column G_1, G_2,...., G_m: Genome 1, 2,...., m
Last column: b 
n rows for number of k-mer
- Cell: number of k-mer in G_i

Output:
- logfile.txt: contains running time of each step.
- sol.csv: solution for x_1, x_2,...,xm (and y_1,y_2,....)

1. metagen_lp.py
- Require Gurobi (http://www.gurobi.com/)
- Command: python metagen_lp.py inputfile.csv logfile.txt sol.csv

2. metagen_le.py
- Require Gurobi (http://www.gurobi.com/)
- Command: python metagen_le.py inputfile.csv logfile.txt sol.csv

3. metagen_l1a.py
- Require CVXOPT (http://cvxopt.org/)
- Command: python metagen_l1a.py inputfile.csv logfile.txt sol.csv
