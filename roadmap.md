# FogFlow Roadmap

FogFlow is a FIWARE Generic Enabler. If
you would like to learn about the overall Roadmap of FIWARE, please check
section "Roadmap" on the FIWARE Catalogue.

## Introduction

FogFlow is an IoT edge computing framework to enable the easy programming and dynamic deployment of IoT services over cloud and edges. It provides easy-to-use edge programming models (e.g., service topology and and fog function) and also developer-friendly user interface, allowing service developers to quickly design and deploy their IoT services with low development and management cost. Currently, we are working on advanced algorithms to ensure service QoS and also improve system scalability and relibility via autonomous service orchestration. 

## Short term

The following list of features are planned to be addressed in the short term,
and incorporated in the next release of the product planned for **28 Feb. 2019**:

-   Update the same entity from multiple brokers: different running task instances might update the same entity from different edge nodes and their entity updates should be propagated to the original host of the entity.  

-   Fulfill all GE requirements: to finish a simplify tutorial that can show how data source can be flowed into Orion Context Broker after triggering some fog function in the FogFlow system


## Medium term

The following list of features are planned to be addressed in the medium term,
typically within the subsequent release(s) generated in the next **9 months**
after next planned release:

- refine the current programming model with a expressive intent model, which allows service developers to specify their expected service QoS and configurable parameters of operators

- provide more use case examples with the new intent-based programming model


## Long term

The following list of features are proposals regarding the longer-term evolution
of the product even though development of these features has not yet been
scheduled for a release in the near future. Please feel free to contact us if
you wish to get involved in the implementation or influence the roadmap

- adapt to the new standard NGSI-LD

- make orchestration decisions to support edge AI and knowledge extraction, especially for enabling dynamic and efficient image processing pipelines at edges

