import subprocess

def run_command(hash_workers_value=4, data_workers_value=1, comp_workers_value=0, input_value="simple", output_file=None, runs=10):
    total_hashGroupTime = 0
    total_compareTreeTime = 0
    
    for _ in range(runs):
        input_v = f"input/{input_value}"
        cmd = ["./BST", "-hash-workers", str(hash_workers_value), "-data-workers", str(data_workers_value), "-comp-workers", str(comp_workers_value), "-input", input_v, "-print", str(0)]
        
        process = subprocess.run(cmd, check=True, universal_newlines=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)

        # Extract timings from stdout
        for line in process.stdout.splitlines():
            if "hashGroupTime:" in line:
                total_hashGroupTime += float(line.split(":")[1].strip())
            elif "compareTreeTime:" in line:
                total_compareTreeTime += float(line.split(":")[1].strip())
                
        # If there's any error, append it to the file
        if process.stderr:
            output_file.write(f"--- Error for -hash-workers {hash_workers_value} ---\n")
            output_file.write(process.stderr + "\n")

    # Calculate average timings
    avg_hashGroupTime = total_hashGroupTime / runs
    avg_compareTreeTime = total_compareTreeTime / runs
    
    # Append the average timings to the file
    output_file.write(f"--- -hash-workers: {hash_workers_value} -data-workers: {data_workers_value} -comp-workers: {comp_workers_value} -input: {input_value} ----\n")
    output_file.write(f"Average hashGroupTime: {avg_hashGroupTime:.6f}\n")
    output_file.write(f"Average compareTreeTime: {avg_compareTreeTime:.6f}\n")

if __name__ == "__main__":
    hash_worker_values = [1,2,4,8,16,32,64,128]
    input_values = ["coarse.txt"]
    comp_workers_value = [1,2,4,8,16,32,64,128]
    with open("combined_output_data_coast.txt", 'w') as output_file:
        for value in hash_worker_values:
            for input_val in input_values:
                run_command(hash_workers_value=16, data_workers_value=value, input_value=input_val, output_file=output_file)
