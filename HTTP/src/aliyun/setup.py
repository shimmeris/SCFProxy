import argparse
import base64
import os
import sys
import zipfile

from alibabacloud_fc_open20210406.client import Client as FC_Open20210406Client
from alibabacloud_fc_open20210406.models import Code
from alibabacloud_tea_openapi import models as open_api_models
from alibabacloud_fc_open20210406 import models as fc__open_20210406_models
from alibabacloud_tea_util import models as util_models

# 修改下面三个变量
ACCOUNT_ID = ""
ACCESS_KEY_ID = ""
ACCESS_KEY_SECRET = ""

# 以下可以不用修改
SERVICE_NAME = "aliyunfc_http_proxy"
FUNCTION_NAME = "http_proxy"
TRIGGER_NAME = "trigger"

domestic_areas = [
    "cn-qingdao",
    "cn-beijing",
    "cn-zhangjiakou",
    "cn-huhehaote",
    "cn-hangzhou",
    "cn-shanghai",
    "cn-shenzhen",
    "cn-chengdu"
]
foreign_areas = [
    "cn-hongkong",
    "ap-northeast-1",
    "ap-southeast-1",
    "ap-southeast-2",
    "ap-southeast-3",
    "ap-southeast-5",
    "us-east-1",
    "us-west-1",
    "eu-west-1",
    "eu-central-1",
    "ap-south-1"
]
areas_dict = {
    "domestic": domestic_areas,
    "foreign": foreign_areas,
    "all": domestic_areas + foreign_areas,
}


def create_client(endpoint, account_id=ACCOUNT_ID, access_key_id=ACCESS_KEY_ID, access_key_secret=ACCESS_KEY_SECRET):
    config = open_api_models.Config(
        access_key_id=access_key_id,
        access_key_secret=access_key_secret
    )
    config.endpoint = f'{account_id}.{endpoint}.fc.aliyuncs.com'
    return FC_Open20210406Client(config)


# 创建服务
def create_service(client):
    create_service_headers = fc__open_20210406_models.CreateServiceHeaders()
    create_service_request = fc__open_20210406_models.CreateServiceRequest()
    create_service_request.service_name = SERVICE_NAME
    runtime = util_models.RuntimeOptions()
    return client.create_service_with_options(create_service_request, create_service_headers, runtime)


# 获取 base64 编码的代码包
def get_zip():
    with zipfile.ZipFile("code.zip", "w", zipfile.ZIP_DEFLATED) as f:
        f.write("server.py")

    with open("code.zip", "rb") as f:
        data = f.read()

    try:
        os.remove("code.zip")
    except Exception as e:
        print(e)

    return base64.b64encode(data).decode("utf-8")


# 创建函数
def create_function(zip_code, client):
    create_function_headers = fc__open_20210406_models.CreateFunctionHeaders()
    create_function_request = fc__open_20210406_models.CreateFunctionRequest()
    runtime = util_models.RuntimeOptions()
    create_function_request.function_name = FUNCTION_NAME
    create_function_request.code = Code(zip_file=zip_code)
    create_function_request.runtime = "python3.9"
    create_function_request.handler = "server.handler"
    return client.create_function_with_options(SERVICE_NAME, create_function_request, create_function_headers, runtime)


# 创建 http 触发器
def create_trigger(client):
    create_trigger_headers = fc__open_20210406_models.CreateTriggerHeaders()
    create_trigger_request = fc__open_20210406_models.CreateTriggerRequest()
    runtime = util_models.RuntimeOptions()
    create_trigger_request.trigger_type = "http"
    create_trigger_request.trigger_name = TRIGGER_NAME
    create_trigger_request.trigger_config = '{"authType": "anonymous", "methods": ["GET", "POST"]}'
    return client.create_trigger_with_options(SERVICE_NAME, FUNCTION_NAME, create_trigger_request,
                                              create_trigger_headers, runtime)


