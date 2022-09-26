from locust import HttpUser, task

class APILoadTesting(HttpUser):
    @task
    def userGet(self):
        self.client.get("/user/anhhd?profile=true&achievements=true&annual-ranks=true")

    @task
    def featuredUserList(self):
        self.client.get("/users/featured")

    @task
    def postList(self):
        self.client.get("/posts/")

    @task
    def eventList(self):
        self.client.get("/events/")
