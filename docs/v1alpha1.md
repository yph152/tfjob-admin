<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [TFJob-admin](#tfjob-admin)
  - [API Authentication](#api-authentication)
  - [Endpoint](#endpoint)
- [API Summary](#api-summary)
  - [Version & Health](#version--health)
  - [TFJob](#tfjob)
- [API Details](#api-details)
  - [list-a-tfjobs](#list-a-tfjobs)
  - [create-a-tfjob](#create-a-tfjob)
  - [get-a-tfjob](#get-a-tfjob)
  - [update-a-tfjob](#update-a-tfjob)
  - [delete-a-tfjob](#delete-a-tfjob)
- [API Objects](#api-objects)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# TFJob-admin

描述 Clever 产品 tfjob-admin 组件的业务 API。

## API Authentication

* OIDC for `/v1alpha1`

## Endpoint


# API Summary

## Version & Health

| API | Path             | Detail               |
| --- | ---------------- | -------------------- |
| Get | GET /api/version | Get version object   |
| Get | GET /api/ping    | Diagnostic JSON info |

## TFJob
| API    | Path                                                                      | Detail                                          |
| ------ | ------------------------------------------------------------------------- | ----------------------------------------------- |
| List   | GET /api/v1alpha1/clusters/{cid}/partitions/{partition}/tfjobs            | [link](#list-of-tfjobs-of-all-users)        |
| Create | POST /api/v1alpha1/clusters/{cid}/partitions/{partition}/tfjobs           | [link](#create-a-tfjob)                     |
| Get    | GET /api/v1alpha1/clusters/{cid}/partitions/{partition}/tfjobs/:ID        | [link](#get-a-tfjob)                        |
| Update | PATCH /api/v1alpha1/clusters/{cid}/partitions/{partition}/tfjobs/:ID      | [link](#update-a-tfjob)                     |
| Delete | DELETE /api/v1alpha1/clusters/{cid}/partitions/{partition}/tfjobs/:ID     | [link](#delete-a-tfjob)                     |
# API Details

## list-a-tfjobs

**Request**

URL: `GET /api/v1alpha2/clusters/{cid}/partitions/{partition}/tfjob`

Header:

```
X-User: admin
```

Args:

| Name      | Type             | Detail                  |
| --------- | ---------------- | ----------------------- |
| cid       | string, required | Cluster id              |
| partition | string, required | Partition(namespace) id |

**Response**

```
200 OK

{
  "metadata": {
    "total"; 1
  },
   "items": [
     {
       "metadata": {
       "uid": "admin",
       "id": "12345678",
       "name": "tfserving-0",
       "creationTime": "2017-08-23T21:17:09.379754063+08:00",
       "lastUpdateTime": "2017-08-23T21:17:09.379754063+08:00"
     },
      "spec": {
      "description": "create a tfjobs",
      "partition": "clever",
      "runtime": "Python_2.7+Tensorflow_1.2.1",
      "workdir": "/clever/admin/xinhe/jobs",
      "bootfile": "/clever/admin/xinhe/file.py",
      "args": ["--test=10","--config=1"],
      "checkpointDir": /clever/admin",
      "eventDir": "/clever/admin/xinhe",
      "modelDir": "/clever/admin/xinhe/logs",
      "logDir": "/clever/admin/xinhe/logs/guideline",
      "replicas": {[
        {
          "type": "worker",
          "count": "1",
          "resource": {
             "cpu": "1.000",
             "memory": "2048Mi",
             "gpu": "1"
             }
         },
         {
           "type": "ps",
           "conut": "0",
           "resource": {
             "cpu": "1.000",
             "memory": "2048Mi",
             "gpu": "1"
           }           
          }]
        }
     },
    "status": {
      "phase": "running",
      "replicas": [
        {
          "id": "12345678-k8s-xyz",
          "restartCount": 0,
          "status": "Running",
          "startTime": "2017-08-23T21:17:09.379754063+08:00"
        }
      ]
     }
   }
 ]
}

```
## create-a-tfjob

**Request**

URL: `POST /api/v1alpha1/clusters/{cid}/partitions/{partition}/tfjob`

Header:

```
X-User: admin
```

Args:

| Name      | Type             | Detail                  |
| --------- | ---------------- | ----------------------- |
| cid       | string, required | Cluster id              |
| partition | string, required | Partition(namespace) id |
| name      | string, required | tfjob name           
| description | string, optional | description tfjob       |
| runtime   | string, required | python and tensorflow version |
| workdir   | string, optional| work directory     |
| bootfile  | string, required | startup script          |
| args      | []string, required | args of command         |
| enviroment| string, optional | enviroment
| checkpointDir  | string, optional | checkpoint directory.   |
| eventDir  | string, optional | event dir               |
| modelDir  | string, optional | model dir               |
| logDir    | string, optional | log dir
| cpu       | string, required | cpu-resource-limit      |
| gpu       | string, required | gpu-resource-limit      |
| memory    | string, required | memory-resource-limit   |
| type      | string, optional | type for tensorflow ps/worker |
| count     | string, optional | ps/worker counts        | 

**Response**

```
200 ok

{
  "metadata": {
    "uid": "admin",
    "id": "12345678",
    "name": "tfserving-0",
    "creationTime": "2017-08-23T21:17:09.379754063+08:00",
    "lastUpdateTime": "2017-08-23T21:17:09.379754063+08:00"
  },
 "spec": {
    "description": "create a tfjobs",
    "partition": "clever",
    "runtime": "Python 2.7 + Tensorflow 1.2.1",
    "workerdir": "/clever/admin/xinhe/logs",
    "bootfile": "/clever/admin/xinhe/file.python",
    "args": "-test=10 -config=1",
    "checkpointDir": /clever/admin",
    "eventDir": "/clever/admin/xinhe",
    "modelDir": "/clever/admin/xinhe/logs",
    "logDir": "/clever/admin/xinhe/logs/guideline",
    "replicas": {
    "count": "1",
    "replicas": {[
      {
        "type": "worker",
        "count": "1",
        "resource": {
           "cpu": "1.000",
           "memory": "2048Mi",
           "gpu": "1"
           }
       },
       {
         "type": "ps",
         "conut": "0",
         "resource": {
           "cpu": "1.000",
           "memory": "2048Mi",
           "gpu": "1"
         }           
        }]
      }
   },
  "status": {
    "phase": "running",
    "replicas": [
      {
        "id": "12345678-k8s-xyz",
        "restartCount": 0,
        "status": "Running",
        "startTime": "2017-08-23T21:17:09.379754063+08:00"
      }
    ]
  }
}
```

## get-a-tfjob

**Request**

URL: `GET /api/v1alpha2/clusters/{cid}/partitions/{partition}/tfservings/:ID`

Header:

```
X-User: admin
```

Args:

| Name      | Type             | Detail                  |
| --------- | ---------------- | ----------------------- |
| cid       | string, required | Cluster id              |
| partition | string, required | Partition(namespace) id |
| id        | string, required | tfjob-id            |

**Response**

```
200 ok

{
  "metadata": {
    "uid": "admin",
    "id": "12345678",
    "name": "tfserving-0",
    "creationTime": "2017-08-23T21:17:09.379754063+08:00",
    "lastUpdateTime": "2017-08-23T21:17:09.379754063+08:00"
  },
 "spec": {
    "description": "create a tfjobs",
    "partition": "clever",
    "runtime": "Python 2.7 + Tensorflow 1.2.1",
    "workerdir": "/clever/admin/xinhe/logs",
    "bootfile": "/clever/admin/xinhe/file.python",
    "args": "-test=10 -config=1",
    "checkpointDir": /clever/admin",
    "eventDir": "/clever/admin/xinhe",
    "modelDir": "/clever/admin/xinhe/logs",
    "logDir": "/clever/admin/xinhe/logs/guideline",
    "replicas": {
    "count": "1",
    "replicas": {[
      {
        "type": "worker",
        "count": "1",
        "resource": {
           "cpu": "1.000",
           "memory": "2048Mi",
           "gpu": "1"
           }
       },
       {
         "type": "ps",
         "conut": "0",
         "resource": {
           "cpu": "1.000",
           "memory": "2048Mi",
           "gpu": "1"
         }           
        }]
      }
   },
  "status": {
    "phase": "running",
    "replicas": [
      {
        "id": "12345678-k8s-xyz",
        "restartCount": 0,
        "status": "Running",
        "startTime": "2017-08-23T21:17:09.379754063+08:00"
      }
    ]
  }
}
```

## update-a-tfjob

**Request**

URL: `PATCH /api/v1alpha2/clusters/{cid}/partitions/{partition}/tfservings/:ID`

Header:

```
X-User: admin
```

Args:

| Name      | Type             | Detail                  |
| --------- | ---------------- | ----------------------- |
| cid       | string, required | Cluster id              |
| partition | string, required | Partition(namespace) id |
| name      | string, required | tfjob name           
| description | string, optional | description tfjob       |
| runtime   | string, required | python and tensorflow version |
| workdir   | string, optional| work directory     |
| bootfile  | string, required | startup script          |
| args      | []string, required | args of command         |
| enviroment| string, optional | enviroment
| checkpointDir  | string, optional | checkpoint directory.   |
| eventDir  | string, optional | event dir               |
| modelDir  | string, optional | model dir               |
| logDir    | string, optional | log dir
| cpu       | string, required | cpu-resource-limit      |
| gpu       | string, required | gpu-resource-limit      |
| memory    | string, required | memory-resource-limit   |
| type      | string, optional | type for tensorflow ps/worker |
| count     | string, optional | ps/worker counts        |
**Response**

```
200 ok

{
  "metadata": {
    "uid": "admin",
    "id": "12345678",
    "name": "tfserving-0",
    "creationTime": "2017-08-23T21:17:09.379754063+08:00",
    "lastUpdateTime": "2017-08-23T21:17:09.379754063+08:00"
  },
 "spec": {
    "description": "create a tfjobs",
    "partition": "clever",
    "runtime": "Python 2.7 + Tensorflow 1.2.1",
    "workerdir": "/clever/admin/xinhe/logs",
    "bootfile": "/clever/admin/xinhe/file.python",
    "args": "-test=10 -config=1",
    "checkpointDir": /clever/admin",
    "eventDir": "/clever/admin/xinhe",
    "modelDir": "/clever/admin/xinhe/logs",
    "logDir": "/clever/admin/xinhe/logs/guideline",
    "replicas": {
    "count": "1",
    "replicas": {[
      {
        "type": "worker",
        "count": "1",
        "resource": {
           "cpu": "1.000",
           "memory": "2048Mi",
           "gpu": "1"
           }
       },
       {
         "type": "ps",
         "conut": "0",
         "resource": {
           "cpu": "1.000",
           "memory": "2048Mi",
           "gpu": "1"
         }           
        }]
      }
   },
  "status": {
    "phase": "running",
    "replicas": [
      {
        "id": "12345678-k8s-xyz",
        "restartCount": 0,
        "status": "Running",
        "startTime": "2017-08-23T21:17:09.379754063+08:00"
      }
    ]
  }
}
```

## delete-a-tfjob

**Request**

URL: `DELETE /api/v1alpha2/clusters/{cid}/partitions/{partition}/tfservings/:ID`

Header:

```
X-User: admin
```

Args:

| Name      | Type             | Detail                  |
| --------- | ---------------- | ----------------------- |
| cid       | string, required | Cluster id              |
| partition | string, required | Partition(namespace) id |
| id        | string, required | tfjob-id            |

**Response**

```
200 ok
```


# API Objects
