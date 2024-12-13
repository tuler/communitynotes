import numpy as np

# Check if NumPy can see the M1's GPU
np.__config__.show()

import torch

# Check if MPS (Metal Performance Shaders) is available
print("Is MPS available?", torch.backends.mps.is_available())

# Check if your current PyTorch installation was built with MPS support
print("Is MPS built?", torch.backends.mps.is_built())

# Get current device
device = torch.device("mps" if torch.backends.mps.is_available() else "cpu")
print("Current device:", device)

# You can also try a simple operation to verify
x = torch.rand(5, 3)
if torch.backends.mps.is_available():
    x = x.to("mps")
    print("Successfully moved tensor to MPS device")
