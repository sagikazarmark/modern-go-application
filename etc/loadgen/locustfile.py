#!/usr/bin/python

import json
import random
import urllib3
from locust import HttpLocust, TaskSet

clientErrorCodes = [400, 403, 404]
serverErrorCodes = [500, 503]


def index(l):
    l.client.get("/")

def clientErrors(l):
    l.client.get("/httpbin/status/" + str(random.choice(clientErrorCodes)))

def serverErrors(l):
    l.client.get("/httpbin/status/" + str(random.choice(serverErrorCodes)))

def responseSize(l):
    l.client.get("/httpbin/bytes/" + str(random.randint(100, 500000)))

class DemoBehavior(TaskSet):
    def on_start(self):
        self.client.verify = False
        urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
        index(self)

    tasks = {
        index: 3,
        clientErrors: 6,
        serverErrors: 4,
        responseSize: 2,
    }


class WebsiteUser(HttpLocust):
    task_set = DemoBehavior
    min_wait = 1000
    max_wait = 10000
