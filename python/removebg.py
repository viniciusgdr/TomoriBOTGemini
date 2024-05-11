from rembg import remove
from PIL import Image
import sys
input = Image.open(sys.argv[1])
output = remove(input)
output.save(sys.argv[2])