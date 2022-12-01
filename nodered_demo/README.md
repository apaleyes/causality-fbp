This demo showcases application of causal effect estimation to Node-RED.

## Demo specific requirements
Obviously, Node-RED. Running it with Docker is sufficient for the demo, instructions here: https://nodered.org/docs/getting-started/docker

We make extensive use of msg profiler node, which is found here: https://github.com/bartbutenaers/node-red-contrib-msg-profiler . Excellent documentation makes it easy to add this node to your Node-RED and to use it too!

The flow is a slightly extended version of Example 3 in OpenEEW demos: https://github.com/openeew/openeew-nodered . We add a few more nodes for added complexity.

## Description

The demo displays a dashboard on which user can plot seismic data for certain regions, mostly in Mexico. We extend it in several ways:
* Allow auto-play function to display the max of 600 minutes instead of 60.
* Add a node that detects earthquake and displays corresponding message.
* Add a delay node to simulate additional delay.
* Add profiler node.
* Add ability to simulate delay in the node that builds graph data.

We generate requests in the flow by running auto-play for its new max duration, and store all profiler data in `/data/trace.log`. Profiler data gives us exact time stamp of each node's message output, and in aggregation we can find out what was each node's throughput in a give period of time. Throughput can be manipulated by adding delays in the flow. The causal question we ask is: "Give the observed drop in throughput, which nodes caused it?"

Running demo flow multiple times and saving `trace.log` file after each run, we collect profiling information in following modes:
* Normal, before additional delays
* Randomized delay in "Build graphs" node
* Randomized delay in "Additional delay" node

That gives us the raw data to work with. Data aggregation and application of DoWhy to it can be found in the accompanying notebook.