# 删除服务
def delete_service(client):
    delete_service_headers = fc__open_20210406_models.DeleteServiceHeaders()
    runtime = util_models.RuntimeOptions()
    return client.delete_service_with_options(SERVICE_NAME, delete_service_headers, runtime)


# 删除函数
def delete_function(client):
    delete_function_headers = fc__open_20210406_models.DeleteFunctionHeaders()
    runtime = util_models.RuntimeOptions()
    return client.delete_function_with_options(SERVICE_NAME, FUNCTION_NAME, delete_function_headers, runtime)


# 删除触发器
def delete_trigger(client):
    delete_trigger_headers = fc__open_20210406_models.DeleteTriggerHeaders()
    runtime = util_models.RuntimeOptions()
    return client.delete_trigger_with_options(SERVICE_NAME, FUNCTION_NAME, TRIGGER_NAME, delete_trigger_headers,
                                              runtime)


def remove(client):
    try:
        print(delete_trigger(client))
    except Exception as e:
        print(e)

    try:
        print(delete_function(client))
    except Exception as e:
        print(e)

    try:
        print(delete_service(client))
    except Exception as e:
        print(e)


def install(client, zip_code):
    try:
        print(create_service(client))
        print(create_function(zip_code, client))
        ret = create_trigger(client)
        print(ret)
        return ret.body.url_internet
    except Exception as e:
        print(e)


def get_parser():
    parser = argparse.ArgumentParser(
        description="""阿里云函数 HTTP 代理一键配置

# 部署单个城市
python setup.py install -c cn-shanghai

# 部署区域内所有城市
python setup.py install -a domestic

# 删除所有通过 setup.py 部署的云函数
python setup.py delete

# 强制删除所有区域的云函数
python setup.py force_delete

建议：
1. 大陆外地区部署的云函数延迟较高，推荐只使用大陆的
2. 随用随装，用完删除
""",
        add_help=False,
        usage="python3 %(prog)s action [-c city] [-a area] [-h] [-f]",

        formatter_class=argparse.RawTextHelpFormatter,
    )
    parser.add_argument(
        "action",
        nargs="?",
        choices=("install", "delete", "force_delete"),
        default="install",
        metavar="action",
        help="install 或 delete 或 force_delete",
    )
    parser.add_argument(
        "-h", "--help", action="help", default=argparse.SUPPRESS, help="展示帮助信息"
    )
    parser.add_argument(
        "-c", "--city", dest="city", metavar="city",
        help=f"""云函数部署城市

可选城市:
    大陆地区: {", ".join(domestic_areas)}
    非大陆地区: {", ".join(foreign_areas)}

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

        zip_code = get_zip()
        with open("cities.txt", "a") as f:
            if args.area in ["all", "domestic", "foreign"]:
                for city in areas_dict[args.area]:
                    client = create_client(city)
                    trigger = install(client, zip_code)
                    if trigger:
                        f.write(f"{city} {trigger}\n")
                        print(f"{city} 区域部署成功")
                    else:
                        print(f"{city} 区域部署失败")
                        remove(client)

            elif args.city in areas_dict['all']:

                client = create_client(args.city)
                trigger = install(client, zip_code)
                if trigger:
                    f.write(f"{args.city} {trigger}\n")
                    print(f"{args.city} 区域部署成功")
                else:
                    print(f"{args.city} 区域部署失败")
                    remove(client)

            else:
                print(f"请输入有效的城市或区域")
                exit(-1)

    elif args.action == "delete":
        with open("cities.txt", "r") as f:
            for line in f.readlines():
                client = create_client(line.split()[0])
                remove(client)

                print(f"{line.split()[0]} 区域卸载成功")
        os.remove("cities.txt")
    elif args.action == "force_delete":
        for area in areas_dict['all']:
            client = create_client(area)
            remove(client)

            print(f"{area} 区域卸载成功")
        try:
            os.remove("cities.txt")
        except Exception as e:
            pass
