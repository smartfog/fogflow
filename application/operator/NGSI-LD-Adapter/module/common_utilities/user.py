import requests

USERS_URL = 'http://jsonplaceholder.typicode.com/use'


def get_users():
    """Get list of users"""
    response = requests.get(USERS_URL)
    print(response.status_code)
    if response.ok:
        return response
    else:
        return None
