import pandas as pd
import json
import matplotlib.pyplot as plt

# Read the JSON file
with open('test_results.json', 'r') as file:
    lines = file.readlines()

# Parse the JSON data
data = []
for line in lines:
    data.append(json.loads(line))

# Convert to DataFrame
df = pd.json_normalize(data)
# Extract relevant columns
df['time'] = pd.to_datetime(df['data.time'])

# Plot each metric
metrics = {
    'heap_alloc_bytes': 'Heap Allocation (Bytes)',
    'heap_objects': 'Heap Objects (#)',
    'gc_runs': 'GC Runs (#)',
    'last_gc_pause_ns': 'Last GC Pause (ns)'
}

# Filter out rows where the time index is NaT
df = df.dropna(subset=['time'])

# Filter rows where metric is one of the specified metrics
df = df[df['metric'].isin(metrics.keys())]

# Filter out rows where data.value is NaN
df = df.dropna(subset=['data.value'])

# Grab necessary metrics, and set time as index
df = df[['time', 'metric', 'data.value']]
df.set_index('time', inplace=True)

for metric, title in metrics.items():
    metric_df = df[df['metric'] == metric]
    plt.figure()
    if metric == 'heap_alloc_bytes':
        metric_df.loc[:, 'data.value'] = metric_df['data.value'] / (1024 * 1024)  # Convert bytes to megabytes
        plt.ylabel('Heap Allocation (MB)')
    else:
        plt.ylabel(title)
    plt.fill_between(metric_df.index, metric_df['data.value'], alpha=0.3)
    plt.plot(metric_df.index, metric_df['data.value'])
    plt.title(title)
    plt.xlabel('Time')
    plt.grid(True)
    plt.savefig(f'rungraphs/{metric}.png')
    plt.show()
