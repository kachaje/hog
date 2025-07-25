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



def signal_handler(sig, frame):
    sys.exit(0)


if __name__ == "__main__":
    signal.signal(signal.SIGINT, signal_handler)

    filename = "flower.jpg"

    if len(sys.argv) > 1:
        filename = str(sys.argv[1])

    fig, axes = plt.subplots(2, 2, figsize=(10, 8))

    img = resize(color.rgb2gray(io.imread(filename)), (128, 64))

    axes[0, 0].imshow(img, cmap="gray")
    axes[0, 0].axis("off")
    axes[0, 0].set_title(filename)

    img = np.array(img)

    with open(f"./backups/dump.json", 'w') as f:
      json.dump(img.tolist(), f, indent=4)

    mag = []
    theta = []
    for i in range(128):
      magnitudeArray = []
      angleArray = []
      for j in range(64):
        # Condition for axis 0
        if j-1 <= 0 or j+1 >= 64:
          if j-1 <= 0:
            # Condition if first element
            Gx = img[i][j+1] - 0
          elif j + 1 >= len(img[0]):
            Gx = 0 - img[i][j-1]
        # Condition for first element
        else:
          Gx = img[i][j+1] - img[i][j-1]
        
        # Condition for axis 1
        if i-1 <= 0 or i+1 >= 128:
          if i-1 <= 0:
            Gy = 0 - img[i+1][j]
          elif i +1 >= 128:
            Gy = img[i-1][j] - 0
        else:
          Gy = img[i-1][j] - img[i+1][j]

        # Calculating magnitude
        magnitude = math.sqrt(pow(Gx, 2) + pow(Gy, 2))
        magnitudeArray.append(round(magnitude, 9))

        # Calculating angle
        if Gx == 0:
          angle = math.degrees(0.0)
        else:
          angle = math.degrees(abs(math.atan(Gy / Gx)))
        angleArray.append(round(angle, 9))
      mag.append(magnitudeArray)
      theta.append(angleArray)

    with open(f"./backups/magnitudes.json", 'w') as f:
      json.dump(mag, f, indent=4)

    with open(f"./backups/angles.json", 'w') as f:
      json.dump(theta, f, indent=4)

    mag = np.array(mag)
    theta = np.array(theta)

    axes[0, 1].imshow(mag, cmap="gray")
    axes[0, 1].axis("off")
    axes[0, 1].set_title("magnitude")

    axes[1, 0].imshow(theta, cmap="gray")
    axes[1, 0].axis("off")
    axes[1, 0].set_title("theta")

    histogram_points_nine = []
    for i in range(0, 128, 8):
      temp = []
      for j in range(0, 64, 8):
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

        temp.append(bins)

        with open(f"./backups/bins/{j}_{i}.json", 'w') as f:
              json.dump(bins, f, indent=4)

      with open(f"./backups/bins/temp_{i}.json", 'w') as f:
            json.dump(temp, f)

      histogram_points_nine.append(temp)

    axes[1, 0].bar(x=np.arange(9), height=histogram_points_nine[0][0], align="center", width=0.8)

    print(len(histogram_points_nine))
    print(len(histogram_points_nine[0]))
    print(len(histogram_points_nine[0][0]))

    with open(f"./backups/hist.json", 'w') as f:
          json.dump(histogram_points_nine, f)

    epsilon = 1e-05

    feature_vectors = []
    for i in range(0, len(histogram_points_nine) - 1, 1):
      temp = []
      for j in range(0, len(histogram_points_nine[0]) - 1, 1):
        values = [[histogram_points_nine[i][x] for x in range(j, j+2)] for i in range(i, i+2)]

        with open(f"./backups/features/values_{i}_{j}.json", 'w') as f:
              json.dump(values, f)

        final_vector = []
        for k in values:
          for l in k:
            for m in l:
              final_vector.append(m)

        with open(f"./backups/features/vector_round_1_{i}_{j}.json", 'w') as f:
              json.dump(final_vector, f)

        k = round(math.sqrt(sum([pow(x, 2) for x in final_vector])), 9)

        with open(f"./backups/features/vector_k_{i}_{j}.json", 'w') as f:
              json.dump(k, f)

        final_vector = [round(x/(k + epsilon), 9) for x in final_vector]

        with open(f"./backups/features/vector_round_2_{i}_{j}.json", 'w') as f:
              json.dump(final_vector, f)

        temp.append(final_vector)
      feature_vectors.append(temp)

    print("----------------------")

    print(len(feature_vectors))
    print(len(feature_vectors[0]))
    print(len(feature_vectors[0][0]))

    print("----------------------")

    print(f'Number of HOG features = {len(feature_vectors) * len(feature_vectors[0]) * len(feature_vectors[0][0])}')

    with open(f"./backups/features.json", 'w') as f:
          json.dump(feature_vectors, f)


    # axes[1, 1].set_visible(False)

    plt.show()
