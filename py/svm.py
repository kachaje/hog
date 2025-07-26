#!/usr/bin/env python3

import cv2
from sklearn.svm import SVC
from sklearn.model_selection import train_test_split
from sklearn.metrics import accuracy_score

# Assume you have 'images' (list of image arrays) and 'labels' (list of corresponding labels)

# 1. Split data into training and testing sets
X_train, X_test, y_train, y_test = train_test_split(images, labels, test_size=0.2, random_state=42)

# 2. HOG Feature Extraction
hog = cv2.HOGDescriptor() # Initialize HOG descriptor
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
