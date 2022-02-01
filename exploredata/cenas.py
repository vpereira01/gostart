# %% [markdown]
# Generic imports

# %%
import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
import seaborn as sns

# Make NumPy printouts easier to read.
np.set_printoptions(precision=3, suppress=True)

import tensorflow as tf

from tensorflow import keras
from tensorflow.keras import layers

print(tf.__version__)

# %% [markdown]
# Data load

# %%
recordsUrl = '../data/records_150_64000.csv'

dataset = pd.read_csv(recordsUrl)

# dataset.pop("BiggerPrime")

# %% [markdown]
# Split dataset

# %%
train_dataset = dataset.sample(frac=0.8, random_state=0)
test_dataset = dataset.drop(train_dataset.index)

# %% [markdown]
# Visualize columns correlations

# %%
from mpl_toolkits.mplot3d import Axes3D
from matplotlib.colors import ListedColormap

fig = plt.figure(figsize=(6,6))
# ax = fig.add_subplot(projection='3d')
# ax.scatter(train_dataset['NumberN'], train_dataset['SmallerPrime'], train_dataset['BiggerPrime'], c=train_dataset['BiggerPrime'], cmap='gray')
ax = Axes3D(fig, auto_add_to_figure=True)
fig.add_axes(ax)
cmap = ListedColormap(sns.color_palette("husl", 256).as_hex())

sc = ax.scatter(train_dataset['NumberN'], train_dataset['SmallerPrime'], train_dataset['BiggerPrime'], c=train_dataset['SmallerPrime'], alpha=1)
ax.set_title("3D plot")
ax.set_xlabel('NumberN')
ax.set_ylabel('SmallerPrime')
ax.set_zlabel('BiggerPrime')

plt.show()