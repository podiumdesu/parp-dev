import csv
import re
def extract_and_write_data(start_time, lines, output_file, num_records=60, context_lines=10):
    """
    Extracts a number of records from the monitoring data starting from a given time and writes them to a CSV file.
    Also includes additional lines before and after the start time for context.
    
    :param start_time: Time to start recording data (format: "HH:MM:SS")
    :param input_file: Path to the file containing the monitoring data.
    :param output_file: Path to the output CSV file.
    :param num_records: Number of records to write after the start time.
    :param context_lines: Number of lines to include before and after the data.
    """
    
    # Find the start index
    start_index = None
    for i, line in enumerate(lines):
        if start_time in line:
            start_index = i
            break
    
    if start_index is None:
        print(f"Start time {start_time} not found in the file.")
        return
    
    # Adjust start and end index to include context lines
    start_index_with_context = max(start_index - context_lines, 0)
    end_index_with_context = start_index + num_records + context_lines

    # Extract the desired records with context lines
    records = lines[start_index_with_context:end_index_with_context]
    
    # Write to CSV
    with open(output_file, 'w', newline='') as csvfile:
        writer = csv.writer(csvfile)
        writer.writerow(['CPU (%)', 'Memory (%)'])  # Header
        for record in records:
            print(record)
            pattern = re.compile(r"\d+\.\d+")

            # Find all occurrences of the pattern
            matches = pattern.findall(record)

            cpu_usage = matches[0]
            memory_usage = matches[1]
            writer.writerow([cpu_usage, memory_usage])

    print(f"Data successfully written to {output_file}")

# Example usage
if __name__ == '__main__':
    num = 25

    input_log_file = f'real-geth-{num}.txt'
    with open(input_log_file, 'r') as f:
        lines = f.readlines()
    start_time = str(lines[0].rstrip('\n'))
    lines.pop(0)

    # start_time = '18:05:36'
    output_csv_file = f'output/output-real-geth-{num}.csv'
    # start_time = '18:05:34'  # Time you specify
    extract_and_write_data(start_time, lines, output_csv_file)
