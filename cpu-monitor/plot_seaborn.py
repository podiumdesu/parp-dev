import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns

# List of file names and corresponding labels
files = ['output/output-real-1.csv', 'output/output-real-5.csv', 'output/output-real-10.csv', 
         'output/output-real-15.csv', 'output/output-real-20.csv', 'output/output-real-25.csv']
labels = [1, 5, 10, 15, 20, 25]  # Corresponding number of light clients

# Prepare to hold the data
all_data = []

# Load and process each file
for file, label in zip(files, labels):
    df = pd.read_csv(file)
    
    # Convert 'CPU (%)' and 'Memory (%)' columns to numeric, coercing errors to NaN
    df['CPU (%)'] = pd.to_numeric(df['CPU (%)'], errors='coerce')
    df['Memory (%)'] = pd.to_numeric(df['Memory (%)'], errors='coerce')
    
    df.dropna(inplace=True)  # Drop rows with NaN values
    df = df.iloc[10:70]  # Extract the relevant rows
    df['Number of Light Clients'] = label  # Add a column for the number of light clients
    all_data.append(df)

# Concatenate all data into a single DataFrame
combined_df = pd.concat(all_data)

# Reshape DataFrame for better plotting with seaborn
melted_df = pd.melt(combined_df, id_vars=['Number of Light Clients'], value_vars=['CPU (%)', 'Memory (%)'],
                    var_name='Metric', value_name='Usage')

# Create a figure with a single subplot
plt.figure(figsize=(10, 6))

# Plot using seaborn
sns.boxplot(x='Number of Light Clients', y='Usage', hue='Metric', data=melted_df, palette="Set2")

plt.title('CPU and Memory Usage Across Different Light Client Instances')
plt.ylabel('Usage (%)')
plt.xlabel('Number of Light Clients')

plt.show()
