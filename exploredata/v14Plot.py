import matplotlib.pyplot as plt
import numpy as np
import math
import pandas as pd
import sys
from collections import namedtuple
# from bigfloat import *

# base definitions
numberN = 11089180634223337999
startSubEst = 192582402-10000000

subEstimations = np.arange(startSubEst, startSubEst+5000, 2, dtype=int)

# solve([p1*p2 = numberN, p2-p1=subEst], [p1, p2]);
# p1(subEst) :=  (sqrt(subEst^2  + 4 * numberN) - subEst)/2;
# p2(subEst) :=  (sqrt(subEst^2  + 4 * numberN) + subEst)/2;

# we want to find perfect squares made by 4*numberN + subEst^2
# so we first start by doing the startNum = math.ceil(math.sqrt(4*numberN))
# then we check value1 = startNum**2 - 4*numberN so we get a candidate for value1 = subEst^2
# if sqrt(value1) % 1 != 0 ignore value, increment the startNum
# initial startNum 6660084275
# final startNum 6662868040, required at most 2783765 iterations
# because square of odd numbers are odd and startNum needs to be even we need only to test even startNums

def pFunction(subEst):
    return math.sqrt(subEst**2 + 4*numberN)
myPFunction = np.vectorize(pFunction, otypes=['float'])
p = myPFunction(subEstimations)

def pDiffFunction(subEst):
    pValue = pFunction(subEst)
    diffValue = pValue - round(pValue)
    return abs(diffValue)
myPDiffFunction = np.vectorize(pDiffFunction, otypes=['float'])
dp = myPDiffFunction(subEstimations)




def p1Function(subEst):
    return  (math.sqrt(subEst**2 + 4 * numberN ) - subEst)/2
myP1Function = np.vectorize(p1Function, otypes=['float'])
p1 = myP1Function(subEstimations)

def p1DiffFunction(subEst):
    p1Value = p1Function(subEst)
    diffValue = p1Value - round(p1Value)
    return abs(diffValue)
myP1DiffFunction = np.vectorize(p1DiffFunction, otypes=['float'])
dp1 = myP1DiffFunction(subEstimations)



def p2Function(subEst):
    return  (math.sqrt(subEst**2 + 4 * numberN ) + subEst)/2
myP2Function = np.vectorize(p2Function, otypes=['float'])
p2 = myP2Function(subEstimations)



def diffFunction(subEst):
    p1AsInt = round(p1Function(subEst))
    p2AsInt = round(p2Function(subEst))
    diffValue = numberN - (p1AsInt * p2AsInt)
    return int(abs(diffValue))
myDiffFunction = np.vectorize(diffFunction, otypes=['int'])
df = myDiffFunction(subEstimations)

# find all local mins
MinValuePosition = namedtuple('MinValuePosition', ['value', 'x', 'index', 'xdiff'])
mins = []
for i in range(1, len(subEstimations)-1):
    # sumEstimations can be a numpy int32/int64 which overflows, converting to python3 int fixes the issue
    prevSumEstimation = int(subEstimations[i-1])
    curSumEstimation = int(subEstimations[i])
    nextSumEstimation = int(subEstimations[i+1])

    prevValue = diffFunction(prevSumEstimation)
    curValue = diffFunction(curSumEstimation)
    nextValue = diffFunction(nextSumEstimation)

    if prevValue > curValue and curValue < nextValue:
        if len(mins) == 0:
            mins.append(MinValuePosition(curValue, curSumEstimation, i, 0))
        else:
            mins.append(MinValuePosition(curValue, curSumEstimation, i, mins[-1].x - curSumEstimation))

minsArray = np.array(mins)
# print(sumEstimations)
# print(p1)

# setting the axes at the centre
fig = plt.figure()
ax = fig.add_subplot(1, 1, 1)
ax.spines['left'].set_position('center')
ax.spines['bottom'].set_position('zero')
ax.spines['right'].set_color('none')
ax.spines['top'].set_color('none')
ax.xaxis.set_ticks_position('bottom')
ax.yaxis.set_ticks_position('left')

# plot the function
plt.plot(subEstimations,  dp, '.', label='p1')
# plt.plot(subEstimations,  p1, '.', label='p1')
# plt.plot(subEstimations, dp1, 'r', label='dp1')
# plt.plot(subEstimations, dp1, '.', label='dp1')
# plt.plot(subEstimations,  df,  '.'    ,label='df')
# plt.plot(subEstimations,  df,  'r'    ,label='df')
# plt.plot(c3,df3, label='df3')
# plt.plot(c4,df4, label='df4')
# plt.plot(minsArray[:, 1], minsArray[:, 0], label='mins')

# show the plot
plt.legend()
plt.show()