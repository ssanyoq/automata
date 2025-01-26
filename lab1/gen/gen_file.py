import random
import gen_fields

def gen_file(path="input.txt", n=100, error_chance=0.1):
    servers_amt = random.randint(1, n)
    servers = [gen_fields.gen_servername() for _ in range(servers_amt)]
    faulty = 0
    with open(path, "w") as file:
        for i in range(n):
            if random.random() < error_chance:
                s = gen_fields.gen_faulty_string(random.choice(servers))
                faulty += 1
            else:
                s = gen_fields.gen_string(random.choice(servers))
            file.write(s + "\n")
    print(f"Generated file at {path} with:")
    print(f"{n} entries")
    print(f"{servers_amt} servers")
    print(f"{faulty} faulty entries")

if __name__ == '__main__':
    gen_file()