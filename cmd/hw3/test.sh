pip3 install locust

locust -f tests/locust_runner.py --headless -u 1000 -r 10 --run-time 5m --host http://arch.homework