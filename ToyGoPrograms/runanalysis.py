import pandas as pd
import sys
import matplotlib.pyplot as plt

# Load data
if len(sys.argv) != 2:
    print("Usage: python runanalysis.py <filename>")
    sys.exit(1)

filename = sys.argv[1]
df = pd.read_csv(filename)

# Convert time to numeric
df["time"] = df["time"].astype(int)

# Create figure with two subplots
fig, ax1 = plt.subplots(figsize=(10, 5))

# Plot Heap in-use (Left Y-axis)
ax1.set_xlabel("Time (seconds)")
ax1.set_ylabel("Heap in Use (bytes)", color="green")
ax1.plot(df["time"], df["heap_inuse"], marker="o", linestyle="-", color="green", label="Heap in Use")
ax1.tick_params(axis="y", labelcolor="green")

# Create second Y-axis for GC CPU usage
ax2 = ax1.twinx()
ax2.set_ylabel("GC CPU Fraction", color="blue")
ax2.plot(df["time"], df["gc_cpu_fraction"], marker="s", linestyle="--", color="blue", label="GC CPU Usage")
ax2.tick_params(axis="y", labelcolor="blue")

# Title and grid
plt.title("Heap In-Use and GC CPU Usage Over Time")
fig.tight_layout()
plt.grid()

# Show the plot
plt.show()