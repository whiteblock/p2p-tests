import os
import re
import json
import logging
from datetime import datetime, timedelta

def writeFile(data):
    f= open("aggregatedTimeData.txt","a+")
    f.write(str(data))
    f.close()

def getTime(rxlist):
    data = {}
    for i in range(len(rxlist)):
        msgid = rxlist[i]['data']['id']
        time = rxlist[i]['timestamp'] - rxlist[i]['data']['timestamp']
        data[msgid]= str(time/1000000000)
    return data

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

dirname = "./data/test3/"

l = [[]]*100
lSn = [[]]*100
lRx = [[]]*100

for filename in os.listdir(dirname):
    if filename.endswith(".log"):
        i = re.findall(r'\d+', filename)
        with open(dirname + filename, "r") as f:
            l[int(i[0])] = list(f)
        continue
    else:
        continue

for i in range(len(l)):
    parseMsgType(l[i][:], i)

for i in range(len(lRx)):
    writeFile(getTime(lRx[i]))