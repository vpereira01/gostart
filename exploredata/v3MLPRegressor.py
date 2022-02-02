from matplotlib import pyplot as plt
import numpy as np
from sklearn import linear_model
from sklearn.multioutput import RegressorChain
from sklearn.linear_model import LogisticRegression
from sklearn.neural_network import MLPRegressor
from sklearn.model_selection import GridSearchCV

# Load the data
f = open("../data/records_150_64000.csv")
f.readline()  # skip the header
data = np.loadtxt(f, delimiter=',')

regressor = MLPRegressor(
    hidden_layer_sizes=(16384)
    #,n_iter_no_change=32
    # activation='relu', 
    # solver='lbfgs', 
    # alpha=0.0001,
    )
rawX = data[:, 0].reshape(-1, 1)
# print('X')
# print(X)
rawY = data[:, 1:]
# print('Y')
# print(Y)
x_test = rawX[-2000:] # last 200 records
y_test = rawY[-2000:] # last 200 records
x_train = rawX[:-2000] # all except last 200 records
y_train = rawY[:-2000] # all except last 200 records

print("len(x_train)=" + str(len(x_train)))
print("len(y_train)=" + str(len(y_train)))
regressor.fit(x_train, y_train)

y_pred = regressor.predict(x_test)

print('regressor.score(x_test, y_test)=' + str(regressor.score(x_test, y_test)))
print('y_pred[0]=' + str(y_pred[0]) + " y_test[0]=" + str(y_test[0]))

plt.scatter(x_test, y_test[:, 0], color="black")
plt.scatter(x_test, y_pred[:, 0], color="blue")
plt.show()

plt.scatter(x_test, y_test[:, 1], color="black")
plt.scatter(x_test, y_pred[:, 1], color="blue")
plt.show()

plt.plot(regressor.loss_curve_)
plt.title("Loss Curve", fontsize=14)
plt.xlabel('Iterations')
plt.ylabel('Cost')
plt.show()

# from sklearn.experimental import enable_halving_search_cv  # noqa
# from sklearn.model_selection import HalvingRandomSearchCV
# param_distributions = {
#     'hidden_layer_sizes': [(512), (512,256), (512,256,128), (512,256,128,64), (512,256,128,64,32), (512,256,128,64,32), (512,256,128,64,32,16), (512,256,128,64,32,16,8)],
#     'max_iter': [32, 64, 128, 256, 512, 1024],
#     'activation': ['logistic', 'tanh', 'relu'],
#     'solver': ['lbfgs', 'sgd', 'adam'],
#     'alpha': [0.0001, 0.001, 0.01, 0.05],
#     'learning_rate': ['constant', 'invscaling', 'adaptive'],
# }

# randomSearch = HalvingRandomSearchCV(regressor, param_distributions).fit(x_train, y_train)

# print(randomSearch.best_params_) 
# this took almost 24hours
# #{'solver': 'lbfgs', 'max_iter': 256, 'learning_rate': 'constant', 'hidden_layer_sizes': (512, 256, 128, 64, 32, 16), 'alpha': 0.0001, 'activation': 'relu'}