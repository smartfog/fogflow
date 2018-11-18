# FogFlow

[![FIWARE Security](https://nexus.lab.fiware.org/static/badges/chapters/processing.svg)](https://www.fiware.org/developers/catalogue/)
[![License: BSD-4-Clause](https://img.shields.io/badge/license-BSD%204%20Clause-blue.svg)](https://spdx.org/licenses/BSD-4-Clause.html)
[![Docker Status](https://img.shields.io/docker/pulls/fogflow/discovery.svg)](https://hub.docker.com/r/fogflow/discovery/)
[![Support badge](https://img.shields.io/badge/tag-fiware-orange.svg?logo=stackoverflow)](https://stackoverflow.com/questions/tagged/fiware)
<br>
[![Documentation badge](https://img.shields.io/readthedocs/fogflow.svg)](http://fogflow.readthedocs.org/en/latest/)
![Status](https://nexus.lab.fiware.org/repository/raw/public/static/badges/statuses/fogflow.svg)
[![Build Status](https://travis-ci.org/smartfog/fogflow.svg?branch=master)](https://travis-ci.org/smartfog/fogflow)
[![Swagger Validator](https://img.shields.io/swagger/valid/2.0/https/raw.githubusercontent.com/OAI/OpenAPI-Specification/master/examples/v2.0/json/petstore-expanded.json.svg)](https://app.swaggerhub.com/apis/fogflow/broker/1.0.0)

FogFlow is an IoT edge computing framework to automatically orchestrate dynamic
data processing flows over cloud and edges driven by context, including system
context on the available system resources from all layers, data context on the
registered metadata of all available data entities, and also usage context on
the expected QoS defined by users.

This project is part of [FIWARE](https://www.fiware.org/). For more information
check the FIWARE Catalogue entry for
[Processing](https://github.com/Fiware/catalogue/tree/master/processing).

## Content

-   [Background](#background)
-   [Installation](#installation)
-   [Usage](#usage)
-   [API](#api)
-   [Testing](#testing)
-   [Quality Assurance](#quality-assurance)
-   [More Information](#more-information)
-   [License](#license)

## Background

FogFlow is a standard-based data processing framework for service providers to
easily programming and managing IoT services over cloud and edges. Below are the
motivation, functionalities, and benefits of FogFlow.

-   _Why do we need FogFlow?_

    -   the cost of a cloud-only solution is too high to run a large scale IoT
        system with >1000 geo-distributed devices
    -   many IoT services require fast response time, such as <10ms end-to-end
        latency
    -   service providers are facing huge complexity and cost to fast design and
        deploy their IoT services in a cloud-edge environment
    -   business demands are changing fast over time and service providers need
        to try out and release any new services over their shared cloud-edge
        infrastructure at a fast speed
    -   lack of programming model to fast design and deploy IoT services over
        geo-distributed ICT infrastructure
    -   lack of interoperability and openness to share and reuse data and
        dervied results across various applications

-   _What does FogFlow provide?_

    -   efficient programming model: programming a service is like building lego
        blocks
    -   dynamic service orchestration: launching necessary data processing only
        when it is required
    -   optimized task deployment: assigning tasks between cloud and edges based
        on the locality of producers and consumers - scalable context
        management: allowing flexible information exchanging (both topic-based
        and scope-based) between producers and consumers

-   _How can customers benefit from FogFlow?_
    -   fast time-to-market when realizing and releasing new services over the
        shared, geo-distributed ICT infrastructure
    -   reduced operation cost and management complexity when operating variou
        services
    -   being able to provide services that require low latency and fast
        response time


## Installation

The instructions to install FogFlow can be found in the
[Installation Guide](https://fogflow.readthedocs.io/en/latest/setup.html)

## Usage

Information about how to use FogFlow can be found in the
[User & Programmers Manual](https://fogflow.readthedocs.io/en/latest/example1.html)

## API

APIs and examples of their usage can be found
[here](https://fogflow.readthedocs.io/en/latest/api.html)

## Testing

For performing a basic end-to-end test, you can follow the detailed instructions [here](https://fogflow.readthedocs.io/en/latest/test.html).

## Quality Assurance

This project is part of [FIWARE](https://fiware.org/) and has been rated as
follows:

-   **Version Tested:**
    ![ ](https://img.shields.io/badge/dynamic/json.svg?label=Version&url=https://fiware.github.io/catalogue/json/fogflow.json&query=$.version&colorB=blue)
-   **Documentation:**
    ![ ](https://img.shields.io/badge/dynamic/json.svg?label=Completeness&url=https://fiware.github.io/catalogue/json/fogflow.json&query=$.docCompleteness&colorB=blue)
    ![ ](https://img.shields.io/badge/dynamic/json.svg?label=Usability&url=https://fiware.github.io/catalogue/json/fogflow.json&query=$.docSoundness&colorB=blue)
-   **Responsiveness:**
    ![ ](https://img.shields.io/badge/dynamic/json.svg?label=Time%20to%20Respond&url=https://fiware.github.io/catalogue/json/fogflow.json&query=$.timeToCharge&colorB=blue)
    ![ ](https://img.shields.io/badge/dynamic/json.svg?label=Time%20to%20Fix&url=https://fiware.github.io/catalogue/json/fogflow.json&query=$.timeToFix&colorB=blue)
-   **FIWARE Testing:**
    ![ ](https://img.shields.io/badge/dynamic/json.svg?label=Tests%20Passed&url=https://fiware.github.io/catalogue/json/fogflow.json&query=$.failureRate&colorB=blue)
    ![ ](https://img.shields.io/badge/dynamic/json.svg?label=Scalability&url=https://fiware.github.io/catalogue/json/fogflow.json&query=$.scalability&colorB=blue)
    ![ ](https://img.shields.io/badge/dynamic/json.svg?label=Performance&url=https://fiware.github.io/catalogue/json/fogflow.json&query=$.performance&colorB=blue)
    ![ ](https://img.shields.io/badge/dynamic/json.svg?label=Stability&url=https://fiware.github.io/catalogue/json/fogflow.json&query=$.stability&colorB=blue)

## More Information

-   [Tutorial](http://fogflow.readthedocs.io/en/latest/index.html)
-   [IoT-J paper](http://ieeexplore.ieee.org/document/8022859/)

## License

FogFlow is licensed under
[BSD-4-Clause](https://spdx.org/licenses/BSD-4-Clause.html).

© 2017 NEC
