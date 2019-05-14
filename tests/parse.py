import os
import re
import json
import logging
from datetime import datetime, timedelta

def totimestamp(dt, epoch=datetime(1970,1,1)):
    dt = datetime.strptime(dt, '%Y-%m-%dT%H:%M:%SZ')
    td = dt - epoch
#     return td.total_seconds()
    return (td.microseconds + (td.seconds + td.days * 86400) * 10**6) / 10**6 

def getTime(rxlist):
    for i in range(len(rxlist)):
        time = rxlist[i]['timestamp'] - rxlist[i]['data']['timestamp']
        print("message: " + rxlist[i]['data']['id'])
        print("time from send to receive: " + str(time/1000000000))

def getMsgType(msg):
    print(msg)
    if msg['msg'] == "Received a message":
        return "received from: " + msg['from'] + "\n" + "at: " + str(totimestamp(msg['time']))
    if msg['msg'] == "Sending a message":
        return "sent message"
    
def parseMsgType(msgList, index):
    Rx = []
    Sn = []
    for i in range(len(msgList)):
        try:
            json_object = json.loads(msgList[i])
            if(json_object['msg']=="Received a message"):
                Rx.append(json_object)
            elif(json_object['msg']=="Sending a message"):
                Sn.append(json_object)
            else:
                continue
        except:
            continue
    lRx[index] = Rx
    lSn[index] = Sn
    

dirname = "./data/"

l = [[]]*100
lSn = [[]]*100
lRx = [[]]*100

for filename in os.listdir(dirname):
    if filename.endswith(".txt"):
        i = re.findall(r'\d+', filename)
        with open(dirname + filename, "r") as f:
            l[int(i[0])] = list(f)
        continue
    else:
        continue

for i in range(len(l)):
    parseMsgType(l[i][:], i)

for i in range(len(lRx)):
    print("node" + str(i))
    getTime(lRx[i])