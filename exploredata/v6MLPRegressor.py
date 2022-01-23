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

abort # still useless

# Load the data
f = open("../data/rawRecords_64_256000.csv")
f.readline()  # skip the header
data = np.loadtxt(f, delimiter=',')

rawX = data[:, 0].reshape(-1, 1)
rawY = data[:, 1:]
# add sum column
sumColumn = np.sum(rawY,1).reshape(-1, 1)
# not raw anymore but OK
# rawY = np.append(rawY, sumColumn, 1)
newX = rawX
newY = sumColumn

# interesting
X_train, X_test, y_train, y_test = train_test_split(newX, newY)
transformer = MaxAbsScaler()
regressor = MLPRegressor()
pipeline = make_pipeline(transformer, regressor)
pipeline.fit(X_train, y_train)
print('R2 score: {0:.2f}'.format(pipeline.score(X_test, y_test)))

y_pred = regressor.predict(X_test)
print('regressor.score(x_test, y_test)=' + str(regressor.score(X_test, y_test)))
print(f'y_pred[0]= {int(y_pred[0])} y_test[0]= {int(y_test[0])}')