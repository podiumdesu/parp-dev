import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns

# List of file names for both data sets
real_files = ['sslip/output-real-1.csv', 'sslip/output-real-5.csv', 'sslip/output-real-10.csv', 
              'sslip/output-real-15.csv', 'sslip/output-real-20.csv', 'sslip/output-real-25.csv']
geth_files = ['geth/output-real-geth-1.csv', 'geth/output-real-geth-5.csv', 'geth/output-real-geth-10.csv', 
              'geth/output-real-geth-15.csv', 'geth/output-real-geth-20.csv', 'geth/output-real-geth-25.csv']
labels = [1, 5, 10, 15, 20, 25]  # Corresponding number of light clients

def process_files(files, label, client_type):
    all_data = []
    for file, label in zip(files, labels):
        df = pd.read_csv(file)
        df['CPU (%)'] = pd.to_numeric(df['CPU (%)'], errors='coerce')
        df['Memory (%)'] = pd.to_numeric(df['Memory (%)'], errors='coerce')
        df.dropna(inplace=True)  # Drop rows with NaN values
        df = df.iloc[10:70]  # Extract the relevant rows
        df['Number of Light Clients'] = label  # Add a column for the number of light clients
        df['Type'] = client_type  # Add a column for the type of data
        all_data.append(df)
    return pd.concat(all_data)

# Process both sets of files
real_data = process_files(real_files, labels, 'PARP')
geth_data = process_files(geth_files, labels, 'Geth')

# Combine all data for plotting
combined_data = pd.concat([real_data, geth_data])

# Reshape DataFrame for better plotting with seaborn
melted_data = pd.melt(combined_data, id_vars=['Number of Light Clients', 'Type'], value_vars=['CPU (%)', 'Memory (%)'],
                      var_name='Metric', value_name='Usage')

# Separate plots for CPU and Memory
for metric in ['CPU (%)', 'Memory (%)']:
    plt.figure(figsize=(10, 6))
    sns.boxplot(x='Number of Light Clients', y='Usage', hue='Type', 
                data=melted_data[melted_data['Metric'] == metric], palette="Set2")
    plt.title(f'{metric} Usage Comparison Across Different Light Client Instances')
    plt.ylabel('Server Memory Resource Usage (%)')
    plt.xlabel('Number of Light Clients')
    plt.legend(title='Server Type')
    plt.show()
