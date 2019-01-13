#!/usr/bin/python

import json
import random
import urllib3
from locust import HttpLocust, TaskSet

clientErrorCodes = [400, 403, 404]
serverErrorCodes = [500, 503]


def index(l):
    l.client.get("/")

def helloWorld(l):
    l.client.get("/hello")

def sayHello(l):
    payload = {'who': 'John'}
    headers = {'content-type': 'application/json'}

    l.client.post("/hello", data=json.dumps(payload), headers=headers)

def clientErrors(l):
    l.client.get("/errors/" + str(random.choice(clientErrorCodes)))

def serverErrors(l):
    l.client.get("/errors/" + str(random.choice(serverErrorCodes)))


class DemoBehavior(TaskSet):
    def on_start(self):
        self.client.verify = False
        urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
        index(self)

    tasks = {
        index: 1,
        helloWorld: 6,
        sayHello: 6,
        clientErrors: 3,
        serverErrors: 2,
    }


class WebsiteUser(HttpLocust):
    task_set = DemoBehavior
    min_wait = 1000
    max_wait = 10000
