# Causal fault localisation in dataflow systems

This repo aims to provide a realistic demonstration of how causal inference can be used for fault localisation in dataflow systems. It contains demo projects built with modern dataflow frameworks. In each project we setup experiments in which we introduce some disturbance in the system , either changing its logic or shifting the input distribution. We then observe the shift in the output, and compute attributions scores of each node in the dataflow graph. Attribution score reflects how much a given node impacted the overall output distribution shift. The experiment succeeds if the node at which we intervened received the highest attribution score.

This repo accompanies paper of the same name, to be presented at EuroMLSys'23. Arxiv version is [here](https://arxiv.org/abs/2304.11987).

At the moment the repo contains three demos. Each subfolder contains instructions that explain how to install and run the demo.

## SciPipe demo

[SciPipe](https://scipipe.org/) is a framework for scientific workflows, developed in Go. It is heavily command-line oriented, with most steps taking in files and producing files with a specified command. It was originally developed for computational biology, although it can undoubtly be used in other fields too. A particularly interesting feature of SciPipe is [audit files](https://rillabs.com/posts/provenance-reports-in-scientific-workflows), which we use extensively. Our demo in SciPipe is concerned with calculating GC-ratio in DNA fragments.
built 
## Node-RED demo

[Node-RED](https://nodered.org/) is a flow-based visual programming platform that is highly popular in IoT space. Our demo in Node-RED extends existing example built by OpenEEW team and found [here](https://github.com/openeew/openeew-nodered). It is a dashboard that allows to examine seismic activity in a specified area and detect signs of earthquackes. We enhance it with a [custom profiler node](https://github.com/bartbutenaers/node-red-contrib-msg-profiler) that implements vital graph traversal and data tracing capabilities.

# Seldon Core demo

[Seldon Core](https://github.com/SeldonIO/seldon-core) is a ML inference platform. Second version of Seldon Core is built with dataflow-oriented approach using Kafka streams. Despite strong emphasis on model serving, second version of the platform can be used as a generic dataflow engine. It also provides some APIs that leverage dataflow nature of the platform, namely `inspect` which we use. Our demo in Seldon Core v2 is re-implementation of the insurance claims processing demo found [here](https://github.com/mlatcl/fbp-vs-soa/).
