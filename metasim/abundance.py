import sys
import os

# taxon profile creation -- creates 5 taxon profiles with variance in abundance scores

count2 = 1
count3 = 1
count4 = 1
count5 = 1

numOfFastaFiles = int(sys.argv[1])

directory = str(sys.argv[numOfFastaFiles + 2]).split('.')

if os.path.isdir(str(directory[0])) == False:
    os.mkdir(str(directory[0]))

for x in range(2,numOfFastaFiles + 2):
    file = open(str(sys.argv[x]),'r')
    
    if x == 2:
        fileOutput = open(str(sys.argv[numOfFastaFiles + 2]),'w')
        fileOutput2 = open('2-' + str(sys.argv[numOfFastaFiles + 2]),'w')
        fileOutput3 = open('3-' + str(sys.argv[numOfFastaFiles + 2]),'w')
        fileOutput4 = open('4-' + str(sys.argv[numOfFastaFiles + 2]),'w')
        fileOutput5 = open('5-' + str(sys.argv[numOfFastaFiles + 2]),'w')
    else:
        fileOutput = open(str(sys.argv[numOfFastaFiles + 2]),'a')
        fileOutput2 = open('2-' + str(sys.argv[numOfFastaFiles + 2]),'a')
        fileOutput3 = open('3-' + str(sys.argv[numOfFastaFiles + 2]),'a')
        fileOutput4 = open('4-' + str(sys.argv[numOfFastaFiles + 2]),'a')
        fileOutput5 = open('5-' + str(sys.argv[numOfFastaFiles + 2]),'a')

    for line in file:
        if line[0] == '>':
            fileOutput.write('1 name "' + line[1:].strip() + '"\n')
            fileOutput2.write(str(count2) + ' name "' + line[1:].strip() + '"\n')
            fileOutput3.write(str(count3) + ' name "' + line[1:].strip() + '"\n')
            fileOutput4.write(str(count4) + ' name "' + line[1:].strip() + '"\n')
            fileOutput5.write(str(count5) + ' name "' + line[1:].strip() + '"\n')
            count2 = count2 + 1
            count3 = count3 * 2
            count4 = count4 * 3
            count5 = count5 * 4

    file.close()
    fileOutput.close()
    fileOutput2.close()
    fileOutput3.close()
    fileOutput4.close()
    fileOutput5.close()

# runs the Metasim program with each taxon profile

for y in range(2,numOfFastaFiles + 2):
    os.system('./Metasim cmd --add-files ' + str(sys.argv[y]))

os.system('./Metasim cmd -d ' + directory[0] + ' --454 -r5000 ' + str(sys.argv[numOfFastaFiles + 2]))
os.system('./Metasim cmd -d ' + directory[0] + ' --454 -r5000 2-' + str(sys.argv[numOfFastaFiles + 2]))
os.system('./Metasim cmd -d ' + directory[0] + ' --454 -r5000 3-' + str(sys.argv[numOfFastaFiles + 2]))
os.system('./Metasim cmd -d ' + directory[0] + ' --454 -r5000 4-' + str(sys.argv[numOfFastaFiles + 2]))
os.system('./Metasim cmd -d ' + directory[0] + ' --454 -r5000 5-' + str(sys.argv[numOfFastaFiles + 2]))

# delete the taxon profiles after simulation

# os.remove(str(sys.argv[numOfFastaFiles + 2]))
# os.remove('2-' + str(sys.argv[numOfFastaFiles + 2]))
# os.remove('3-' + str(sys.argv[numOfFastaFiles + 2]))
# os.remove('4-' + str(sys.argv[numOfFastaFiles + 2]))
# os.remove('5-' + str(sys.argv[numOfFastaFiles + 2]))
