import json
import string

from random import *
from locust import HttpUser, task, between


def random_string(length: int) -> str:
    return ''.join(choices(string.ascii_uppercase + string.digits, k=length))


class User(HttpUser):
    @task(10)
    def get_user(self):
        self.client.get("/user/" + str(randint(1, 10_000)), name="get_user")

    @task(1000)
    def create_user(self):
        payload = {
            "username": random_string(6),
            "first_name": random_string(4),
            "last_name": random_string(4),
            "email": random_string(10),
            "phone": random_string(11)
        }

        headers = {"content-type": "application/json"}

        self.client.post("/user", data=json.dumps(payload), headers=headers, name="create_user")

    @task(10)
    def delete_user(self):
        self.client.delete("/user/" + str(randint(1, 10_000)), name="delete_user")

    @task(10)
    def update_user(self):
        payload = {
            "username": random_string(6),
            "first_name": random_string(4),
            "last_name": random_string(4),
            "email": random_string(10),
            "phone": random_string(11)
        }

        headers = {"content-type": "application/json"}

        self.client.put(
            '/user/' + str(randint(1, 10_000)),
            data=json.dumps(payload),
            headers=headers,
            name="update_user"
        )

    wait_time = between(1, 5)