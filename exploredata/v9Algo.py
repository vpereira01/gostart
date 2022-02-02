# Evaluate trying to determine all local difference minimums

import matplotlib.pyplot as plt
import numpy as np
import math
import pandas as pd
import sys
from collections import namedtuple
from bigfloat import *

# base definitions
numberN = 1137501950159415429602603162643696967245621181
# sumEst = 67568508516220459435322 # real estimation
sumEst = 67551153748445296705730 + 1234567
# realSum = 67551153748445296705730

# the sum of primes has to be even
startNumOddFix = (sumEst % 2)

def p1Function(c1):
     return -(sqrt(sumEst**2 - 2 * c1 * sumEst - 4 * numberN + c1**2 ) - sumEst + c1)/2

def p2Function(c1):
     return (sqrt(sumEst**2 - 2 * c1 * sumEst - 4 * numberN + c1**2 ) + sumEst - c1)/2

def diffFunction(c1):
    p1AsInt = int(round(p1Function(c1)))
    p2AsInt = int(round(p2Function(c1)))
    diffValue = numberN - (p1AsInt * p2AsInt)
    return int(abs(diffValue))

# proposed mins
MinValuePosition = namedtuple('MinValuePosition', ['value', 'x', 'index', 'xdiff'])
algoMins = []
for c2Value in np.arange(startNumOddFix, 512, 2):
    # print(f'diffFunction(c2Value)={str(diffFunction(c2Value))}')
    if diffFunction(c2Value-2) > diffFunction(c2Value) and diffFunction(c2Value) < diffFunction(c2Value+2):
        if len(algoMins) == 0:
            algoMins.append(MinValuePosition(diffFunction(c2Value), c2Value, -42, 0))
        else:
            algoMins.append(MinValuePosition(diffFunction(c2Value), c2Value, -42, c2Value - algoMins[-1].x))

        if len(algoMins) == 5:
            break

algoMinsArray = np.array(algoMins)
algoMinsArray = algoMinsArray[1:] # remove the first because of missing xdiff
# print(algoMinsArray)
algoMinsMinXDiff = algoMinsArray[:, 3].min()
algoMinsMaxXDiff = algoMinsArray[:, 3].max()
if (algoMinsMaxXDiff == algoMinsMinXDiff):
    sys.exit("algoMinsMaxXDiff == algoMinsMinXDiff")
if (algoMinsMinXDiff + 2 != algoMinsMaxXDiff):
    sys.exit("algoMinsMinXDiff + 2 != algoMinsMaxXDiff")

# minDiffC2 = algoMinsArray[1][0]
minDiffC2 = 9841402020351421

c2Value = algoMinsArray[1][1]
lastC2Value = c2Value-algoMinsMinXDiff
step = c2Value-lastC2Value
print(f'step={step}')
while True:
    diffC2MinusTwo = diffFunction(c2Value-2)
    if (diffC2MinusTwo == 0):
        print(f'diffC2MinusTwo={str(diffC2MinusTwo)} c2Value-2={str(c2Value-2)}')
        sys.exit("stopping 04")

    diffC2 = diffFunction(c2Value)
    if (diffC2 == 0):
        print(f'diffC2={str(diffC2)} c2Value={str(c2Value)}')
        sys.exit("stopping 01")
    
    diffC2PlusTwo = diffFunction(c2Value+2)
    if (diffC2PlusTwo == 0):
        print(f'diffC2PlusTwo={str(diffC2PlusTwo)} c2Value+2={str(c2Value+2)}')
        sys.exit("stopping 02")

    if diffC2 < diffC2PlusTwo and diffC2 < diffC2MinusTwo:
        step = c2Value-lastC2Value
        if (diffC2 < minDiffC2):
            print(f'c2Value={str(c2Value):>25} step={str(step):>2} diffC2={str(diffC2):>25} minDiffC2={str(minDiffC2):>25}')
            # minDiffC2 = diffC2
    elif diffC2PlusTwo < diffC2 and diffC2PlusTwo < diffC2MinusTwo:
        step = (c2Value+2)-lastC2Value
        if (diffC2PlusTwo < minDiffC2):
            print(f'c4Value={str(c2Value+2):>25} step={str(step):>2} diffC2={str(diffC2PlusTwo):>25} minDiffC2={str(minDiffC2):>25}')
            # minDiffC2 = diffC2PlusTwo
    elif diffC2MinusTwo < diffC2 and diffC2MinusTwo < diffC2PlusTwo:
        step = (c2Value-2)-lastC2Value
        if (diffC2MinusTwo < minDiffC2):
            print(f'c0Value={str(c2Value-2):>25} step={str(step):>2} diffC2={str(diffC2MinusTwo):>25} minDiffC2={str(minDiffC2):>25}')
            # minDiffC2 = diffC2MinusTwo
        if (step < 6):
            step = 6
            # print(f'adjusted step to 6 at {str(c2Value-2)}')
    else:
        sys.exit("no minimum found")
    # if (step <= 8):
    #     sys.exit(f'bug 01 {str(step)} {str(c2Value)} {str(lastC2Value)} {str(diffC2MinusTwo)} {str(diffC2)} {str(diffC2PlusTwo)}')
    lastC2Value = c2Value
    c2Value += step


    if (c2Value >= 17354767775162729000):
        print(f'c2Value={str(c2Value)}')

    if (c2Value >= 17354767775162729999):
        sys.exit("stopping 03")
    
    if c2Value % 5000 <= 15:
        print(f'c2Value={str(c2Value):>25} step={str(step):>2} diffC2={str(diffC2):>25}')
