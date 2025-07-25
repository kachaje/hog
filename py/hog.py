#!/usr/bin/env python3

from skimage.io import imread
from skimage.transform import resize
from skimage.feature import hog
from skimage import exposure
import matplotlib.pyplot as plt
import sys
import signal
import json

def signal_handler(sig, frame):
    sys.exit(0)


if __name__ == "__main__":
    signal.signal(signal.SIGINT, signal_handler)

    filename = "flower.jpg"

    if len(sys.argv) > 1:
        filename = str(sys.argv[1])

    img = imread(filename)
    plt.axis("off")
    plt.imshow(img)
    print(img.shape)

    plt.show()

    resized_img = resize(img, (128*4, 64*4))
    plt.axis("off")
    plt.imshow(resized_img)
    print(resized_img.shape)

    plt.show()

    fd, hog_image = hog(resized_img, orientations=9, pixels_per_cell=(8, 8), cells_per_block=(2, 2), visualize=True, channel_axis=-1)

    print(len(fd))

    with open(f"./outputFd.json", 'w') as f:
        json.dump(fd.tolist(), f, indent=4)

    plt.axis("off")
    plt.imshow(hog_image, cmap="gray")
    plt.show()