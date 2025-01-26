import random
import string


def gen_servername() -> str:
    return ''.join(random.choices(string.ascii_letters + string.digits, k=random.randint(3, 15)))

def gen_string(servername="") -> str:
    """
    Generates random line that matches template specified in the task.
    :param servername: Server name. If not specified, a random one will be generated.
    """

    contact_type = random.choice(["talk", "chat"])

    username = ''.join(random.choices(string.digits + string.ascii_letters, k=random.randint(3, 15)))

    if servername == "":
        servername = gen_servername()

    zone_name = ''.join(random.choices(string.ascii_letters, k=random.randint(1, 5)))

    output = f"gtalk:{contact_type}?jid={username}@{servername}.{zone_name}"
    return output


def gen_faulty_string(servername=""):
    """
    Generates string that is incorrect in at least one field
    """
    was_incorrect = False  # not evenly distributed but who cares

    if random.random() > 0.5:
        was_incorrect = True
        contact_type = "incorrect"
    else:
        contact_type = random.choice(["talk", "chat"])

    if random.random() > 0.5:
        was_incorrect = True
        username = ''.join(random.choices(string.punctuation, k=random.randint(0, 15)))
    else:
        username = ''.join(random.choices(string.digits + string.ascii_letters, k=random.randint(3, 15)))

    if random.random() > 0.5:
        was_incorrect = True
        servername = ''.join(random.choices(string.punctuation, k=random.randint(0, 15)))
    else:
        if servername == "":
            servername = gen_servername()

    if random.random() > 0.5 or not was_incorrect:
        zone_name = ''.join(random.choices(string.ascii_letters, k=random.randint(6, 10)))
    else:
        zone_name = ''.join(random.choices(string.ascii_letters, k=random.randint(1, 5)))

    output = f"gtalk:{contact_type}?jid={username}@{servername}.{zone_name}"
    return output

if __name__ == '__main__':
    print(gen_string())
    print(gen_faulty_string())
