import re
import collections
import time

def match_lines_and_count_servers():
    input_file_name = "input.txt"

    pattern = re.compile(r'gtalk:(talk|chat)\?jid=[a-z0-9]+@([a-z0-9]+)\.[a-z]{1,5}', re.IGNORECASE)

    server_counts = collections.Counter()

    try:
        with open(input_file_name, 'r') as file:
            start = time.time()
            for line in file:
                match = pattern.search(line)
                if match:
                    server_name = match.group(2)
                    server_counts[server_name] += 1
            overall = time.time() - start
        output_file_name = "regex_res.txt"

        with open(output_file_name, 'w') as file:
            for server, count in server_counts.items():
                file.write(f"{server}: {count}\n")

        print(f"Regex - Elapsed time: {overall}")

    except FileNotFoundError:
        print(f"The file {input_file_name} does not exist")

if __name__ == "__main__":
    match_lines_and_count_servers()