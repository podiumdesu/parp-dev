import pandas as pd
import matplotlib.pyplot as plt

# List of file names
files = ['output/output-real-1.csv', 'output/output-real-5.csv', 'output/output-real-10.csv', 'output/output-real-15.csv', 'output/output-real-20.csv', 'output/output-real-25.csv']

# Prepare to hold the data
data_cpu = []
data_memory = []

# Load and process each file
for file in files:
    
    # print(file.split('-')[2].split('.')[0])
    # Read the dataset

    df = pd.read_csv(file, header=None, names=['CPU', 'Memory'])

    # print(df)    
    df['CPU'] = pd.to_numeric(df['CPU'], errors='coerce')
    df['Memory'] = pd.to_numeric(df['Memory'], errors='coerce')

    df.dropna(inplace=True)

    # Extract rows 11 to 70 (adjust indices for zero-based index)
    df = df.iloc[10:70]  # pandas is zero-indexed, so row 11 is index 10
    
    # Append to list
    data_cpu.append(df['CPU'])
    data_memory.append(df['Memory'])

    # print(data_cpu)
    # print(data_memory)

# # Create boxplots
# fig, axs = plt.subplots(2, 1, figsize=(10, 12))  # 2 plots, one for CPU and one for Memory

# # Plot CPU Usage
# axs[0].boxplot(data_cpu, labels=[file.split('-')[2].split('.')[0] for file in files])
# axs[0].set_title('CPU Usage Across Different Light Client Instances')
# axs[0].set_xlabel('Number of Light Clients')
# axs[0].set_ylabel('CPU Usage (%)')

# # Plot Memory Usage
# axs[1].boxplot(data_memory, labels=[file.split('-')[2].split('.')[0] for file in files])
# axs[1].set_title('Memory Usage Across Different Light Client Instances')
# axs[1].set_xlabel('Number of Light Clients')
# axs[1].set_ylabel('Memory Usage (%)')

# # Show the plots
# plt.tight_layout()
# plt.show()


# Create a figure and a set of subplots
fig, ax1 = plt.subplots(figsize=(10, 6))

# Plot CPU Usage on ax1
bp1 = ax1.boxplot(data_cpu, positions=range(1, len(files) * 2, 2), widths=0.6, patch_artist=True, boxprops=dict(facecolor='blue'))
ax1.set_xlabel('Number of Light Clients')
ax1.set_ylabel('CPU Usage (%)', color='blue')
ax1.tick_params(axis='y', labelcolor='blue')
ax1.set_title('CPU and Memory Usage Across Different Light Client Instances')

# Create ax2 for Memory Usage with shared x-axis and independent y-axis
ax2 = ax1.twinx()
bp2 = ax2.boxplot(data_memory, positions=range(2, len(files) * 2 + 1, 2), widths=0.6, patch_artist=True, boxprops=dict(facecolor='green'))
ax2.set_ylabel('Memory Usage (%)', color='green')
ax2.tick_params(axis='y', labelcolor='green')

# Add legend
ax1.legend([bp1["boxes"][0], bp2["boxes"][0]], ['CPU', 'Memory'], loc='upper left')

# Configure the x-axis to show labels for both CPU and Memory
labels = [file.split('-')[2].split('.')[0] for file in files]
ax1.set_xticks(range(1, len(files) * 2, 2))
ax1.set_xticklabels(labels)

plt.show()