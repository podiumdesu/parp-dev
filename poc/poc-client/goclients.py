import subprocess
import time
import datetime
import os
from concurrent.futures import ThreadPoolExecutor, as_completed

def run_client_instance():
    """Function to start a client and handle its output."""
    process = subprocess.Popen(['./myclienttogeth'], stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    stdout, stderr = process.communicate()
    return stdout, stderr

def write_current_datetime_to_file(client_number):
    # Get the current date and time
    current_datetime = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")

    # Define the directory path where the logs will be stored
    directory = f"client-logs/geth/pyscript/"
    
    # Ensure the directory exists; if not, create it
    if not os.path.exists(directory):
        os.makedirs(directory)
    
    # Define the full path for the log file
    file_path = f"{directory}client-{client_number}-{current_datetime}.txt"
    
    # Open a file. If the file does not exist, it will be created. If it exists, it will be overwritten.
    try:
        with open(file_path, "w") as file:
            file.write(f"Current Date and Time: {current_datetime}\n")
            file.write(f"Current client number: {client_number}")
        print("Date and time written successfully")
    except IOError as e:
        print(f"Error opening or writing to file: {e}")

def main(num_instances):
    number_of_instances = num_instances  # Define how many clients you want to run
    results = []

    with ThreadPoolExecutor(max_workers=200) as executor:
        futures = [executor.submit(run_client_instance) for _ in range(number_of_instances)]
        for future in as_completed(futures):
            stdout, stderr = future.result()
            print(f"Output: {stdout}")
            if stderr:
                print(f"Error: {stderr}")

if __name__ == '__main__':
    num_instances = 25  # Number of times you want to run the program
    write_current_datetime_to_file(num_instances)

    main(num_instances)