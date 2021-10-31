import joblib
import pandas as pd

loaded_rf = joblib.load("./predictor.joblib")

my_data = {
  'age': [64],
  'sex': [1],
  'cp': [3],
  'trestbps': [170],
  'chol': [227],
  'fbs': [0],
  'restecg': [0],
  'thalach':[155],
  'exang': [0],
  'oldpeak': [0.6],
  'slope': [1],
  'ca': [0],
  'thal': [3]
}

myvar = pd.DataFrame(my_data)

prediction = loaded_rf.predict(myvar)
print("The prediction is ",prediction)

print('You are not at risk') if prediction[0] == 0 else print('You are at Risk')

