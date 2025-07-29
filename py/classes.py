#!/usr/bin/env python3

# %%
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

# %%
root = '../utils/fixtures'
style_file = '/styles.csv'
image_folder = root + '/images/'
print(root+style_file)
styles = pd.read_csv(Path(root+style_file))

# %%
print("Style shape: ", str(styles.shape))
styles.head()

# %%
styles.nunique()

# %%
import matplotlib.pyplot as plt
from tqdm import tqdm
categories = ['gender','masterCategory','subCategory','season','usage']
for cat in categories:
    df_cat = styles.groupby(cat,as_index=False).size().sort_values(ascending=True, by=cat).head(10)
    df_cat.plot(kind='bar',title = cat)
    plt.show()

# %%
all_images = []
#labels = []

def load_image(ids,path=image_folder):
    img = cv.imread(image_folder+ids+'.jpg',cv.IMREAD_GRAYSCALE) #load at gray scale
    #img = cv.cvtColor(img, cv.COLOR_BGR2GRAY) #convert to gray scale
    return img,ids

for ids in tqdm(list(styles.id)[:20000]):
    if os.path.exists(image_folder+str(ids)+'.jpg'):
      img,ids = load_image(str(ids))
      if img is not None:
          all_images.append([img,int(ids)])
      #labels.append(ids)
len(all_images)

# %%
# np.array(all_images[0])

# %%
def resize_image(img,ids):
    return cv.resize(img, (60, 80),interpolation =cv.INTER_LINEAR)
    
all_images_resized = [[resize_image(x,y),y] for x,y in all_images]
len(all_images_resized)

# %%
styles.head()

# %%
[styles.masterCategory.value_counts().index]

# %%
df_labels = pd.DataFrame(all_images_resized,columns=['image','id'])

target = 'masterCategory'
categories = ['Apparel', 'Accessories', 'Footwear', 'Personal Care', 'Free Items']
df_train = styles[styles[target].isin(categories)][['id',target]]

df_labels = pd.merge(df_labels,df_train,how='left',on=['id'])
df_labels = df_labels.fillna('Others')
df_labels['class'] = pd.factorize(df_labels[target])[0]
print("Data Shape: ", str(df_labels.shape))
print(df_labels[target].value_counts())

# %%
mapper = df_labels[['class',target]].drop_duplicates()

# %%
for image in df_labels.image[:20]:
    print(image.shape)
    plt.imshow(image)
    fast = cv.FastFeatureDetector_create(50)
    kp = fast.detect(image,None)
    img2 = cv.drawKeypoints(image, kp, None, color=(255,0,0))
    # Print all default params
    #print( "Threshold: {}".format(fast.getThreshold()) )
    #print( "neighborhood: {}".format(fast.getType()) )
    print( "Total Keypoints with nonmaxSuppression: {}".format(len(kp)))
    fast_image=cv.drawKeypoints(image,kp,image)
    plt.imshow(fast_image);plt.title('FAST Detector')
    plt.show()

# %%
train_images = np.stack(df_labels.image.values,axis=0)
n_samples = len(train_images)
data_images = train_images.reshape((n_samples, -1))

# %%
ppcr = 8
ppcc = 8
hog_images = []
hog_features = []
for image in tqdm(train_images):
    blur = cv.GaussianBlur(image,(5,5),0)
    fd,hog_image = hog(blur, orientations=8, pixels_per_cell=(ppcr,ppcc),cells_per_block=(2,2),block_norm= 'L2',visualize=True)
    hog_images.append(hog_image)
    hog_features.append(fd)

hog_features = np.array(hog_features)

hog_features.shape

# %%
edges = [cv.Canny(image,50,150,apertureSize = 3) for image in train_images]
edges = np.array(edges)
n_samples_edges = len(edges)
edge_images = edges.reshape((n_samples, -1))
edge_images.shape

# %%
train_images.shape, hog_features.shape, edge_images.shape

# %%
edge_hog = np.hstack([hog_features,edge_images])
edge_hog.shape

# %%
histr = [cv.calcHist([img],[0],None,[256],[0,256]) for img in train_images]
histr = np.array(histr)
n_samples_histr = len(histr)
image_hist = histr.reshape((n_samples_histr, -1))
image_hist.shape

# %%
edge_hog = np.hstack([hog_features,edge_images,image_hist])
edge_hog.shape

# %%
from sklearn.model_selection import train_test_split
X_train, X_test, y_train, y_test = train_test_split(hog_features,df_labels['class'],test_size=0.2,stratify=df_labels['class'])
print('Training data and target sizes: \n{}, {}'.format(X_train.shape,y_train.shape))
print('Test data and target sizes: \n{}, {}'.format(X_test.shape,y_test.shape))

# %%
y_train.value_counts(),y_test.value_counts()

