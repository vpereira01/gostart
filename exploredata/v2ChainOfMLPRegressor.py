from matplotlib import pyplot as plt
import numpy as np
from sklearn import linear_model
from sklearn.multioutput import RegressorChain
from sklearn.linear_model import LogisticRegression
from sklearn.neural_network import MLPRegressor

# Load the data
f = open("../data/records_64_256000.csv")
f.readline()  # skip the header
data = np.loadtxt(f, delimiter=',')

# base_estimator = linear_model.LinearRegression()
base_estimator = MLPRegressor()
rawX = data[:, 0].reshape(-1, 1)
rawY = data[:, 1:]
negativeTwentyPercentOfOtalRecords = int(len(rawX) * 0.20 * -1)
x_test = rawX[negativeTwentyPercentOfOtalRecords:] # last 20% records
y_test = rawY[negativeTwentyPercentOfOtalRecords:] # last 20% records
x_train = rawX[:negativeTwentyPercentOfOtalRecords] # all except last 20% records
y_train = rawY[:negativeTwentyPercentOfOtalRecords] # all except last 20% records


chain = RegressorChain(base_estimator=base_estimator)
chain = chain.fit(x_train, y_train)

y_pred = chain.predict(x_test)
print(x_test)
print(y_pred)

print(chain.score(x_test, y_test))

plt.scatter(x_test, y_test[:, 0], color="black")
plt.scatter(x_test, y_pred[:, 0], color="blue")
plt.show()

plt.scatter(x_test, y_test[:, 1], color="black")
plt.scatter(x_test, y_pred[:, 1], color="blue")
plt.show()

# logreg = LogisticRegression(solver='lbfgs',multi_class='multinomial')
# X, Y = [[1, 0], [0, 1], [1, 1]], [[0, 2], [1, 1], [2, 0]]
# print(X)
# print(Y)
# chain = RegressorChain(base_estimator=logreg, order=[0, 1]).fit(X, Y)
# print(chain.predict(X))