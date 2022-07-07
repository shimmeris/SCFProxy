# -*- coding: utf8 -*-
import sys
import json
import time
import base64
import zipfile
import argparse
from pathlib import Path
from concurrent.futures import ThreadPoolExecutor

from tencentcloud.common import credential
from tencentcloud.scf.v20180416 import scf_client, models
from tencentcloud.common.profile.http_profile import HttpProfile
from tencentcloud.common.profile.client_profile import ClientProfile
from tencentcloud.common.exception.tencent_cloud_sdk_exception import (
    TencentCloudSDKException,
)


# 填入腾讯云的 SecretId 和 SecretKey
SecretId = ""
SecretKey = ""


domestic_areas = [
    "ap-beijing",
    "ap-chengdu",
    "ap-guangzhou",
    "ap-shanghai",
]
foreign_areas = [
    "ap-hongkong",
    "ap-mumbai",
    "ap-singapore",
    "ap-bangkok",
    "ap-seoul",
    "ap-tokyo",
    "eu-frankfurt",
    "na-ashburn",
    "na-toronto",
    "na-siliconvalley",
]
areas_dict = {
    "domestic": domestic_areas,
    "foreign": foreign_areas,
    "all": domestic_areas + foreign_areas,
}


def get_zip():
    with zipfile.ZipFile("code.zip", "w", zipfile.ZIP_DEFLATED) as f:
        f.write("server.py")

    with open("code.zip", "rb") as f:
        data = f.read()
    return base64.b64encode(data).decode("utf-8")


def remove_file(filename):
    f = Path(filename)
    f.unlink(missing_ok=True)


def create_client(city):
    cred = credential.Credential(SecretId, SecretKey)
    httpProfile = HttpProfile()
    httpProfile.endpoint = "scf.tencentcloudapi.com"

    clientProfile = ClientProfile()
    clientProfile.httpProfile = httpProfile
    return scf_client.ScfClient(cred, city, clientProfile)


def create_scf(client, function_name):
    try:
        req = models.CreateFunctionRequest()
        params = {
            "FunctionName": function_name,
            "Code": {
                "ZipFile": get_zip(),
            },
            "Handler": "server.main_handler",
            "Runtime": "Python3.6",
            "Timeout": 60
        }

        req.from_json_string(json.dumps(params))
        client.CreateFunction(req)

    except TencentCloudSDKException as err:
        print(err)


def create_trigger(client, function_name):
    req = models.CreateTriggerRequest()
    params = {
        "FunctionName": function_name,
        "TriggerName": "http_trigger",
        "Type": "apigw",
        "TriggerDesc": """{
            "api":{
                "authRequired":"FALSE",
                "requestConfig":{
                    "method":"ANY"
                },
                "isIntegratedResponse":"TRUE"
            },
            "service":{
                "serviceName":"SCF_API_SERVICE"
            },
            "release":{
                "environmentName":"release"
            }
        }""",
    }
    req.from_json_string(json.dumps(params))
    resp = client.CreateTrigger(req)

    data = json.loads(resp.to_json_string())
    trigger_url = json.loads(data["TriggerInfo"]["TriggerDesc"])["service"]["subDomain"]
    return trigger_url


def delete_scf(city):
    client = create_client(city)
    try:
        req = models.DeleteFunctionRequest()
        params = {"FunctionName": f"http_{city}", "Namespace": "default"}
        req.from_json_string(json.dumps(params))
        client.DeleteFunction(req)
        print(f"{city} 区域云函数删除成功")
    except TencentCloudSDKException as err:
        print(err)


def install(city):
    function_name = f"http_{city}"
    client = create_client(city)
    create_scf(client, function_name)
    time.sleep(5)
    while True:
        try:
            trigger_url = create_trigger(client, function_name)
        except TencentCloudSDKException as err:
            if err.code == "FailedOperation":
                time.sleep(20)
        else:
            if trigger_url:
                print(f"{city} 区域云函数创建成功")
                break
    return city, trigger_url


def get_parser():
    parser = argparse.ArgumentParser(
        description="""腾讯云函数 HTTP 代理一键配置

# 部署单个城市
python setup.py install -c ap-beijing

# 部署区域内所有城市
python setup.py install -a domestic

# 删除所有通过 setup.py 部署的云函数
python setup.py delete

建议：
1. 大陆外地区部署的云函数延迟较高，推荐只使用国内的
2. 随用随装，用完删除
""",
        add_help=False,
        usage="python3 %(prog)s action [-c city] [-a area] [-h]",
        
        formatter_class=argparse.RawTextHelpFormatter,
    )
    parser.add_argument(
        "action",
        nargs="?",
        choices=("install", "delete"),
        default="install",
        metavar="action",
        help="install 或 delete",
    )
    parser.add_argument(
        "-h", "--help", action="help", default=argparse.SUPPRESS, help="展示帮助信息"
    )
    parser.add_argument(
        "-c", "--city", dest="city", metavar="city", 
        help="""云函数部署城市

可选城市:
    大陆地区: ap-beijing, ap-chengdu, ap-guangzhou, ap-shanghai
    亚太地区: ap-hongkong, ap-mumbai, ap-singapore, ap-bangkok, ap-seoul, ap-tokyo
    欧洲地区: eu-frankfurt
    北美地区: na-siliconvalley, na-toronto, na-ashburn
        
"""
    )
    parser.add_argument(
        "-a", "--area",
        metavar="area",
        dest="area",
        help="""云函数部署区域（包含多个城市）

可选区域: 
    大陆地区: domestic
    非大陆地区: foreign
    所有地区: all
    """,
    )
    return parser


if __name__ == "__main__":
    parser = get_parser()
    if len(sys.argv) == 1:
        parser.print_help()
        exit()

    args = parser.parse_args()

    if args.action == 'install':
        if not (args.city or args.area):
            print(f"请输入城市或区域")
            exit()
        
        with open("cities.txt", "a") as f:
            if args.area in ["all", "domestic", "foreign"]:
                with ThreadPoolExecutor(max_workers=5) as executor:
                    results = list(executor.map(install, areas_dict[args.area]))

                for city, trigger in results:
                    if trigger:
                        f.write(f"{city} {trigger}\n")
            
            elif args.city in areas_dict['all']:
                city, trigger = install(args.city)
                if trigger:
                    f.write(f"{city} {trigger}\n")
            
            else:
                print(f"请输入有效的城市或区域")
                exit()

        remove_file("code.zip")


    elif args.action == "delete":
        with open("cities.txt", "r") as f:
            for line in f.readlines():
                delete_scf(line.split()[0])
        remove_file("cities.txt")
