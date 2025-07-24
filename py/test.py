#!/usr/bin/env python3

import matplotlib.pyplot as plt
from skimage import io
from skimage import color
from skimage.transform import resize
import math
from skimage.feature import hog
import numpy as np
import sys
import signal
import json
     
number_of_bins = 9
step_size = 180 / number_of_bins

def calculate_j(angle):
  temp = (angle / step_size) - 0.5
  j = math.floor(temp)
  return j

def calculate_Cj(j):
  Cj = step_size * (j + 0.5)
  return round(Cj, 9)

def calculate_value_j(magnitude, angle, j):
  Cj = calculate_Cj(j+1)
  Vj = magnitude * ((Cj - angle) / step_size)
  return round(Vj, 9)

i = 0
j = 0

mag = []
theta = []

with open("./backups/magnitudes.json", 'r') as f:
  mag = json.load(f) 

with open("./backups/angles.json", 'r') as f:
  theta = json.load(f) 

mag = np.array(mag)
theta = np.array(theta)

magnitude_values = [[mag[i][x] for x in range(j, j+8)] for i in range(i,i+8)]
angle_values = [[theta[i][x] for x in range(j, j+8)] for i in range(i, i+8)]
for k in range(len(magnitude_values)):
  for l in range(len(magnitude_values[0])):
    bins = [0.0 for _ in range(number_of_bins)]
    value_j = calculate_j(angle_values[k][l])
    Vj = calculate_value_j(magnitude_values[k][l], angle_values[k][l], value_j)
    Vj_1 = magnitude_values[k][l] - Vj

    with open(f"./backups/points/{j}_{i}_{l}_{k}.json", 'w') as f:
      json.dump({"i":i, "j": j, "k": k, "l": l, "value_j": value_j, "Vj": Vj, "Vj_1": Vj_1, "magnitude": magnitude_values[k][l], "angle": angle_values[k][l]}, f, indent=4)

    bins[value_j]+=Vj
    bins[value_j+1]+=Vj_1
    bins = [round(x, 9) for x in bins]


with open(f"./outputBin0_0.json", 'w') as f:
      json.dump(bins, f, indent=4)

