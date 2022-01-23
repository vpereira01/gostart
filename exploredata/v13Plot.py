import matplotlib.pyplot as plt
import numpy as np
import math
import pandas as pd
import sys
from collections import namedtuple
# from bigfloat import *

# base definitions
numberN = 11089180634223337999
startSumEst = 6662868040-5000

sumEstimations = np.arange(startSumEst, startSumEst+10000, 2, dtype=int)

# solve([p1*p2 = numberN, p1+p2=sumEst], [p1, p2]);
# p1(sumEst) := -(sqrt(sumEst^2  - 4 * numberN ) - sumEst)/2;
# p2(sumEst) :=  (sqrt(sumEst^2  - 4 * numberN ) + sumEst)/2;

def p1Function(sumEst):
    return  -(math.sqrt(sumEst**2 - 4 * numberN ) - sumEst)/2
myP1Function = np.vectorize(p1Function, otypes=['float'])
p1 = myP1Function(sumEstimations)

def p2Function(sumEst):
    return  (math.sqrt(sumEst**2 - 4 * numberN ) + sumEst)/2
myP2Function = np.vectorize(p2Function, otypes=['float'])
p2 = myP2Function(sumEstimations)

def p1DiffFunction(sumEst):
    p1Value = p1Function(sumEst)
    diffValue = p1Value - round(p1Value)
    return diffValue
myP1DiffFunction = np.vectorize(p1DiffFunction, otypes=['float'])
dp1 = myP1DiffFunction(sumEstimations)

def diffFunction(c1):
    p1AsInt = round(p1Function(c1))
    p2AsInt = round(p2Function(c1))
    diffValue = numberN - (p1AsInt * p2AsInt)
    return int(abs(diffValue))
myDiffFunction = np.vectorize(diffFunction, otypes=['int'])
df = myDiffFunction(sumEstimations)

# find all local mins
MinValuePosition = namedtuple('MinValuePosition', ['value', 'x', 'index', 'xdiff'])
mins = []
for i in range(1, len(sumEstimations)-1):
    # sumEstimations can be a numpy int32/int64 which overflows, converting to python3 int fixes the issue
    prevSumEstimation = int(sumEstimations[i-1])
    curSumEstimation = int(sumEstimations[i])
    nextSumEstimation = int(sumEstimations[i+1])

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
# plt.plot(sumEstimations,  p1, '.', label='p1')
# plt.plot(sumEstimations, p1, 'r', label='dp1')
# plt.plot(sumEstimations, dp1, '.', label='dp1')
plt.plot(sumEstimations,  df,  '.'    ,label='df')
plt.plot(sumEstimations,  df,  'r'    ,label='df')
# plt.plot(c3,df3, label='df3')
# plt.plot(c4,df4, label='df4')
plt.plot(minsArray[:, 1], minsArray[:, 0], label='mins')

# show the plot
plt.legend()
plt.show()