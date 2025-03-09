import matplotlib.pyplot as plt
import numpy as np

# Data for AST Parsing Program
ast_tuners = ["no tuner", "aimd", "rolling avg", "linear tuner", "flip flop", "threshold 1GB"]
ast_avg_heap = [85884133, 99374745, 164198771, 169568590, 163609801, 382257446]
ast_p99_heap = [106474368, 145753072, 294568000, 310584640, 292403408, 590616704]
ast_gc_cpu = [0.215638, 0.177323, 0.097974, 0.094690, 0.095302, 0.039836]

# Data for Nested Pointer Map Program
nested_tuners = ["No Tuner", "AIMD", "Rolling Avg", "Linear Tuner", "Flip Flop", "Threshold 1GB"]
nested_avg_heap = [4142343925, 4253992109, 6827388932, 6811350721, 6887911609, 6916982121]
nested_p99_heap = [4694352792, 5346867176, 11903151600, 11814333864, 11948653160, 12265711216]
nested_gc_cpu = [0.238416, 0.206628, 0.088090, 0.089675, 0.079936, 0.072610]

# Function to create bar charts
def create_bar_chart(tuners, ast_data, nested_data, title, ylabel, ast_label, nested_label):
    x = np.arange(len(tuners))  # the label locations
    width = 0.35  # the width of the bars

    fig, ax = plt.subplots()
    rects1 = ax.bar(x - width/2, ast_data, width, label=ast_label)
    rects2 = ax.bar(x + width/2, nested_data, width, label=nested_label)

    ax.set_xlabel('Tuner')
    ax.set_ylabel(ylabel)
    ax.set_title(title)
    ax.set_xticks(x)
    ax.set_xticklabels(tuners, rotation=45)
    ax.legend()

    plt.tight_layout()
    plt.show()

# Create charts for each metric
# create_bar_chart(ast_tuners, ast_avg_heap, nested_avg_heap, "Average HeapAlloc Comparison", "Bytes", "AST Parsing", "Nested Maps")
# create_bar_chart(ast_tuners, ast_p99_heap, nested_p99_heap, "P99 HeapAlloc Comparison", "Bytes", "AST Parsing", "Nested Maps")
create_bar_chart(ast_tuners, ast_gc_cpu, nested_gc_cpu, "GC CPU Fraction Comparison", "GC CPU Fraction", "AST Parsing", "Nested LL Maps")