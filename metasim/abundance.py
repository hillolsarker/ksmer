import sys
import os

# taxon profile creation -- creates 5 taxon profiles with variance in abundance scores

count2 = 1
count3 = 1
count4 = 1
count5 = 1
countSum2 = 0
countSum3 = 0
countSum4 = 0
countSum5 = 0

numOfFastaFiles = int(sys.argv[1])

wordCount = []
totalWordCount = 0
directory = str(sys.argv[numOfFastaFiles + 2]).split('.')

if os.path.isdir(str(directory[0])) == False:
    os.mkdir(str(directory[0]))

for x in range(2,numOfFastaFiles + 2):
    file = open(str(sys.argv[x]),'r')
    temp, temp2 = str(os.popen("wc -m " + str(sys.argv[x])).read()).split()
    wordCount.append(int(temp))
    
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
        count2 = count2 + 1
        count3 = count3 * 2
        count4 = count4 * 3
        count5 = count5 * 4
        countSum2 = countSum2 + count2
        countSum3 = countSum3 + count3
        countSum4 = countSum4 + count4
        countSum5 = countSum5 + count5



    for line in file:
        if line[0] == '>':
            fileOutput.write('1 name "' + line[1:].strip() + '"\n')
            fileOutput2.write(str(count2) + ' name "' + line[1:].strip() + '"\n')
            fileOutput3.write(str(count3) + ' name "' + line[1:].strip() + '"\n')
            fileOutput4.write(str(count4) + ' name "' + line[1:].strip() + '"\n')
            fileOutput5.write(str(count5) + ' name "' + line[1:].strip() + '"\n')


    file.close()
    fileOutput.close()
    fileOutput2.close()
    fileOutput3.close()
    fileOutput4.close()
    fileOutput5.close()

for z in range(0,len(wordCount)):
    totalWordCount =  totalWordCount + wordCount[z]
countSum1 = totalWordCount
countSum2 = countSum2 + totalWordCount
countSum3 = countSum3 + totalWordCount
countSum4 = countSum4 + totalWordCount
countSum5 = countSum5 + totalWordCount

# runs the Metasim program with each taxon profile

for y in range(2,numOfFastaFiles + 2):
    os.system('./Metasim cmd --add-files ' + str(sys.argv[y]))

os.system('./Metasim cmd -d ' + directory[0] + ' --454 -r ' + str(countSum1) + ' -f 100 -t 0 ' + str(sys.argv[numOfFastaFiles + 2]))
os.system('./Metasim cmd -d ' + directory[0] + ' --454 -r ' + str(countSum2) + ' -f 100 -t 0 2-' + str(sys.argv[numOfFastaFiles + 2]))
os.system('./Metasim cmd -d ' + directory[0] + ' --454 -r ' + str(countSum3) + ' -f 100 -t 0 3-' + str(sys.argv[numOfFastaFiles + 2]))
os.system('./Metasim cmd -d ' + directory[0] + ' --454 -r ' + str(countSum4) + ' -f 100 -t 0 4-' + str(sys.argv[numOfFastaFiles + 2]))
os.system('./Metasim cmd -d ' + directory[0] + ' --454 -r ' + str(countSum5) + ' -f 100 -t 0 5-' + str(sys.argv[numOfFastaFiles + 2]))

# delete the taxon profiles after simulation

# os.remove(str(sys.argv[numOfFastaFiles + 2]))
# os.remove('2-' + str(sys.argv[numOfFastaFiles + 2]))
# os.remove('3-' + str(sys.argv[numOfFastaFiles + 2]))
# os.remove('4-' + str(sys.argv[numOfFastaFiles + 2]))
# os.remove('5-' + str(sys.argv[numOfFastaFiles + 2]))