# %%
from sklearn import datasets, svm, metrics
from sklearn.neighbors import KNeighborsClassifier
from sklearn import metrics
# # Create a classifier: a support vector classifier
# classifier = svm.SVC(gamma=0.001)
# #fit to the trainin data
# classifier.fit(X_train,y_train)

# %%
from sklearn.preprocessing import StandardScaler
from sklearn.ensemble import RandomForestClassifier

test_accuracy = []
scaler = StandardScaler()
X_scaled = scaler.fit_transform(X_train)

classifier = KNeighborsClassifier(n_neighbors=3,algorithm='brute')
classifier.fit(X_scaled, y_train)
test_accuracy = classifier.score(scaler.transform(X_test), y_test)
print(test_accuracy)

# #FOR TUNING
# print(search_params)
# for p in tqdm(search_params):
#     #classifier = svm.SVC(gamma=p)
#     classifier = RandomForestClassifier(max_depth=8,n_estimators=600)

#     classifier.fit(X_scaled, y_train)
#     test_accuracy.append([p,classifier.score(scaler.transform(X_test), y_test)])

# df_accuracy = pd.DataFrame(test_accuracy,columns =['gamma','accuracy'])
# df_accuracy.index = df_accuracy.gamma
# df_accuracy[['accuracy']].plot()
# plt.show()

# %%
##PCA
from sklearn.decomposition import PCA
# pca = PCA(.90)
# principalComponents = pca.fit_transform(X = X_scaled)

# %%
mapper= mapper.reset_index(drop=True)

# %%
y_pred = classifier.predict(scaler.transform(X_test))

df_result = pd.DataFrame(y_test)
df_result['id'] = df_result.index
df_result = df_result.rename(columns={'class':'actual'})
df_result['predicted'] = y_pred
df_result = df_result.reset_index(drop=True)
df_result = pd.merge(df_result,mapper,left_on='predicted',right_on = 'class',how='inner')
df_result = df_result.drop(columns=['class'],axis=1)
df_result = df_result.rename(columns={'gender':'predicted_category'})

df_result = pd.merge(df_result,mapper,left_on='actual',right_on = 'class',how='inner')
df_result = df_result.drop(columns=['class'],axis=1)
df_result.shape

# %%
#some references for debugging
kd = df_result[df_result.actual!=df_result.predicted]
print(kd.shape)
kd.head()

# %%
# image_id = styles[styles.index==2663]['id'].reset_index(drop=True)
# k = str(image_id)
# print(k)
#print(image_folder+str(image_id)+'.jpg')

#debug image with id
img = cv.imread(image_folder+str(7347)+'.jpg')
img.shape
plt.imshow(img)

# %%
list_of_categories = categories +['Others']

print("Classification Report: \n Target: %s \n Labels: %s \n Classifier: %s:\n%s\n"
      % (target,list_of_categories,classifier, metrics.classification_report(y_test, y_pred)))

df_report = pd.DataFrame(metrics.confusion_matrix(y_test, y_pred),columns = list_of_categories )
df_report.index = [list_of_categories]
df_report

# %%
categories

# %%


# %%
df_report

# %%

#test image with id
test_data_location = root+'/test/'

img = cv.imread(test_data_location+'2093.jpg',cv.IMREAD_GRAYSCALE) #load at gray scale
image = cv.resize(img, (60, 80),interpolation =cv.INTER_LINEAR)

ppcr = 8
ppcc = 8
hog_images_test = []
hog_features_test = []

blur = cv.GaussianBlur(image,(5,5),0)
fd_test,hog_img = hog(blur, orientations=8, pixels_per_cell=(ppcr,ppcc),cells_per_block=(2,2),block_norm= 'L2',visualize=True)
hog_images_test.append(hog_img)
hog_features_test.append(fd)

hog_features_test = np.array(hog_features_test)
y_pred_user = classifier.predict(scaler.transform(hog_features_test))
#print(plt.imshow(hog_images_test))
print(y_pred_user)
print("Predicted MaterCategory: ", mapper[mapper['class']==int(y_pred_user)]['masterCategory'])

# %%
plt.imshow(hog_img)
plt.show()

# %%
hog_features.shape

# %%
from sklearn.preprocessing import MinMaxScaler
from sklearn.neighbors import NearestNeighbors

scaler_global = MinMaxScaler()
final_features_scaled = scaler_global.fit_transform(hog_features)
    
neighbors = NearestNeighbors(n_neighbors=20, algorithm='brute')
neighbors.fit(final_features_scaled)

distance,potential = neighbors.kneighbors(scaler_global.transform(hog_features_test))
print("Potential Neighbors Found!")
neighbors = []
for i in potential[0]:
    neighbors.append(i)

recommendation_list = list(df_labels.iloc[neighbors]['id'])
recommendation_list


