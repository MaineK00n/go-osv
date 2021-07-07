import argparse
import logging
from typing import Tuple
from deepdiff import DeepDiff
import requests
from requests.adapters import HTTPAdapter
from urllib3.util import Retry
from urllib.parse import quote
import pprint
from concurrent.futures import ThreadPoolExecutor
import os


def diff_response(args: Tuple[str, str, str]):
    # Endpoint
    # /ids/:id
    # /:type/ids/:id
    # /pkgs/:name
    # /:type/pkgs/:name

    path = ''
    if args[0] == 'id':
        if args[1] == 'all':
            path = f'ids/{args[2]}'
        else:
            path = f'{args[1]}/ids/{args[2]}'
    if args[0] == 'package':
        if args[1] == 'all':
            path = f'pkgs/{args[2]}'
        else:
            path = f'{args[1]}/pkgs/{args[2]}'

    session = requests.Session()
    retries = Retry(total=5,
                    backoff_factor=1,
                    status_forcelist=[503, 504])
    session.mount("http://", HTTPAdapter(max_retries=retries))

    try:
        response_old = requests.get(
            f'http://127.0.0.1:1325/{path}', timeout=(10.0, 10.0)).json()
        response_new = requests.get(
            f'http://127.0.0.1:1326/{path}', timeout=(10.0, 10.0)).json()
    except requests.ConnectionError as e:
        logger.error(f'Failed to Connection..., err: {e}')
        exit(1)
    except Exception as e:
        logger.error(f'Failed to GET request..., err: {e}')
        exit(1)

    diff = DeepDiff(response_old, response_new, ignore_order=True)
    if diff != {}:
        logger.warning(
            f'There is a difference between old and new(or RDB and Redis):\n {pprint.pformat({"mode": args[0], "args": args, "diff": diff}, indent=2)}')


parser = argparse.ArgumentParser()
parser.add_argument('mode', choices=['id', 'package'],
                    help='Specify the mode to test.')
parser.add_argument('fetchtype', choices=['all', 'crates.io', 'DWF', 'Go', 'Linux', 'OSS-Fuzz', 'PyPI'],
                    help='Specify the Fetch Type to be started in server mode when testing.')
parser.add_argument(
    '--debug', action=argparse.BooleanOptionalAction, help='print debug message')
args = parser.parse_args()

logger = logging.getLogger(__name__)
stream_handler = logging.StreamHandler()

if args.debug:
    logger.setLevel(logging.DEBUG)
    stream_handler.setLevel(logging.DEBUG)
else:
    logger.setLevel(logging.INFO)
    stream_handler.setLevel(logging.INFO)

formatter = logging.Formatter(
    '%(levelname)s[%(asctime)s] %(message)s', "%m-%d|%H:%M:%S")
stream_handler.setFormatter(formatter)
logger.addHandler(stream_handler)

logger.info(
    f'start server mode test(mode: {args.mode}, os: {args.fetchtype})')

list_path = None
if args.mode == 'id':
    list_path = f"integration/id/{args.fetchtype}.txt"
if args.mode == 'package':
    list_path = f"integration/package/{args.fetchtype}.txt"

if not os.path.isfile(list_path):
    logger.error(f'Failed to find list path..., list_path: {list_path}')
    exit(1)

with open(list_path) as f:
    list = [s.strip() for s in f.readlines()]
    with ThreadPoolExecutor() as executor:
        ins = ((args.mode, args.fetchtype, e) for e in list)
        executor.map(diff_response, ins)
