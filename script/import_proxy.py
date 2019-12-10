#!/usr/bin/python3.6

import sys, getopt
import csv
import json
import xml
import urllib.parse
import urllib.request
import datetime
import random


def post( hostport, url, data ):
    urlpath = f"{hostport}{url}"
    headers = { 'accept' : 'application/json', 'Content-Type': 'application/json' }
    body = bytes(json.dumps(data), 'utf8')
    req = urllib.request.Request(urlpath, body, headers)
    response = urllib.request.urlopen(req)
    result = response.read()
    d=result.decode("utf8")
    j = json.loads(d)
    return j

def import_proxy(hostport,recordfile,tag):
    data = {}
    with open(recordfile, 'r') as csvfile:
        reader = csv.reader(csvfile, delimiter=':')
        records = []
        first = True
        for row in reader:
            host = row[0]
            port = int(row[1])
            record = {}
            record['host'] = host
            record['port'] = port
            records.append(record)
        data["proxies"] = records
        data["tag"] = tag
        rsp = post(hostport, "/api/v1/proxies", data)
        print(rsp)
   
def main(argv):
   help = "import_proxy.py -i <ifile> -u <url> -t <tag>"
   try:
      opts, args = getopt.getopt(argv,"hi:u:t:",["ifile=","url=","tag="])
   except getopt.GetoptError:
      print(help)
      sys.exit(2)
   for opt, arg in opts:
      if opt == '-h':
         print(help)
         sys.exit()
      elif opt in ("-i", "--ifile"):
         inputfile = arg
      elif opt in ("-u", "--url"):
         url = arg
      elif opt in ("-t", "--tag"):
         tag = arg
   import_proxy(url, inputfile, tag)

if __name__ == "__main__":
   main(sys.argv[1:])
