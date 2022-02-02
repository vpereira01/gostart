from matplotlib import pyplot as plt
import numpy as np
from sklearn import linear_model
from sklearn.multioutput import RegressorChain
from sklearn.linear_model import LogisticRegression
from sklearn.neural_network import MLPRegressor
from sklearn.model_selection import GridSearchCV
from sklearn.pipeline import Pipeline
from sklearn.preprocessing import StandardScaler
from sklearn.compose import TransformedTargetRegressor
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import PowerTransformer
from sklearn.preprocessing import MaxAbsScaler
from sklearn.pipeline import make_pipeline

# Load the data
f = open("../data/rawRecords_64_256000.csv")
f.readline()  # skip the header
data = np.loadtxt(f, delimiter=',')

rawX = data[:, 0].reshape(-1, 1)
rawY = data[:, 1:]

# xScaler = MaxAbsScaler()
# xScaler.fit(rawX)
# transformedRawX = xScaler.transform(rawX)
# yScaler = MaxAbsScaler()
# yScaler.fit(rawY)
# transformedRawY = yScaler.transform(rawY)
# X_train, X_test, y_train, y_test = train_test_split(transformedRawX, transformedRawY)
# regressor = MLPRegressor(max_iter=400)
# regressor.fit(X_train, y_train)
# print(regressor.score(X_test, y_test)) # = 0.8102673449412665

# interesting, convergence stopped so it could go on
X_train, X_test, y_train, y_test = train_test_split(rawX, rawY)
transformer = MaxAbsScaler()
regressor = MLPRegressor()
pipeline = make_pipeline(transformer, regressor)
pipeline.fit(X_train, y_train)
print('R2 score: {0:.2f}'.format(pipeline.score(X_test, y_test)))

y_pred = regressor.predict(X_test)
print('regressor.score(x_test, y_test)=' + str(regressor.score(X_test, y_test)))
print(f'y_pred[0]= {int(y_pred[0][0])},{int(y_pred[0][1])} y_test[0]= {int(y_test[0][0])},{int(y_test[0][0])}')

# this seems to suck but don't know why
# transformer = MaxAbsScaler()
# regressor = MLPRegressor()
# transformedTargetRegressor = TransformedTargetRegressor(regressor=regressor, transformer=transformer)
# X_train, X_test, y_train, y_test = train_test_split(rawX, rawY)
# transformedTargetRegressor.fit(X_train, y_train)
# print('R2 score: {0:.2f}'.format(transformedTargetRegressor.score(X_test, y_test)))
# raw_target_regr = regressor.fit(X_train, y_train)
# print('R2 score: {0:.2f}'.format(raw_target_regr.score(X_test, y_test)))
