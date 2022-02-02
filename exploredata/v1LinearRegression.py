import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
from sklearn import datasets, linear_model
from sklearn.metrics import mean_squared_error, r2_score

f = open("../data/records_64_256000.csv")
f.readline()  # skip the header
data = np.loadtxt(f, delimiter=',')
print(data)
X = data[:, 0].reshape(-1, 1)
y = data[:, 1].reshape(-1, 1)

# Split the data into training/testing sets
X_train = X[:-2000]
X_test = X[-2000:]

# Split the targets into training/testing sets
y_train = y[:-2000]
y_test = y[-2000:]

# Create linear regression object
regr = linear_model.LinearRegression()

# Train the model using the training sets
regr.fit(X_train, y_train)

# Make predictions using the testing set
y_pred = regr.predict(X_test)

# The coefficients
print("Coefficients: \n", regr.coef_)
# The mean squared error
print("Mean squared error: %.2f" % mean_squared_error(y_test, y_pred))
# The coefficient of determination: 1 is perfect prediction
print("Coefficient of determination: %.2f" % r2_score(y_test, y_pred))

# Plot outputs
plt.scatter(X_test, y_test, color="black")
plt.plot(X_test, y_pred, color="blue", linewidth=3)

plt.xticks(())
plt.yticks(())

plt.show()