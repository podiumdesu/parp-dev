import psutil
import subprocess
import time
import datetime

def get_pid_by_unique_identifier(unique_identifier):
    try:
        command = ["ps", "aux"]
        ps = subprocess.Popen(command, stdout=subprocess.PIPE)
        grep = subprocess.Popen(["grep", unique_identifier], stdin=ps.stdout, stdout=subprocess.PIPE)
        ps.stdout.close()
        output = grep.communicate()[0]
        
        for line in output.decode().split('\n'):
            if line:
                return line.split()[1]  # PID is typically the second column
    except subprocess.CalledProcessError as e:
        print("Failed to fetch process PID", e)
    return None

def get_cpu_usage(pid):
    try:
        process = psutil.Process(pid)
        return process.cpu_percent(interval=1)
    except psutil.NoSuchProcess:
        print(f"No process with PID {pid} found.")
        return None

def get_memory_usage(pid):
    try:
        process = psutil.Process(pid)
        memory_info = process.memory_info()
        return memory_info.rss, process.memory_percent()
    except psutil.NoSuchProcess:
        print(f"No process with PID {pid} found.")
        return None, None

def monitor_usage(pid):
    start_time = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    file_path = f"monitoring_results_{start_time.replace(':', '-')}.txt"
    with open(file_path, 'w') as log_file:
        try:
            while True:
                cpu_usage = get_cpu_usage(pid)
                rss_memory, memory_percentage = get_memory_usage(pid)
                if cpu_usage is not None and rss_memory is not None:
                    current_time = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
                    print(f"CPU Usage for PID {pid}: {cpu_usage}%, memory: {memory_percentage:.2f}%")

                    log_entry = f"{current_time} - CPU: {cpu_usage}%, Memory: {memory_percentage:.2f}%\n"
                    print(log_entry.strip())
                    log_file.write(log_entry)

                    # print(f"Memory Usage for PID {pid}: {rss_memory} bytes (RSS), ")
                time.sleep(1)  # Sleep for 1 second before checking again
        except KeyboardInterrupt:
            print("Monitoring stopped.")

if __name__ == "__main__":
    pid = get_pid_by_unique_identifier("ws.port=8100")
    if pid:
        print("PID:", pid)
        print(f"Monitoring CPU and memory usage for PID {pid}. Press Ctrl+C to stop.")
        monitor_usage(int(pid))
    else:
        print("Process not found.")
