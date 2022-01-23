

import sys
from typing import NamedTuple
import matplotlib.pyplot as plt
import numpy as np
import math

# base definitions
numberN = 13438310478517603073
sumEst = 7387124034+1234567

# the sum of primes has to be even
startNumOddFix = (sumEst % 2)
c1 = np.arange(startNumOddFix-0, 15000,  2)
c2 = np.arange(startNumOddFix-0, 15000,  6)
c3 = np.arange(startNumOddFix-0, 15000,  4)
c4 = np.arange(startNumOddFix-0, 15000, 10)

# solve([p1*p2 = numberN, p1+p2+c1=sumEst], [p1, p2]);
# p1(x) := -(sqrt(sumEst^2  - 2 * c1 * sumEst - 4 * numberN + c1^2 ) - sumEst + c1)/2;
# p2(x) :=  (sqrt(sumEst^2  - 2 * c1 * sumEst - 4 * numberN + c1^2 ) + sumEst - c1)/2;

def p1Function(c1):
     return -(math.sqrt(sumEst**2 - 2 * c1 * sumEst - 4 * numberN + c1**2 ) - sumEst + c1)/2
myP1Function = np.vectorize(p1Function)
p1 = myP1Function(c1)

def p2Function(c1):
     return (math.sqrt(sumEst**2 - 2 * c1 * sumEst - 4 * numberN + c1**2 ) + sumEst - c1)/2
myP2Function = np.vectorize(p2Function)
p2 = myP2Function(c1)

def diffFunction(c1):
    p1AsInt = round(p1Function(c1))
    p2AsInt = round(p2Function(c1))
    diffValue = numberN - (p1AsInt * p2AsInt)
    # if diffValue > 1274550:
    #     return -1274550
    return abs(diffValue)

from collections import namedtuple
MinValuePosition = namedtuple('MinValuePosition', ['value', 'x', 'index'])
# find first local min
for i in range(1, len(c1)-1):
    curValue = diffFunction(c1[i])
    if diffFunction(c1[i-1]) > diffFunction(c1[i]) and diffFunction(c1[i]) < diffFunction(c1[i+1]):
        firstMin = MinValuePosition(diffFunction(c1[i]), c1[i], i)
        break
# find second local min
for i in range(firstMin.index+1, len(c1)-1):
    curValue = diffFunction(c1[i])
    if diffFunction(c1[i-1]) > diffFunction(c1[i]) and diffFunction(c1[i]) < diffFunction(c1[i+1]):
        secondMin = MinValuePosition(diffFunction(c1[i]), c1[i], i)
        break

# find all local min
mins = []
for i in range(1, len(c1)-1):
    curValue = diffFunction(c1[i])
    if diffFunction(c1[i-1]) > diffFunction(c1[i]) and diffFunction(c1[i]) < diffFunction(c1[i+1]):
        mins.append(MinValuePosition(diffFunction(c1[i]), c1[i], i))

minsArray = np.array(mins)
print (firstMin)
print (secondMin)
# print(mins)
# abort

myDiffFunction = np.vectorize(diffFunction)
df1 = myDiffFunction(c1)
df2 = myDiffFunction(c2)
df3 = myDiffFunction(c3)
df4 = myDiffFunction(c4)

def sumFunction(c1):
    p1AsInt = round(p1Function(c1))
    p2AsInt = round(p2Function(c1))
    diffValue = p1AsInt + p2AsInt
    # if diffValue > 1274550:
    #     return -1274550
    return diffValue
mySumFunction = np.vectorize(sumFunction)
sf = mySumFunction(c1)

# setting the axes at the centre
fig = plt.figure()
ax = fig.add_subplot(1, 1, 1)
ax.spines['left'].set_position('center')
ax.spines['bottom'].set_position('zero')
ax.spines['right'].set_color('none')
ax.spines['top'].set_color('none')
ax.xaxis.set_ticks_position('bottom')
ax.yaxis.set_ticks_position('left')

# print(c1)
# print(c2[369])

# abort("cenas")

# plot the function
# plt.plot(c1,p1, 'r')
# plt.plot(c1,p2, 'x')
# plt.plot(c1,sf, 'r')
plt.plot(c1,df1, '.', label='df')
# plt.plot(c1,df1, label='df1')
# plt.plot(c3,df3, label='df3')
# plt.plot(c4,df4, label='df4')
plt.plot(minsArray[:, 1], minsArray[:, 0], label='mins')

# show the plot
plt.legend()
plt.show()