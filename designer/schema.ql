name: string @index(term) .
formattype: string @index(term) .
description: string @index(term) .
filepath: string @index(term) .
url: string @index(term) .
flavor: string @index(term) .
inputdata: string @index(term) .
version: string @index(term) .
attribute: string @index(term) .
internalType: string @index(term) .
updateAction: string @index(term) .


tag: string @index(term) .
hwType: string @index(term) .
osType: string @index(term) .
operatorName: string @index(term) .
prefetched: bool  . 
values: string @index(term) .
parameters: [uid] @reverse .

type DockerImage {
    name
    tag
    hwType
    osType
    operatorName
    prefetched
}


type Operator {
    name
    description
    parameters
}


id: string @index(term) .
topology: string @index(term) .
stype: string @index(term) .
qos: string @index(term) .
priority: string  .
geoscope: string . 

type ServiceIntent {
    id
    topology    
    stype
    priority
    qos
    geoscope
}

designboard: string . 

type ServiceTopology {
    topology
    designboard
}

intent: string . 
status: string . 

type FogFunction {
    name
    topology
    designboard
    intent 
    status
}
