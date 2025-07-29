#!/usr/bin/env python3

from sklearn.svm import SVC
from sklearn.model_selection import train_test_split
from sklearn.metrics import accuracy_score
import os
import numpy as np
import pandas as pd
import cv2 as cv
from pathlib import Path
import warnings
from skimage.feature import hog
import tqdm

warnings.filterwarnings("ignore")
pd.options.display.max_columns = None

root = '../utils/fixtures'
style_file = '/styles.csv'
image_folder = root + '/images/'
print(root+style_file)
styles = pd.read_csv(Path(root+style_file))

styles.nunique()

import matplotlib.pyplot as plt
from tqdm import tqdm
categories = ['gender','masterCategory','subCategory','season','usage']
for cat in categories:
    df_cat = styles.groupby(cat,as_index=False).size().sort_values(ascending=True, by=cat).head(10)
    df_cat.plot(kind='bar',title = cat)
    plt.show()

images = []
labels = []

def load_image(ids,path=image_folder):
    img = cv.imread(image_folder+ids+'.jpg',cv.IMREAD_GRAYSCALE) #load at gray scale
    # img = cv.cvtColor(img, cv.COLOR_BGR2GRAY) #convert to gray scale
    return img,ids

for ids in tqdm(list(styles.id)[:20000]):
    if os.path.exists(image_folder+str(ids)+'.jpg'):
      img,ids = load_image(str(ids))
      if img is not None:
          images.append([img,int(ids)])
      labels.append(ids)

# 1. Split data into training and testing sets
X_train, X_test, y_train, y_test = train_test_split(images, labels, test_size=0.2, random_state=42)

# 2. HOG Feature Extraction
hog = cv.HOGDescriptor() # Initialize HOG descriptor
hog_features_train = []
for img in X_train:
    hog_features_train.append(hog.compute(img))

hog_features_test = []
for img in X_test:
    hog_features_test.append(hog.compute(img))

# Convert lists to NumPy arrays
import numpy as np
hog_features_train = np.array(hog_features_train).squeeze()
hog_features_test = np.array(hog_features_test).squeeze()

# 3. SVM Training
svm_model = SVC(kernel='linear') # Or other kernels like 'rbf'
svm_model.fit(hog_features_train, y_train)

# 4. Prediction and Evaluation
y_pred = svm_model.predict(hog_features_test)
accuracy = accuracy_score(y_test, y_pred)
print(f"Accuracy: {accuracy * 100:.2f}%")
