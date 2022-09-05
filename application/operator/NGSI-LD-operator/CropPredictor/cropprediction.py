import joblib
from sklearn.tree import DecisionTreeRegressor
from sklearn.preprocessing import StandardScaler
import pandas as pd
import numpy as np
from sklearn.model_selection import train_test_split

# Reading the csv file
data = pd.read_csv('cpdata.csv')
print(data.head(1))

# Creating dummy variable for target i.e label
label = pd.get_dummies(data.label).iloc[:, 1:]
data = pd.concat([data, label], axis=1)
data.drop('label', axis=1, inplace=True)
print('The data present in one row of the dataset is')
print(data.head(1))
train = data.iloc[:, 0:4].values
test = data.iloc[:, 4:].values

# Dividing the data into training and test set
X_train, X_test, y_train, y_test = train_test_split(train, test, test_size=0.3)

sc = StandardScaler()
X_train = sc.fit_transform(X_train)
X_test = sc.transform(X_test)

# Importing Decision Tree classifier
clf = DecisionTreeRegressor()

# Fitting the classifier into training set
clf.fit(X_train, y_train)
# pred=clf.predict(X_test)

joblib.dump(clf, "./croppredictor.joblib")

print("*****************X_test is *************")
print(X_test)
Y_pred_rf = clf.predict(X_test)
print("*****************Y_pred_rf is *************")
print(Y_pred_rf)


'''from sklearn.metrics import accuracy_score
# Finding the accuracy of the model
a=accuracy_score(y_test,pred)
print("The accuracy of this model is: ", a*100)'''


# Using firebase to import data to be tested
'''from firebase import firebase
firebase =firebase.FirebaseApplication('https://cropit-eb156.firebaseio.com/')
tp=firebase.get('/Realtime',None)

ah=tp['Air Humidity']
atemp=tp['Air Temp']
shum=tp['Soil Humidity']
pH=tp['Soil pH']
rain=tp['Rainfall']


l=[]
l.append(ah)
l.append(atemp)
l.append(pH)
l.append(rain)
predictcrop=[l]

# Putting the names of crop in a single list
crops=['wheat','mungbean','Tea','millet','maize','lentil','jute','cofee','cotton','ground nut','peas','rubber','sugarcane','tobacco','kidney beans','moth beans','coconut','blackgram','adzuki beans','pigeon peas','chick peas','banana','grapes','apple','mango','muskmelon','orange','papaya','watermelon','pomegranate']
cr='rice'

#Predicting the crop
predictions = clf.predict(predictcrop)
count=0
for i in range(0,30):
    if(predictions[0][i]==1):
        c=crops[i]
        count=count+1
        break;
    i=i+1
if(count==0):
    print('The predicted crop is %s'%cr)
else:
    print('The predicted crop is %s'%c)
'''
