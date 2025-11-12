"""Module to handle slack calling functions"""
from time import sleep
from requests import Session

def upload_emoji(session:Session, emoji_name:str, filename:str):
    """uploads the emoji data to slack"""
    data = {
        'mode': 'data',
        'name': emoji_name,
        'token': session.api_token
    }

    while True:
        with open(filename, 'rb') as f:
            files = {'image': f}
            resp = session.post(
                session.url_add,
                data=data,
                files=files,
                allow_redirects=False,
                verify=False)

            if resp.status_code == 429:
                wait = int(resp.headers.get('retry-after', 1))
                print(f"429 Too Many Requests!, sleeping for {wait} seconds")
                sleep(wait)
                continue

        resp.raise_for_status()

        # Slack returns 200 OK even if upload fails, so check for status.
        response_json = resp.json()
        if not response_json['ok']:
            print(f"Error with uploading {emoji_name}: {response_json}")

        break

def get_current_emoji_list(session:Session, all=False):
    """List currently uploaded emoji to filter on"""
    page = 1
    result = []
    while True:
        data = {
            'query': '',
            'page': page,
            'count': 1000,
            'token': session.api_token
        }
        resp = session.post(session.url_list, data=data, verify=False)
        resp.raise_for_status()
        response_json = resp.json()

        if not all:
            result.extend(map(lambda e: e["name"], response_json["emoji"]))
            if page >= response_json["paging"]["pages"]:
                break
        else:
            result.extend(response_json["emoji"])
            if page >= response_json["paging"]["pages"]:
                break

        page = page + 1
    return result
