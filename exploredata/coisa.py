from matplotlib import pyplot as plt
import numpy as np
from sklearn import linear_model
from sklearn.multioutput import RegressorChain
from sklearn.linear_model import LogisticRegression
from sklearn.neural_network import MLPRegressor

# Load the data
f = open("../data/records_150_64000.csv")
f.readline()  # skip the header
data = np.loadtxt(f, delimiter=',')

# base_estimator = linear_model.LinearRegression()
base_estimator = MLPRegressor(random_state=1, max_iter=500, learning_rate="adaptive", hidden_layer_sizes=(9,5,3))
X = data[:, 0].reshape(-1, 1)
# print('X')
# print(X)
Y = data[:, 1:]
# print('Y')
# print(Y)
chain = RegressorChain(base_estimator=base_estimator)
chain = chain.fit(X, Y)

x_test = X[:-200]
y_test = Y[:-200]

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